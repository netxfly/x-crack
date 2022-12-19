/*

Copyright (c) 2017 xsec.io

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THEq
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.

*/

package cmd

import (
	"github.com/urfave/cli"

	"x-crack/util"
)

var Scan = cli.Command{
	Name:        "scan",
	Usage:       "start to crack weak password",
	Description: "start to crack weak password",
	Action:      util.Scan,
	Flags: []cli.Flag{
		boolFlag("debug, d", "debug mode"),
		intFlag("timeout, t", 8, "timeout"),
		intFlag("scan_num, n", 100, "thread num"),

		stringFlag("ip_list, I", "iplist.txt", "iplist"),        //TODO: should notice
		stringFlag("user_dict, U", "user.dic", "user dict"),     //TODO: should notice
		stringFlag("pass_dict, P", "pass.dic", "password dict"), //TODO: should notice
		stringFlag("outfile, o", "x_crack.txt", "scan result file"),

		stringFlag("ip, i", "", "service ip addr, eg: ssh://127.0.0.1:22"),
		stringFlag("username, u", "", "username, eg: root"),
		stringFlag("password, p", "", "password, eg: 123456"),
	},
}

//var Dump = cli.Command{
//	Name:        "dump",
//	Usage:       "dump result to a text file",
//	Description: "dump result to a text file",
//	Action:      models.Dump,
//	Flags: []cli.Flag{
//		stringFlag("outfile, o", "x_crack.txt", "scan result file"),
//	},
//}

func stringFlag(name, value, usage string) cli.StringFlag {
	return cli.StringFlag{
		Name:  name,
		Value: value,
		Usage: usage,
	}
}

func boolFlag(name, usage string) cli.BoolFlag {
	return cli.BoolFlag{
		Name:  name,
		Usage: usage,
	}
}

func intFlag(name string, value int, usage string) cli.IntFlag {
	return cli.IntFlag{
		Name:  name,
		Value: value,
		Usage: usage,
	}
}
