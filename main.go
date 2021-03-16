package main

import (
	"net"
	"os"

	"github.com/apex/log"
	"github.com/jklaiber/dockconman/dockconman"
	"github.com/jklaiber/dockconman/pkg/rsahelper"
	"github.com/urfave/cli/v2"
)

var VERSION string

const DefaultSshKey = `-----BEGIN RSA PRIVATE KEY-----
MIIJKQIBAAKCAgEA0W67ase+PfZtTUxrhWmws+hYpdRTZj0CPEeNgRrsEcAO/xKy
X9lY2ETqjtJJtah80ce9i0FdFZNYKn8HIwz/fnKVZh5LIdcGi6x3jf01NuDG0Ek8
9m3og6DxOG78XmU8F3cFbroXLNPC6Q8r1zihqoeuQXsqt2b/1Eqqbx1BGWLzw8Xy
0ASZoeIMmwTG+jhIeQnHiCOBil2NfKG6LvmDvaZ6HR4JUpzd7npvtWkEsoiPnLs/
u/QaSEBE2BV+s+4y/fpKu4B6fBjjkMqAidvezOUeBybhUSs0oiTwBCbrpHWC44EA
kRpSIJvSiWB++2iWiU4DOsMQ0ef0kZt4+QoAFYcwKwb6ZrlgtLdtJ6d/P63OMJxL
TqPMnwhqihGIcRJ1LsQP3MzNpzTVltHPJp0WtlFORp9ES9Xd70DHqeN2zAOtxEmT
HwRMutMLTnZB3ETjfwMxB6pxtW+oZR61LP1/tHNk19p+4qlZARrZjU/eY/2EhoEK
13gxkbBvDURl+jvuTxMbtdgdcqNtzG6HFji6hREMchqzLZvuM4AhG2ySruC8hY50
I1G33BoSvZAc5rjb60K0sBoAJzkgl972x2w9JNZioxKUQ6a5WJadxQ55tI0SX95a
fUhRtIhTnpyAAXJqSvbQJsQMgjypLvVBPllSPESdBU/hK0JLn6G9x8D+aRcCAwEA
AQKCAgBelnNg66uZUpXVBoG9NJnQ90wqQTYVg9JhpTNcrusVrTdYrnoPXhuJOb7y
GDmgKOOO33ZU0YWX+/8i2lI/21v9IQUmpOHr+7CCHre0UjrZeTVx1tIIvmT4JhUs
Frw2aeR0+JVkh/l1joKGPgrf3jztxu/FtTn6sTM0DzDglEbVj2Jm9h0PJDS671wg
G00+r3LT773FV7vm4Q/IzUZIdvqwCeM3CVzOifiV/4g6V57+fzzVXaCQY9QG87fk
G/kojJlNKeDcxE8NgzQsLblWCg6bXZNtSXsT7L6NyL76MRXWJhiiZZ366vdSsO2q
jAFtzLPIeYpt3LHZC7jilmRRndmhDjtoUipWWJQ9fqfkALDdOLo/pQbaiXXt4MYQ
mtuhL4hvWz+cY9U29Ndnmu2IEn1ZDvqNiv5lDG0oq9Mv2sUjdWr7TaDOKOY9Ab8x
mrq98k3NqNceVp5bBeIuXCHCFEvqyYNAJtWepmypMdX3F6TrsknTMjOsMpW4sL30
dgxpy1kPOS4YTmVQAZSIWhJVYgbQcE2s/BgpwGGbOcFUdvqaQUEbgE4gNMBtGwGL
gTbvC5f74naQ6v3TsQhTjxN4puTRpcABpL/WtKy0AA5ayIt4MCFl+vLNTwzOI/gP
oHWU5ohKkOJQ6otEixKEGIc6NBIx1IkTiUtXuEHv4izNIfBDmQKCAQEA7i4AzEhx
l0JsxI72268mxaIEo/QCc4R5EBp55orSI65oU15zTyY1g9Vqar3FR3Nsh+PKKTmJ
9IdSl4mJkC+Tki2FNk7PLPR00di3yhar5xZHqrcqlK6ik30p453nTtG5mPLjgp0o
NWdRQxffx70v9MYgt4rCg0V16pt9EwxfXslnqRr6XCNr2vbfBceJs/xveH7Nbnwl
1LTShbI14bmxQ6/wSCKnZ14CBR9tzHbTV5cAv89pLqnrfCTnKcCzq14QhMa3C4uT
Rm+0dJzLfdYCypCqbiXYZti1eAAEhV/ppoZi0LKumwmV4ScVsj7WM1qDf6eQWC0Q
eSiebRRSZg9KiwKCAQEA4RoeJEnoB+X1YM9X5JoV3bQI5XBSoRykJuV81XoCb/gO
ypTs824csqHYH9wnldpDyk0wDTenVFvaNQT6EkgOtRjnESfS+GgBtngSNL27Ea84
0+GLqoj7WaMllqXHan3LwY/SeSozx6VWW0eU1XFFr1FdZ/M2jj4mY8AgmHAVu5Zo
YFkPTjohBCem2sLkYUe4Nlf079b6ASAt4I81Cm+rphEeBP7+ovc/HXJcL6uS/ZKW
iySO2zngT7z9LaXBL/0qExj3B/JRoZVonx3G+DVT2mGexmuQB4kf/gPYSfzsWhGL
451LNn5XyHMiclD545edloyl/seHsB5KMXKAGXVJJQKCAQA+EPLYSRCAsCiT2AVw
HeZmnd/DsbRp0d2SWrPlZct4zNwWzYgS2gwb/KMsiaM9CVEA4FUwBPR0KkdVgdu9
HQjBkOcjzcmjF1jRzj2mhd3p7B5k2DJaaF+pO3aM//rkyTYqKzEqOjXeJLxCVZhU
/nHewTqJWblyZ8lgh4BCVHkNxEIlCQiwtfJHLwnTAbpakq+hoLl7zxI0qaIqgNQV
rEQLNW/R/GXPQ+oW16fPHi/YpVrmoO/x3wmkYiFy+epX/70iPH46nfaU5ksKEEne
0sQLcUNYTLhlpJc1XBvRfbrvUBmz9LwXXpoWAA9hUYqT+0RFIa81qxid2f3ewurt
+ZIdAoIBAQDbhUTvzsNhMHljt9DXNw0r8G7ckfWC+RN8e0CKTzohR5/lH+cUXsXN
zted+m0ATqLdnvjFawjb09ew7PGS8oKlSWvN5zBu378L23ylwoG0dVTODJ7P6FZ2
zAvUJkebKqKSWVfAoc9tW2gkDGKw5I44sviMbzs87I8zqCIhhu0qyztu+mtatoWM
L78gh/+Afxi+pnhPjS6x+lfDLuVjEBQtF3RXGvXop4X9iZEtS/1FHLeDaluGn6KJ
IJ0m7wa/bfyiMy51qXLCSZqF0dxAIoFr7teQWUVUk/2HEujS/rzf+Uya5MJ8mimx
adal9SI9OZaNQwx+ssc4kdF491jFewOhAoIBAQCnAG2m7Y5ZDI7tcEdaeZsaGbFC
vE2WgIOwZ0/ptWOkI8cpPQVrmGT+WQGRU5vDOUwTiSy4EbZdtF0GFwFJVsHpwX64
nOLmTlsHfd72X9WFdW9zBRd2onsiWE3D1SqPUJHj4E20yK3k0GVKvq1y6NyhH3py
gZmuOE49IyhD9XJcl6HnumvgPenlxdUE318NqDBu/gKdqVFK3FaEy1lGHxc4c6SS
NeFNiHhqwhLMvh+cmbaJ7NyKbk2T3PCsrc+94ovbZW+X1tJNUHscJwHHl1sPjLsJ
OHJOxt65tQQADMrmZpye1di+aopRGLcj7hmxQgikUYDdww4oIWbuIosJWOMi
-----END RSA PRIVATE KEY-----`

const DefaultSshPublicKey = `-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEA0W67ase+PfZtTUxrhWmw
s+hYpdRTZj0CPEeNgRrsEcAO/xKyX9lY2ETqjtJJtah80ce9i0FdFZNYKn8HIwz/
fnKVZh5LIdcGi6x3jf01NuDG0Ek89m3og6DxOG78XmU8F3cFbroXLNPC6Q8r1zih
qoeuQXsqt2b/1Eqqbx1BGWLzw8Xy0ASZoeIMmwTG+jhIeQnHiCOBil2NfKG6LvmD
vaZ6HR4JUpzd7npvtWkEsoiPnLs/u/QaSEBE2BV+s+4y/fpKu4B6fBjjkMqAidve
zOUeBybhUSs0oiTwBCbrpHWC44EAkRpSIJvSiWB++2iWiU4DOsMQ0ef0kZt4+QoA
FYcwKwb6ZrlgtLdtJ6d/P63OMJxLTqPMnwhqihGIcRJ1LsQP3MzNpzTVltHPJp0W
tlFORp9ES9Xd70DHqeN2zAOtxEmTHwRMutMLTnZB3ETjfwMxB6pxtW+oZR61LP1/
tHNk19p+4qlZARrZjU/eY/2EhoEK13gxkbBvDURl+jvuTxMbtdgdcqNtzG6HFji6
hREMchqzLZvuM4AhG2ySruC8hY50I1G33BoSvZAc5rjb60K0sBoAJzkgl972x2w9
JNZioxKUQ6a5WJadxQ55tI0SX95afUhRtIhTnpyAAXJqSvbQJsQMgjypLvVBPllS
PESdBU/hK0JLn6G9x8D+aRcCAwEAAQ==
-----END PUBLIC KEY-----`

const DefaultBanner = `#############################
# Docker Connection Manager #
#############################
`

func main() {
	app := &cli.App{
		Name:  "dockconman",
		Usage: "simple ssh portal to containers",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "container_name",
				Aliases:  []string{"n"},
				Usage:    "Target container",
				EnvVars:  []string{"DOCKCONMAN_CONTAINER"},
				Required: true,
			},
			&cli.StringFlag{
				Name:    "command",
				Aliases: []string{"c"},
				Usage:   "Execute command on target container",
				EnvVars: []string{"DOCKCONMAN_COMMAND"},
				Value:   "bash",
			},
			&cli.StringFlag{
				Name:    "key-destination",
				Aliases: []string{"k"},
				Usage:   "Host key destination which should be taken",
				EnvVars: []string{"DOCKCONMAN_SSH_KEY_FILE"},
			},
			&cli.StringFlag{
				Name:    "default-key",
				Aliases: []string{"d"},
				Usage:   "Disable automatic key generation",
				Value:   "false",
			},
			&cli.StringFlag{
				Name:    "banner",
				Aliases: []string{"b"},
				Usage:   "Login banner",
				Value:   DefaultBanner,
			},
			&cli.StringFlag{
				Name:    "shell",
				Aliases: []string{"s"},
				Usage:   "Default shell",
				Value:   "/bin/sh",
			},
			&cli.StringFlag{
				Name:    "port",
				Aliases: []string{"p"},
				Usage:   "Binding port",
				EnvVars: []string{"DOCKCONMAN_PORT"},
				Value:   ":2222",
			},
		},
		Action: func(c *cli.Context) error {
			server, err := dockconman.NewServer()

			if err != nil {
				log.Fatalf("Cannot create server")
			}
			server.Banner = DefaultBanner
			server.DefaultShell = c.String("shell")
			server.DockerContainer = c.String("container_name")
			server.DockerExecArgs = c.String("command")
			server.Banner = c.String("banner")

			if c.String("key-destination") == "" && c.String("default-key") == "true" {
				server.AddHostKey(DefaultSshKey)
			} else if c.String("default-key") == "false" && c.String("key-destination") == "" {
				rsaKey, err := rsahelper.GetRSA(4096)
				if err != nil {
					log.Errorf("RSA key generation failed")
				}
				server.AddHostKey(rsaKey)
			} else if c.String("key-destination") != "" {
				rsaKey, err := rsahelper.RsaSetup(c.String("key-destination"))
				if err != nil {
					log.Errorf("RSA key loading failed")
				}
				server.AddHostKey(rsaKey)
			}

			bindAddress := c.String("port")
			listener, err := net.Listen("tcp", bindAddress)
			if err != nil {
				log.Fatalf("Failed to start listener on %q", bindAddress)
			}
			log.Infof("Listening on %q", bindAddress)

			if err = server.Init(); err != nil {
				log.Fatalf("Failed to initialize the server")
			}

			for {
				conn, err := listener.Accept()
				if err != nil {
					log.Errorf("Client acception failed")
					continue
				}
				go server.Handle(conn)
			}
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal("Application cannot be started!")
	}
}
