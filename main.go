package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a command: switch or status")
		os.Exit(1)
	}

	command := os.Args[1]

	switchHostIP := os.Getenv("SWITCH_HOST_IP")
	if switchHostIP == "" {
		fmt.Println("SWITCH_HOST_IP environment variable is not set.")
		os.Exit(1)
	}

	switch command {
	case "switch":
		if err := switchEnv(switchHostIP); err != nil {
			fmt.Println("Failed to switch environment:", err)
			os.Exit(1)
		}
	case "status":
		status, err := getStatus(switchHostIP)
		if err != nil {
			fmt.Println("Failed to get status:", err)
			os.Exit(1)
		}
		fmt.Println(status)
	default:
		fmt.Println("Invalid command. Valid commands are: switch, status")
		os.Exit(1)
	}
}

func readHostsFile() (string, error) {
	hosts, err := ioutil.ReadFile("/private/etc/hosts")
	return string(hosts), err
}

func switchEnv(switchHostIP string) error {
	hosts, err := readHostsFile()
	if err != nil {
		return fmt.Errorf("Error opening hosts file: %w", err)
	}

	lines := strings.Split(hosts, "\n")
	var hasSwitchHost = false
	for i, line := range lines {
		if strings.HasPrefix(line, "# "+switchHostIP) {
			lines[i] = strings.TrimPrefix(line, "# ")
			hasSwitchHost = true
			fmt.Println("Successfully switching to local")
			break
		} else if strings.HasPrefix(line, switchHostIP) {
			lines[i] = "# " + line
			hasSwitchHost = true
			fmt.Println("Successfully switching to dev")
			break
		}
	}

	if !hasSwitchHost {
		return fmt.Errorf("Error cannot find SWITCH_HOST_IP in hosts file: %s", switchHostIP)
	}

	return ioutil.WriteFile("/private/etc/hosts", []byte(strings.Join(lines, "\n")), 0644)
}

func getStatus(switchHostIP string) (string, error) {
	hosts, err := readHostsFile()
	if err != nil {
		return "", fmt.Errorf("Error reading hosts file: %w", err)
	}

	lines := strings.Split(hosts, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "# "+switchHostIP) {
			return "dev", nil
		} else if strings.HasPrefix(line, switchHostIP) {
			return "local", nil
		}
	}

	return "", fmt.Errorf("Error cannot find SWITCH_HOST_IP in hosts file: %s", switchHostIP)
}
