package util

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"time"
)

// create date header for output file
func createHeader() string {
	t := time.Now()
	return ("\n-----------------------------\n" + t.Format("Mon Jan 2 03:04:05 MST 2006") + "\n")
}

// create service and timer files
func createSystemd(systemdType string, unitText string, name string) error {
	fullPath := ("/etc/systemd/system/" + name + "." + systemdType)
	fmt.Printf("Creating systemd '%s'\n", systemdType)
	// Create systemd files
	f, err := os.OpenFile(fullPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("cannot open file: %s", fullPath)
	}
	defer f.Close()

	_, err = f.WriteString(unitText)
	if err != nil {
		return fmt.Errorf(("cannot write the unit content to the file"))
	}
	f.Sync()
	return nil
}

// enables the systemd timer
func enableSystemd(service string) error {
	systemctl, err := exec.LookPath("systemctl")
	if err != nil {
		return fmt.Errorf("cannot find 'systemctl' executable ")
	}
	fmt.Printf("Enabling systemd '%s'", service)
	enableCmd := exec.Command(systemctl, "enable", service, "--now")
	err = enableCmd.Run()
	if err != nil {
		return fmt.Errorf("cannot enable '%s'", service)
	}
	return nil
}

// disables the sytemd timer
func disableSystemd(service string) error {
	unitPath := ("/etc/systemd/system/" + service)
	systemctl, err := exec.LookPath("systemctl")
	if err != nil {
		return fmt.Errorf("cannot find 'systemctl' executable")
	}
	fmt.Println("Disabling systemd timer")
	// Check for systemd timer file
	if _, err := os.Stat(unitPath); err == nil {
		disableCmd := exec.Command(systemctl, "disable", service, "--now")
		err = disableCmd.Run()
		if err != nil {
			return fmt.Errorf("cannot disable '%s'", service)
		}
	} else {
		return fmt.Errorf("cannot disable nonexistent service '%s'", service)
	}
	return nil
}

// deletes the systemd service or timer
func deleteSystemd(service string) error {
	unitPath := ("/etc/systemd/system/" + service)
	if _, err := os.Stat(unitPath); err == nil {
		fmt.Printf("Removing systemd '%s'\n", service)
		err := os.Remove(unitPath)
		if err != nil {
			return fmt.Errorf("cannot remove '%s'", service)
		}
	} else {
		return fmt.Errorf("cannot delete nonexistent service '%s'", service)
	}
	return nil
}

// SidadmValid will return if sidadm is a valid username
func SidadmValid(sidadm string) bool {
	if _, err := user.Lookup(sidadm); err != nil {
		return false
	}
	return true
}

//StartService will create systemd service and enable timer
func StartService(interval string, sidadm string, progName string, exdir string) {

	timer := `[Unit]
Description=hanacall-timer Timer
Requires=hanacall-timer.service
		
[Timer]
OnActiveSec=0
OnUnitActiveSec=` + (fmt.Sprint(interval)) + "\n" +
		`AccuracySec=500msec

[Install]
WantedBy=timers.target`

	service := `[Unit]
Description=hanacall-timer Service
Wants=hanacall.timer
		
[Service]
Type=oneshot
ExecStart=` + exdir + "/hanacall-timer run --sidadm " + (fmt.Sprint(sidadm)) + "\n"

	fmt.Printf("%s is now enabled and will run every %s seconds\n", progName, fmt.Sprint(interval))
	createSystemd("service", service, "hanacall-timer")
	createSystemd("timer", timer, "hanacall-timer")
	enableSystemd("hanacall-timer.timer")
}

func StopService() {
	disableSystemd("hanacall-timer.timer")
	deleteSystemd("hanacall-timer.timer")
	deleteSystemd("hanacall-timer.service")
}

func RunTimer(sidadm string, logfile string) error {
	f, err := os.OpenFile(logfile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0744)
	if err != nil {
		return fmt.Errorf("cannot open file '%s'", logfile)
	}
	defer f.Close()
	_, err = f.WriteString(createHeader() + "\n")
	if err != nil {
		return fmt.Errorf("cannot write to file '%s'", logfile)
	}
	cmd1 := exec.Command("su", "-", sidadm, "-c", "HDBSettings.sh systemReplicationStatus.py;echo rc=$?")
	cmd1.Stdout = f
	start1 := time.Now()
	err = cmd1.Start()
	if err != nil {
		return fmt.Errorf("cannot run cmd")
	}

	cmd2 := exec.Command("su", "-", sidadm, "-c", "HDBSettings.sh landscapeHostConfiguration.py;echo rc=$?")
	cmd2.Stdout = f
	start2 := time.Now()
	err = cmd2.Start()
	if err != nil {
		return fmt.Errorf("cannot run cmd")
	}
	cmd1.Wait()
	cmd2.Wait()
	_, err = fmt.Fprintf(f, "\nzzz #######################################################\nTime spent in systemReplicationStatus.py    : %v\n", time.Since(start1))
	if err != nil {
		return fmt.Errorf("cannot write to file '%s'", logfile)
	}
	_, err = fmt.Fprintf(f, "Time spent in landscapeHostConfiguration.py : %v\n###########################################################\n", time.Since(start2))
	if err != nil {
		return fmt.Errorf("cannot write to file '%s'", logfile)
	}
	return nil
}

// Usage prints help
func Usage() {
	fmt.Println("hanacall-timer: HANA_CALL timer")
	fmt.Println("Start timer:")
	fmt.Println("   hanacall-timer start --sidadm [sidadm]")
	fmt.Println("Stop timer:")
	fmt.Println("   hanacall-timer stop")
	fmt.Println("Run timer one-shot (for systemd):")
	fmt.Println("   hanacall-timer run --sidadm [sidadm]")
	fmt.Println("Print hanacall-timer version: ")
	fmt.Println("   hanacall-timer version")
	fmt.Println("Print this message:")
	fmt.Println("   hanacall-timer help")
}
