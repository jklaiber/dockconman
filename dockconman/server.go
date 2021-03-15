package dockconman

import (
	"io/ioutil"
	"net"
	"os"
	"strings"

	"github.com/apex/log"
	"golang.org/x/crypto/ssh"
)

type Server struct {
	SshConfig     *ssh.ServerConfig
	ClientConfigs map[string]*ClientConfig

	DefaultShell    string
	DockerContainer string
	DockerExecArgs  string
	Banner          string

	initialized bool
}

func NewServer() (*Server, error) {
	server := Server{}
	server.SshConfig = &ssh.ServerConfig{
		Config:       ssh.Config{},
		NoClientAuth: true,
		MaxAuthTries: 0,
	}
	server.ClientConfigs = make(map[string]*ClientConfig, 0)
	server.DefaultShell = "/bin/sh"
	return &server, nil
}

func (s *Server) Init() error {
	if s.initialized {
		return nil
	}

	s.SshConfig.PasswordCallback = nil
	s.initialized = true
	return nil
}

func (s *Server) Handle(netConn net.Conn) error {
	if err := s.Init(); err != nil {
		return err
	}

	log.Debugf("Server.Handle netConn=%v", netConn)

	conn, chans, reqs, err := ssh.NewServerConn(netConn, s.SshConfig)

	if err != nil {
		log.Infof("Received disconnect from %s", netConn.RemoteAddr().String())
		return err
	}

	client := NewClient(conn, chans, reqs, s)

	if err = client.HandleRequests(); err != nil {
		return err
	}

	if err = client.HandleChannels(); err != nil {
		return err
	}

	return nil
}

func (s *Server) AddHostKey(keystring string) error {
	keypath := os.ExpandEnv(strings.Replace(keystring, "~", "$HOME", 2))
	_, err := os.Stat(keypath)
	var keybytes []byte
	if err == nil {
		keybytes, err = ioutil.ReadFile(keypath)
		if err != nil {
			return err
		}
	} else {
		keybytes = []byte(keystring)
	}

	hostkey, err := ssh.ParsePrivateKey(keybytes)
	if err != nil {
		return err
	}

	s.SshConfig.AddHostKey(hostkey)
	return nil
}
