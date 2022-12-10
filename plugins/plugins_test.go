package plugins_test

import (
	"github.com/urfave/cli"
	"testing"
	"x-crack/cmd"
	"x-crack/models"
	"x-crack/plugins"
	"x-crack/vars"
)

func Cli(Args []string) {
	app := cli.NewApp()
	app.Name = "x-crack"
	app.Author = "netxfly"
	app.Email = "x@xsec.io"
	app.Version = "20171227"
	app.Usage = "Weak password scanner, Support: FTP/SSH/MSSQL/MYSQL/PostGreSQL/REDIS/ElasticSearch/MONGODB"
	app.Commands = []cli.Command{cmd.Scan /*cmd.Dump*/}
	app.Flags = append(app.Flags, cmd.Scan.Flags...)
	// app.Flags = append(app.Flags, cmd.Dump.Flags...)
	app.Run(Args)
}

// https://stackoverflow.com/questions/33723300/how-to-test-the-passing-of-arguments-in-golang
func TestCli(t *testing.T) {

	//oldArgs := os.Args
	//defer func() { os.Args = oldArgs }()
	Args := []string{
		"",
		"scan",
		"--ip=ssh://127.0.0.1:22",
		"--username=root",
		"--password=123456",
	}
	Cli(Args)

}

func TestAllProtocols(t *testing.T) {
	s := models.Service{Ip: "127.0.0.1", Port: 110, Protocol: "pop3", Username: vars.USER, Password: vars.PASS}
	t.Log(plugins.ScanPop3(s))
}
