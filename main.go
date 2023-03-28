package main

import (
	"bytes"
	"context"
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

// Tools is being developed for perfroming purple teaming
// Author: Ali.Alwashali
// purpose: to simulate complex command that require custom setup such as w3wp.exe spawning powershell or rundll32.exe making DNS connections

var app = cli.NewApp()

func init() {
	app.Name = "Malware Cli"

	app.Description = `Detection validation tool. 
	 The objective is to generate event with specific conditions to validate detection rule.
	 You can execute commands such as w3wp.exe spawning shell or winword creating file or making DNS queries.`

}

func filewrite(fpath string, binPath string) {

	if binPath != "" {
		directory := filepath.Dir(binPath)
		fileName := filepath.Base(binPath)

		copyBinaryTo(directory, fileName)

		args := []string{"createfile", "--path", fpath}
		output := execute(binPath, args)
		fmt.Println(output)
		return
	}

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

	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Millisecond * time.Duration(10000),
			}
			return d.DialContext(ctx, network, "8.8.8.8:53")
		},
	}
	ips, _ := r.LookupHost(context.Background(), hostname)

	for _, ip := range ips {
		fmt.Println(ip)
	}

}

func copyBinaryTo(toPath string, binaryName string) string {

	//get path of current running process
	srcbinary, err := os.Executable()
	if err != nil {
		panic(err)
	}

	parentBinary := ""
	if !strings.HasSuffix(binaryName, ".exe") {
		if strings.HasSuffix(toPath, "/") {

			parentBinary = toPath + binaryName + ".exe"

		} else {

			parentBinary = toPath + "/" + binaryName + ".exe"
		}

	} else {
		parentBinary = toPath + "/" + binaryName
	}

	in, err := os.Open(srcbinary)
	if err != nil {
		log.Println(err)
	}
	defer in.Close()
	out, err := os.Create(parentBinary)
	if err != nil {
		log.Println(err)
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
			log.Println(err)
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		log.Println(err)
	}
	err = out.Sync()
	if err != nil {
		log.Println(err)
	}

	return parentBinary
}

func prepareCommandArgs(command string, args []string) []string {
	c := []string{}

	c = append(c, "execute")
	c = append(c, "--command")
	c = append(c, command)

	if len(args) > 1 {

		c = append(c, "--arg")
		chainedArgs := ""
		for _, arg := range args {
			chainedArgs = fmt.Sprintf("%s %s", chainedArgs, arg)
		}
		chainedArgs = strings.TrimPrefix(chainedArgs, " ")
		c = append(c, chainedArgs)

	}

	return c

}

func execute(command string, args []string) string {

	log.Println("Executing:", command, args)

	cmd := exec.Command(command, args...)

	stdout, err := cmd.StdoutPipe()

	if err != nil {
		log.Println(err)
	}
	if err := cmd.Start(); err != nil {
		log.Println(err)
	}
	var out bytes.Buffer
	if _, err := io.Copy(&out, stdout); err != nil {
		log.Println(err)
	}
	if err := cmd.Wait(); err != nil {
		log.Println(err)
	}
	output := out.String()
	return output
}

func ExecuteCommand(commandName string, commandline string, parent string, copyPath string) (string, error) {

	args := strings.Split(commandline, " ")

	if copyPath != "C:/Users/Public" {
		_, err := os.Stat(copyPath)
		if os.IsNotExist(err) {
			if err := os.MkdirAll(copyPath, 0755); err != nil {
				log.Println(err)
			}

		}
	}

	// when passed for first time with custom parent process name
	if parent != "" {

		parentBinaryPath := copyBinaryTo(copyPath, parent)

		fullcommand := prepareCommandArgs(commandName, args)

		fmt.Println("outter command:", parentBinaryPath, fullcommand)

		output := execute(parentBinaryPath, fullcommand)

		fmt.Println(output)

		return output, nil

	} else {

		//when executed after the command is prepared with the custom parent
		fmt.Println("Inner command:", commandName, args)
		output := execute(commandName, args)
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

	log.Println("Received arguments: ", os.Args)
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
