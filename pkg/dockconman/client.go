package dockconman

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"

	"github.com/apex/log"
	"github.com/flynn/go-shlex"
	"github.com/jklaiber/dockconman/pkg/envhelper"
	"github.com/jklaiber/dockconman/pkg/ttyhelper"
	"github.com/kr/pty"
	"golang.org/x/crypto/ssh"
)

var clientCounter = 0

type Client struct {
	Idx        int
	ChannelIdx int
	Conn       *ssh.ServerConn
	Chans      <-chan ssh.NewChannel
	Reqs       <-chan *ssh.Request
	Server     *Server
	Pty, Tty   *os.File
	Config     *ClientConfig
	ClientID   string
}

type ClientConfig struct {
	ContainerName          string                `json:"image-name,omitempty"`
	RemoteUser             string                `json:"remote-user,omitempty"`
	Env                    envhelper.Environment `json:"env,omitempty"`
	Command                []string              `json:"command,omitempty"`
	DockerRunArgs          []string              `json:"docker-run-args,omitempty"`
	DockerExecArgs         []string              `json:"docker-exec-args,omitempty"`
	User                   string                `json:"user,omitempty"`
	Keys                   []string              `json:"keys,omitempty"`
	AuthenticationMethod   string                `json:"authentication-method,omitempty"`
	AuthenticationComment  string                `json:"authentication-coment,omitempty"`
	EntryPoint             string                `json:"entrypoint,omitempty"`
	AuthenticationAttempts int                   `json:"authentication-attempts,omitempty"`
	Allowed                bool                  `json:"allowed,omitempty"`
	IsLocal                bool                  `json:"is-local,omitempty"`
	UseTTY                 bool                  `json:"use-tty,omitempty"`
}

func NewClient(conn *ssh.ServerConn, chans <-chan ssh.NewChannel, reqs <-chan *ssh.Request, server *Server) *Client {
	client := Client{
		Idx:        clientCounter,
		ClientID:   conn.RemoteAddr().String(),
		ChannelIdx: 0,
		Conn:       conn,
		Chans:      chans,
		Reqs:       reqs,
		Server:     server,

		Config: &ClientConfig{
			RemoteUser:             "anonymous",
			AuthenticationMethod:   "noauth",
			AuthenticationComment:  "",
			AuthenticationAttempts: 0,
			Env:                    envhelper.Environment{},
			Command:                make([]string, 0),
		},
	}

	if _, found := server.ClientConfigs[client.ClientID]; !found {
		server.ClientConfigs[client.ClientID] = client.Config
	}

	client.Config = server.ClientConfigs[conn.RemoteAddr().String()]
	client.Config.Env.ApplyDefaults()

	clientCounter++

	log.Infof("accepted user %s from %s", conn.User(), conn.RemoteAddr())
	return &client
}

func (c *Client) HandleRequests() error {
	go func(in <-chan *ssh.Request) {
		for req := range in {
			log.Debugf("handleRequest: %v", req)
			if req.WantReply {
				req.Reply(false, nil)
			}
		}
	}(c.Reqs)
	return nil
}

func (c *Client) HandleChannels() error {
	for newChannel := range c.Chans {
		if err := c.HandleChannel(newChannel); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) HandleChannel(newChannel ssh.NewChannel) error {
	if newChannel.ChannelType() != "session" {
		log.Debugf("unknown channel type: %s", newChannel.ChannelType())
		newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
		return nil
	}

	channel, requests, err := newChannel.Accept()
	if err != nil {
		log.Errorf("newChannel.Accept failed")
		return err
	}
	c.ChannelIdx++
	log.Debugf("handleChannel.channel (client=%d channel=%d)", c.Idx, c.ChannelIdx)

	log.Debug("creating pty...")
	c.Pty, c.Tty, err = pty.Open()
	if err != nil {
		log.Errorf("pty.Open failed: PTY creation failed!")
		return nil
	}

	c.HandleChannelRequests(channel, requests)

	return nil
}

func (c *Client) HandleChannelRequests(channel ssh.Channel, requests <-chan *ssh.Request) {
	go func(in <-chan *ssh.Request) {
		defer c.Tty.Close()

		for req := range in {
			ok := false
			switch req.Type {
			case "shell":
				log.Debugf("handleChannelRequests.req shell")
				if len(req.Payload) != 0 {
					break
				}
				ok = true

				entrypoint := ""
				if c.Config.EntryPoint != "" {
					entrypoint = c.Config.EntryPoint
				}

				var args []string
				if c.Config.Command != nil {
					args = c.Config.Command
				}

				if entrypoint == "" && len(args) == 0 {
					args = []string{c.Server.DefaultShell}
				}
				c.runCommand(channel, entrypoint, args)

			case "exec":
				command := string(req.Payload[4:])
				log.Debugf("handleChannelRequests.req exec: %q", command)
				ok = true

				args, err := shlex.Split(command)
				if err != nil {
					log.Errorf("failed to parse command %q: %v", command, args)
				}
				c.runCommand(channel, c.Config.EntryPoint, args)

			case "pty-req":
				ok = true
				c.Config.UseTTY = true
				termLen := req.Payload[3]
				c.Config.Env["TERM"] = string(req.Payload[4 : termLen+4])
				c.Config.Env["USE_TTY"] = "1"
				w, h := ttyhelper.ParseDims(req.Payload[termLen+4:])
				ttyhelper.SetWinsize(c.Pty.Fd(), w, h)
				log.Debugf("handleChannelRequests.req pty-req: TERM=%q w=%q h=%q", c.Config.Env["TERM"], int(w), int(h))

			case "window-change":
				w, h := ttyhelper.ParseDims(req.Payload)
				ttyhelper.SetWinsize(c.Pty.Fd(), w, h)
				continue

			case "env":
				keyLen := req.Payload[3]
				key := string(req.Payload[4 : keyLen+4])
				valueLen := req.Payload[keyLen+7]
				value := string(req.Payload[keyLen+8 : keyLen+8+valueLen])
				log.Debugf("handleChannelRequets.req 'env': %s=%q", key, value)
				c.Config.Env[key] = value

			default:
				log.Debugf("unhandled request type: %q: %v", req.Type, req)
			}

			if req.WantReply {
				if !ok {
					log.Debugf("declining %s request...", req.Type)
				}
				req.Reply(ok, nil)
			}
		}
	}(requests)
}

func (c *Client) runCommand(channel ssh.Channel, entrypoint string, command []string) {
	var cmd *exec.Cmd
	var err error

	args := []string{"exec"}
	args = append(args, "-it")
	args = append(args, c.Server.DockerContainer)
	splitDockerExecArgs := strings.Split(c.Server.DockerExecArgs, " ")
	args = append(args, splitDockerExecArgs...)
	cmd = exec.Command("docker", args...)
	cmd.Env = c.Config.Env.List()

	if c.Server.Banner != "" {
		banner := c.Server.Banner
		banner = strings.Replace(banner, "\r", "", -1)
		banner = strings.Replace(banner, "\n", "\n\r", -1)
		fmt.Fprintf(channel, "%s\n\r", banner)
	}

	cmd.Stdout = channel
	cmd.Stdin = channel
	cmd.Stderr = channel
	var wg sync.WaitGroup
	if c.Config.UseTTY {
		cmd.Stdout = c.Tty
		cmd.Stdin = c.Tty
		cmd.Stderr = c.Tty

		wg.Add(1)
		go func() {
			io.Copy(channel, c.Pty)
			wg.Done()
		}()
		wg.Add(1)
		go func() {
			io.Copy(c.Pty, channel)
			wg.Done()
		}()
		defer wg.Wait()
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setctty: c.Config.UseTTY,
		Setsid:  true,
	}

	err = cmd.Start()
	if err != nil {
		log.Warnf("cmd.Start failed: %v", err)
		channel.Close()
		return
	}

	if err := cmd.Wait(); err != nil {
		log.Warnf("cmd.Wait failed: %v", err)
	}
	channel.Close()
	log.Debugf("cmd.Wait done")
}
