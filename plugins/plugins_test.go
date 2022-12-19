package plugins_test

import (
	"fmt"
	"github.com/urfave/cli"
	"testing"
	"x-crack/cmd"
)

func Cli(Args []string) {
	app := cli.NewApp()
	app.Usage = "Weak password scanner, Support: FTP/SSH/MSSQL/MYSQL/PostGreSQL/REDIS/ElasticSearch/MONGODB"
	app.Commands = []cli.Command{cmd.Scan /*cmd.Dump*/}
	app.Flags = append(app.Flags, cmd.Scan.Flags...)
	// app.Flags = append(app.Flags, cmd.Dump.Flags...)
	app.Run(Args)
}

func StartCli(Args []string, debug int) {
	// debug=0 by default
	Final := []string{
		"",
		"scan",
		fmt.Sprintf("-d=%v", debug),
	}
	for _, arg := range Args {
		Final = append(Final, arg)
	}
	Cli(Final)
}

// https://stackoverflow.com/questions/33723300/how-to-test-the-passing-of-arguments-in-golang
func TestCli(t *testing.T) {
	a1 := []string{
		"-u=root",
		"-p=123456",
		"-i=ssh://127.0.0.1:22",
	}
	a2 := []string{
		"-u=root",
		"-P=/Users/dpdu/GolandProjects/x-crack/pass.dic",
		"-i=ssh://127.0.0.1:22",
	}
	//fmt.Println(a1)
	//fmt.Println(a2)
	all := [][]string{
		a1,
		a2,
	}
	for i, a := range all {
		fmt.Println(fmt.Sprintf("------------------------Testing %v starts------------------------", i))
		StartCli(a, 1)
		fmt.Println(fmt.Sprintf("------------------------Testing %v ends------------------------", i))
	}

}

func TestAllProtocols(t *testing.T) {
	// see: vars.PortNames
	services := map[int]string{
		21:  "FTP",
		22:  "SSH",
		161: "SNMP",
		445: "SMB",
		//1433:  "MSSQL",
		3306: "MYSQL",
		//3389:  "RDP",
		5432: "POSTGRES",
		6379: "REDIS",
		//9200:  "ELASTICSEARCH",
		//27017: "MONGODB",
		110: "POP3",
	}
	for port, service := range services {
		ipaddr := fmt.Sprintf("%v://127.0.0.1:%v", service, port)
		input := []string{
			"-u=root",
			"-P=/Users/dpdu/GolandProjects/x-crack/pass.dic",
			"-i=" + ipaddr,
		}
		StartCli(input, 0)
	}
}
