package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

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
