package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rfparedes/hanacall-timer/util"
	flag "github.com/spf13/pflag"
)

const (
	progName string = "hanacall-timer"
	ver      string = "0.0.9"
	exdir    string = "/usr/local/bin"
	interval string = "60"
)

func main() {

	startCommand := flag.NewFlagSet("start", flag.ExitOnError)
	stopCommand := flag.NewFlagSet("stop", flag.ExitOnError)
	statusCommand := flag.NewFlagSet("status", flag.ExitOnError)
	runCommand := flag.NewFlagSet("run", flag.ExitOnError)
	verCommand := flag.NewFlagSet("version", flag.ExitOnError)
	sidadmFlag := flag.String("sidadm", "", "sidadm user")

	if len(os.Args) < 2 {
		util.Usage()
		return
	}

	// Make sure user is running gdg out of /usr/local/bin
	ex, err := os.Executable()
	if err != nil {
		return
	}
	exPath := filepath.Dir(ex)
	if exPath != exdir {
		fmt.Printf("%s binary must be in %s\n", progName, exdir)
		return
	}

	switch os.Args[1] {
	case "start":
		sidadmFlag = startCommand.String("sidadm", "", "sidadm user")
		startCommand.Parse(os.Args[2:])
	case "stop":
		stopCommand.Parse(os.Args[2:])
	case "status":
		statusCommand.Parse(os.Args[2:])
	case "version":
		verCommand.Parse(os.Args[2:])
	case "run":
		sidadmFlag = runCommand.String("sidadm", "", "sidadm user")
		runCommand.Parse(os.Args[2:])
	default:
		util.Usage()
	}

	if verCommand.Parsed() {
		fmt.Println(progName + " v" + ver + " (https://github.com/rfparedes/hanacall-timer)")
	}

	if statusCommand.Parsed() {
		if util.IsStarted() {
			fmt.Println("started")
		} else {
			fmt.Println("stopped")
		}
	}

	if startCommand.Parsed() {
		if *sidadmFlag == "" {
			fmt.Println("Please supply a sidadm username")
			return
		}
		if !util.SidadmValid(*sidadmFlag) {
			fmt.Println("sidadm is not valid. Please supply an existing user")
			return
		}
		util.StartService(interval, *sidadmFlag, progName, exdir)
	}

	if stopCommand.Parsed() {
		util.StopService()
	}

	if runCommand.Parsed() {
		if *sidadmFlag == "" {
			fmt.Println("Please supply a sidadm username, for example: --sidadm rfpadm")
			return
		}
		if !util.SidadmValid(*sidadmFlag) {
			fmt.Println("sidadm is not valid. Please supply an existing user")
			return
		}
		util.RunTimer(*sidadmFlag)
	}
}
