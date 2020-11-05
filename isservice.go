package isservice

import (
	"fmt"
	"time"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

// Isservice takes 3 arguments (service, command and repeate - where repeate
// is number of retries waiting for service to get wanted status 1 min. pr. repeate)
// Example:
// Isservice("ismetering", "install", 10)
func Isservice(svcName string, cmd string, repeate int) (string, error) {
	var err error
	var status string

	switch cmd {
	case "start":
		err = startService(svcName)
		status = " Started"
	case "stop":
		err = controlService(svcName, svc.Stop, svc.Stopped, repeate)
		status = " Stopped"
	case "status":
		err = srvstatus(svcName, svc.Stop)
	}

	// err := controlService(svcName, svc.Stop, svc.Stopped)
	if err == nil {
		fmt.Println(svcName + status)
	} else {
		if err.Error() == "Access is denied." {
			fmt.Println("\n\n###############################################################")
			fmt.Println("ERROR!! ERROR!! ERROR!!")
			status = "ERROR: Could not stop service. Error message:"
			fmt.Println(status)
			fmt.Println(err)
			fmt.Println("Exiting caused lack of prievelegies to do my work")
			fmt.Println("###############################################################")
		}
	}
	return status, err
}

func srvstatus(name string, c svc.Cmd) error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(name)
	if err != nil {
		return fmt.Errorf("could not access service: %v", err)
	}
	defer s.Close()
	status, err := s.Control(c)
	if err != nil {
		return fmt.Errorf("could not send control=%d: %v", c, err)
	}
	fmt.Println("Service status: ")
	fmt.Sprint(status)
	return nil

}
func controlService(name string, c svc.Cmd, to svc.State, repeate int) error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(name)
	if err != nil {
		return fmt.Errorf("could not access service: %v", err)
	}
	defer s.Close()
	status, err := s.Control(c)
	if err != nil {
		return fmt.Errorf("could not send control=%d: %v", c, err)
	}
	timeout := time.Now().Add(time.Duration(repeate) * time.Minute)
	// timeout := time.Now().Add(time.Duration(repeate) * time.Second)
	for status.State != to {
		if timeout.Before(time.Now()) {
			return fmt.Errorf("timeout waiting for service to go to state=%d", to)
		}
		time.Sleep(300 * time.Millisecond)
		status, err = s.Query()
		if err != nil {
			return fmt.Errorf("could not retrieve service status: %v", err)
		}
	}
	return nil
}
func startService(name string) error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(name)
	if err != nil {
		return fmt.Errorf("could not access service: %v", err)
	}
	defer s.Close()
	err = s.Start("is", "manual-started")
	if err != nil {
		return fmt.Errorf("could not start service: %v", err)
	}
	return nil
}
