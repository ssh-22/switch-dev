package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var SWITCH_HOST_IP string

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a command: local or status")
		os.Exit(1)
	}

	command := os.Args[1]

	SWITCH_HOST_IP = os.Getenv("SWITCH_HOST_IP")
	if SWITCH_HOST_IP == "" {
		fmt.Println("SWITCH_HOST_IP environment variable is not set.")
		os.Exit(1)
	}

	switch command {
	case "switch":
		switchEnv()
	case "status":
		printStatus()
	default:
		fmt.Println("Invalid command. Valid commands are: switch, status")
		os.Exit(1)
	}
}

func readHostsFile() ([]string, error) {
	hosts, err := ioutil.ReadFile("/private/etc/hosts")
	if err != nil {
		fmt.Println("Error opening hosts file:", err)
		os.Exit(1)
	}

	lines := strings.Split(string(hosts), "\n")
	return lines, err
}

func switchEnv() error {

	lines, err := readHostsFile()

	var hasSwitchHost = false

	for i, line := range lines {
		{
			if strings.HasPrefix(line, "# "+SWITCH_HOST_IP) {
				lines[i] = strings.TrimPrefix(line, "# ")
				hasSwitchHost = true
				fmt.Println("Successfully switching to local")
				break
			} else if strings.HasPrefix(line, SWITCH_HOST_IP) {
				lines[i] = "# " + line
				hasSwitchHost = true
				fmt.Println("Successfully switching to dev")
				break
			}
		}
	}

	if !hasSwitchHost {
		fmt.Println("Error cannot find SWITCH_HOST_IP in hosts file:", SWITCH_HOST_IP)
		os.Exit(1)
	}

	err = ioutil.WriteFile("/private/etc/hosts", []byte(strings.Join(lines, "\n")), 0644)
	if err != nil {
		fmt.Println("Error writing to hosts file:", err)
		os.Exit(1)
	}
	return err
}

func printStatus() {
	status, err := getStatus()
	if err != nil {
		fmt.Println("Error calling getStatus:", err)
		os.Exit(1)
	}
	fmt.Println(status)
}

func getStatus() (string, error) {
	lines, err := readHostsFile()

	if err != nil {
		fmt.Println("Error reading hosts file:", err)
		os.Exit(1)
	}

	for _, line := range lines {
		{
			if strings.HasPrefix(line, "# "+SWITCH_HOST_IP) {
				return "dev", err
			} else if strings.HasPrefix(line, SWITCH_HOST_IP) {
				return "local", err
			}
		}
	}

	fmt.Println("Error cannot find SWITCH_HOST_IP in hosts file:", SWITCH_HOST_IP)
	os.Exit(1)
	return "", err
}
