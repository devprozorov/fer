package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const chronyConf = "/etc/chrony/chrony.conf"

func ntpMenu(reader *bufio.Reader, _ *Config) {
	items := []string{
		"1) Show tracking status",
		"2) Show sources",
		"3) Add NTP server/pool",
		"4) Remove NTP server/pool",
		"5) Restart Chrony",
		"0) Back",
	}
	for {
		clearScreen()
		drawBox("NTP: Chrony", items)
		switch prompt(reader, "Select") {
		case "1":
			ntpTracking()
		case "2":
			ntpSources()
		case "3":
			ntpAddServer(reader)
		case "4":
			ntpRemoveServer(reader)
		case "5":
			restartService("chrony")
		case "0":
			return
		}
		pause(reader)
	}
}

func ntpTracking() {
	out, err := runCapture("chronyc", "tracking")
	if err != nil {
		fmt.Printf("Error: %s\n", out)
		return
	}
	fmt.Println(out)
}

func ntpSources() {
	out, err := runCapture("chronyc", "sources", "-v")
	if err != nil {
		fmt.Printf("Error: %s\n", out)
		return
	}
	fmt.Println(out)
}

func ntpAddServer(reader *bufio.Reader) {
	server := prompt(reader, "NTP server hostname/IP")
	kind := prompt(reader, "Type: server or pool [server]")
	if kind == "" {
		kind = "server"
	}
	f, err := os.OpenFile(chronyConf, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening %s: %v\n", chronyConf, err)
		return
	}
	defer f.Close()
	_, _ = f.WriteString(fmt.Sprintf("%s %s iburst\n", kind, server))
	restartService("chrony")
	fmt.Printf("%s %q added.\n", kind, server)
}

func ntpRemoveServer(reader *bufio.Reader) {
	server := prompt(reader, "Server/pool to remove")
	data, err := os.ReadFile(chronyConf)
	if err != nil {
		fmt.Printf("Read error: %v\n", err)
		return
	}
	lines := strings.Split(string(data), "\n")
	var out []string
	for _, l := range lines {
		if strings.Contains(l, " "+server) {
			continue
		}
		out = append(out, l)
	}
	if err := os.WriteFile(chronyConf, []byte(strings.Join(out, "\n")), 0644); err != nil {
		fmt.Printf("Write error: %v\n", err)
		return
	}
	restartService("chrony")
	fmt.Printf("Entry for %q removed.\n", server)
}
