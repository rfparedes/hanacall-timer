package main

import (
	"flag"
	"fmt"
	"github.com/rfparedes/hanacall-timer/util"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const timer string = `[Unit]
Description=hanacall-timer Timer
Requires=hanacall.service
	
[Timer]
OnActiveSec=0
OnUnitActiveSec=60
AccuracySec=500msec
	
[Install]
WantedBy=timers.target`

const service string = `[Unit]
Description=hanacall-timer Service
Wants=hanacall.timer
	
[Service]
Type=oneshot
ExecStart=/usr/local/bin/hanacall-timer`

const logfile string = "/var/log/hanacall-timer"

type config struct {
	version  bool
	stop     bool
	start    bool
	interval uint
	run      bool
	sidadm   string
}

func (c *config) setup() {
	flag.BoolVar(&c.version, "version", false, "Output version information")
	flag.BoolVar(&c.start, "start", false, "Start logging HANA_CALL")
	flag.BoolVar(&c.stop, "stop", false, "Stop logging HANA_CALL")
	flag.UintVar(&c.interval, "interval", 60, "How often to run in seconds")
	flag.BoolVar(&c.run, "run", false, "Run once (for systemd timer)")
	flag.StringVar(&c.sidadm, "sidadm", "sidadm", "sidadm user")
}

const (
	progName string = "hanacall-timer"
	ver      string = "0.0.1"
)

var c = config{}

func main() {

	// Make sure user is running gdg out of /usr/local/bbin/
	ex, err := os.Executable()
	if err != nil {
		return
	}
	exPath := filepath.Dir(ex)
	if exPath != "/usr/local/bin" {
		fmt.Printf("%s binary must be in /usr/local/bin\n", progName)
		return
	}

	c.setup()
	flag.Parse()

	// User requests version
	if c.version == true {
		fmt.Println(progName + " v" + ver + " (https://github.com/rfparedes/hanacall-timer)")
		return
	}

	if c.run == true {

		f, err := os.OpenFile(logfile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0744)

		if err != nil {
			fmt.Printf("cannot open file '%s'", logfile)
		}
		defer f.Close()

		_, err = f.WriteString(util.CreateHeader() + "\n")
		if err != nil {
			fmt.Printf("cannot write to file '%s'", logfile)
		}

		cmd1 := exec.Command("su", "-", sidadm, "-c", "HDBSettings.sh systemReplicationStatus.py;echo rc=$?")
		cmd1.Stdout = f
		start1 := time.Now()
		err = cmd1.Start()
		if err != nil {
			fmt.Printf("cannot run cmd")
		}

		cmd2 := exec.Command("su", "-", sidadm, "-c", "HDBSettings.sh landscapeHostConfiguration.py;echo rc=$?")
		cmd2.Stdout = f
		start2 := time.Now()
		err = cmd2.Start()
		if err != nil {
			fmt.Printf("cannot run cmd")
		}

		cmd1.Wait()
		cmd2.Wait()
		_, err = fmt.Fprintf(f, "\nTime spent in systemReplicationStatus.py    : %v\n", time.Since(start1))
		if err != nil {
			fmt.Printf("cannot write to file '%s'", logfile)
		}
		_, err = fmt.Fprintf(f, "Time spent in landscapeHostConfiguration.py : %v\n", time.Since(start2))
		if err != nil {
			fmt.Printf("cannot write to file '%s'", logfile)
		}
	}
}
