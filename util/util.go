package util

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

// CreateHeader will create date header fors output file
func CreateHeader() string {
	t := time.Now()
	return ("\nzzz ################################### " + t.Format("Mon Jan 2 03:04:05 MST 2006") + "\n")
}

// CreateSystemd will create service and timer files
func CreateSystemd(systemdType string, unitText string, name string) error {

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

// EnableSystemd enables the systemd timer
func EnableSystemd(service string) error {
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

// DisableSystemd disables the sytemd timer
func DisableSystemd(service string) error {
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

// DeleteSystemd deletes the systemd service or timer
func DeleteSystemd(service string) error {
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
