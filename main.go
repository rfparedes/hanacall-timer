package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"time"

	"github.com/rfparedes/hanacall-timer/util"
	flag "github.com/spf13/pflag"
)

const logfile string = "/var/log/hanacall-timer"
const exdir string = "/usr/local/bin"

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

	// Make sure user is running gdg out of /usr/local/bin
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

	if c.start == true && isFlagPassed("sidadm") == true {
		startService()
		return
	}

	if c.stop == true {
		stopService()
		return
	}

	if c.run == true && isFlagPassed("sidadm") == true {

		//Make sure sidadm is valid
		if _, err := user.Lookup(c.sidadm); err != nil {
			fmt.Printf("sidadmin '%s' is not valid\n", c.sidadm)
			return
		}

		f, err := os.OpenFile(logfile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0744)

		if err != nil {
			fmt.Printf("cannot open file '%s'", logfile)
		}
		defer f.Close()

		_, err = f.WriteString(util.CreateHeader() + "\n")
		if err != nil {
			fmt.Printf("cannot write to file '%s'", logfile)
		}

		cmd1 := exec.Command("su", "-", c.sidadm, "-c", "HDBSettings.sh systemReplicationStatus.py;echo rc=$?")
		cmd1.Stdout = f
		start1 := time.Now()
		err = cmd1.Start()
		if err != nil {
			fmt.Printf("cannot run cmd")
		}

		cmd2 := exec.Command("su", "-", c.sidadm, "-c", "HDBSettings.sh landscapeHostConfiguration.py;echo rc=$?")
		cmd2.Stdout = f
		start2 := time.Now()
		err = cmd2.Start()
		if err != nil {
			fmt.Printf("cannot run cmd")
		}

		cmd1.Wait()
		cmd2.Wait()

		_, err = fmt.Fprintf(f, "\nzzz\nTime spent in systemReplicationStatus.py    : %v\n", time.Since(start1))
		if err != nil {
			fmt.Printf("cannot write to file '%s'", logfile)
		}
		_, err = fmt.Fprintf(f, "Time spent in landscapeHostConfiguration.py : %v\n", time.Since(start2))
		if err != nil {
			fmt.Printf("cannot write to file '%s'", logfile)
		}
		return
	}

	flag.Usage()
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func startService() {

	timer := `[Unit]
Description=hanacall-timer Timer
Requires=hanacall-timer.service
		
[Timer]
OnActiveSec=0
OnUnitActiveSec=` + (fmt.Sprint(c.interval)) + "\n" +
		`AccuracySec=500msec

[Install]
WantedBy=timers.target`

	service := `[Unit]
Description=hanacall-timer Service
Wants=hanacall.timer
		
[Service]
Type=oneshot
ExecStart=` + exdir + "/hanacall-timer --run --sidadm " + (fmt.Sprint(c.sidadm)) + "\n"

	fmt.Printf("%s is now enabled and will run every %s seconds\n", progName, fmt.Sprint(c.interval))
	util.CreateSystemd("service", service, "hanacall-timer")
	util.CreateSystemd("timer", timer, "hanacall-timer")
	util.EnableSystemd("hanacall-timer.timer")
}

func stopService() {
	util.DisableSystemd("hanacall-timer.timer")
	util.DeleteSystemd("hanacall-timer.timer")
	util.DeleteSystemd("hanacall-timer.service")
}
