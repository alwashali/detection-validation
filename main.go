package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

var app = cli.NewApp()

func init() {
	app.Name = "Malware Cli"

	app.Description = `Detection validation tool. 
	 The objective is to generate event with specific conditions to validate detection rule.
	 You can execute commands such as w3wp.exe spawning shell or winword creating file or making DNS queries.`

}

func main() {
	app.Commands = []cli.Command{
		{
			Name: "argsfree",

			Usage: "Accept any commandline",

			Action: func(c *cli.Context) {

				fmt.Println(os.Args)

			},
		},
		{
			Name: "connect",

			Usage: "Connect to host",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "host",
					Usage:    "hostname or IP Address",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "port",
					Usage:    "port number",
					Required: true,
				},
			},
			Action: func(c *cli.Context) {

				connectToHost(c.String("host"), c.String("port"))

			},
		},
		{
			Name: "download",

			Usage: "Download file",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "url",
					Usage:    "File URL",
					Required: true,
				},
			},
			Action: func(c *cli.Context) {

				downloadFile(c.String("URL"))

			},
		},
		{
			Name: "dnsquery",

			Usage: "Resolve DNS",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "host",
					Usage:    "hostname to resolve",
					Required: true,
				},
			},
			Action: func(c *cli.Context) {

				resolve(c.String("host"))

			},
		},
		{
			Name:  "execute",
			Usage: "Execute command with custom commandline and parent process",
			Flags: []cli.Flag{

				&cli.StringFlag{
					Name:     "command",
					Usage:    "Hostname or IP Address",
					Required: true,
				},
				&cli.StringFlag{
					Name:  "parent",
					Usage: "Optinal parent command to execute",
				},
				&cli.StringFlag{
					Name:  "arg",
					Usage: "Command arguments",
				},
				&cli.StringFlag{
					Name:  "copy",
					Usage: "Copy to path before execution",
					Value: "C:/Users/Public",
				},
			},
			Action: func(c *cli.Context) {

				ExecuteCommand(c.String("command"), c.String("arg"), c.String("parent"), c.String("copy"))

			},
		},
		{
			Name: "createfile",

			Usage: "Create file at a spcific path",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "path",
					Usage:    "full path and file name",
					Required: true,
				},
				&cli.StringFlag{
					Name:  "binpath",
					Usage: "Full path of the binary creating the file. Example: C:/temp/binary.exe",
				},
			},
			Action: func(c *cli.Context) {

				filewrite(c.String("path"), c.String("binpath"))

			},
		},
	}

	//log.Println("Received arguments: ", os.Args)
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
