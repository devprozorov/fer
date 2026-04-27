package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func updatesMenu(reader *bufio.Reader, _ *Config) {
	items := []string{
		"1) Show apt-cacher-ng status",
		"2) Restart apt-cacher-ng",
		"3) Show cache log (last 30 lines)",
		"4) Show proxy config",
		"5) Run apt-get update + upgrade (local)",
		"0) Back",
	}
	for {
		clearScreen()
		drawBox("Updates Server (apt-cacher-ng)", items)
		switch prompt(reader, "Select") {
		case "1":
			updatesStatus()
		case "2":
			restartService("apt-cacher-ng")
		case "3":
			updatesCacheLog()
		case "4":
			updatesShowProxy()
		case "5":
			updatesRunLocal()
		case "0":
			return
		}
		pause(reader)
	}
}

func updatesStatus() {
	out, _ := runCapture("systemctl", "status", "apt-cacher-ng", "--no-pager")
	fmt.Println(out)
}

func updatesCacheLog() {
	data, err := os.ReadFile("/var/log/apt-cacher-ng/apt-cacher.log")
	if err != nil {
		fmt.Printf("Cannot read log: %v\n", err)
		return
	}
	lines := strings.Split(string(data), "\n")
	start := 0
	if len(lines) > 30 {
		start = len(lines) - 30
	}
	for _, l := range lines[start:] {
		fmt.Println(l)
	}
}

func updatesShowProxy() {
	data, err := os.ReadFile("/etc/apt/apt.conf.d/01proxy")
	if err != nil {
		fmt.Printf("No proxy config found: %v\n", err)
		return
	}
	fmt.Println(string(data))
}

func updatesRunLocal() {
	fmt.Println("Running apt-get update...")
	out, err := runCapture("apt-get", "update")
	fmt.Println(out)
	if err != nil {
		fmt.Printf("apt-get update failed\n")
		return
	}
	fmt.Println("Running apt-get upgrade -y...")
	out2, _ := runCapture("apt-get", "upgrade", "-y")
	fmt.Println(out2)
}
