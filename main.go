package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/urfave/cli"
)

// Tools is being developed for perfroming purple teaming and detection validation
// Author: Ali.Alwashali
// purpose: to simulate complex command that require custom setup such as w3wp.exe spawning shell

var app = cli.NewApp()

func filewrite(fpath string) {

	_, err := os.Stat(fpath)
	if os.IsNotExist(err) {
		dir, file := filepath.Split(fpath)
		fmt.Println(dir, "file:", file)
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Println(err)
		}
		_, err := os.Create(fpath)
		if err != nil {
			log.Println(err)

		}

	}
}

func resolve(hostname string) {

	ip, err := net.LookupIP(hostname)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, ip := range ip {
		fmt.Println(ip)
	}
}

func copyFile(src, dst string) (err error) {

	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

func prepareCommandArgs(command string, args []string) string {
	c := ""
	c = fmt.Sprintf("execute --command %s ", command)

	for _, arg := range args {

		c = fmt.Sprintf("%s --arg %s", c, arg)

	}

	return c

}

func ExecuteCommand(commandName string, args []string, parent string) (string, error) {

	// when passed for first time with custom parent process name
	if parent != "" {

		//get path of current running process
		srcbinary, err := os.Executable()
		if err != nil {
			panic(err)
		}

		// exPath := filepath.Dir(srcbinary)
		exPath := "C:\\Users\\Public\\"

		parentBinary := ""
		if !strings.HasSuffix(parent, ".exe") {
			parentBinary = exPath + "\\" + parent + ".exe"
		}
		parentBinary = exPath + "\\" + parent

		copyFile(srcbinary, parentBinary)

		command := prepareCommandArgs(commandName, args)

		fmt.Println(parentBinary, command)

		cmd := exec.Command(parentBinary, command)

		stdout, err := cmd.StdoutPipe()

		if err != nil {
			return "", err
		}
		if err := cmd.Start(); err != nil {
			return "", err
		}
		var out bytes.Buffer
		if _, err := io.Copy(&out, stdout); err != nil {
			return "", err
		}
		if err := cmd.Wait(); err != nil {
			return "", err
		}
		output := out.String()
		fmt.Println(output)
		return output, nil
	} else {

		//when executed after the command is prepared with the custom parent
		fmt.Println("second execution", commandName, args)

		cmd := exec.Command(commandName, args...)

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return "", err
		}
		if err := cmd.Start(); err != nil {
			return "", err
		}
		var out bytes.Buffer
		if _, err := io.Copy(&out, stdout); err != nil {
			return "", err
		}
		if err := cmd.Wait(); err != nil {
			return "", err
		}
		output := out.String()
		fmt.Println(output)
		return output, nil
	}
}

func connectToHost(host string, port string) {
	timeout := time.Second * 3
	conn, err := net.DialTimeout("tcp", host+":"+port, timeout)
	if err != nil {
		fmt.Println("Connecting error:", err)
	}
	if conn != nil {
		defer conn.Close()
		fmt.Println("Opened", net.JoinHostPort(host, port))
	}
}

func main() {
	app.Commands = []cli.Command{
		{
			Name: "argsfree",

			Usage: "use to execute command with any commandline",

			Action: func(c *cli.Context) {

				fmt.Println(os.Args)

			},
		},
		{
			Name: "connect",

			Usage: "connect to host",
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
			Name: "dnsquery",

			Usage: "connect to host",
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
			Name: "execute",

			Usage: "execute command",
			Flags: []cli.Flag{

				&cli.StringFlag{
					Name:     "command",
					Usage:    "hostname or IP Address",
					Required: true,
				},
				&cli.StringFlag{
					Name:  "parent",
					Usage: "Optinal parent command to execute",
				},
				&cli.StringSliceFlag{
					Name:  "arg",
					Usage: "Command arguments",
				},
			},
			Action: func(c *cli.Context) {

				ExecuteCommand(c.String("command"), c.StringSlice("arg"), c.String("parent"))

			},
		},
		{
			Name: "createfile",

			Usage: "connect to host",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "path",
					Usage:    "full path and file name",
					Required: true,
				},
			},
			Action: func(c *cli.Context) {

				filewrite(c.String("path"))

			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
