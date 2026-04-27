package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func dhcpMenu(reader *bufio.Reader, cfg *Config) {
	items := []string{
		"1) Show current leases",
		"2) Add static lease",
		"3) Remove static lease",
		"4) Service status",
		"5) Restart DHCP",
		"0) Back",
	}
	for {
		clearScreen()
		drawBox("DHCP: ISC DHCP Server", items)
		switch prompt(reader, "Select") {
		case "1":
			dhcpShowLeases()
		case "2":
			dhcpAddStatic(reader, cfg)
		case "3":
			dhcpRemoveStatic(reader, cfg)
		case "4":
			dhcpStatus()
		case "5":
			restartService("isc-dhcp-server")
		case "0":
			return
		}
		pause(reader)
	}
}

func dhcpShowLeases() {
	data, err := os.ReadFile("/var/lib/dhcp/dhcpd.leases")
	if err != nil {
		fmt.Printf("Cannot read leases file: %v\n", err)
		return
	}
	fmt.Println(string(data))
}

func dhcpAddStatic(reader *bufio.Reader, cfg *Config) {
	hostname := prompt(reader, "Hostname")
	mac := prompt(reader, "MAC address (xx:xx:xx:xx:xx:xx)")
	ip := prompt(reader, "Fixed IP address")

	block := fmt.Sprintf("\nhost %s {\n  hardware ethernet %s;\n  fixed-address %s;\n}\n",
		hostname, mac, ip)

	if err := os.MkdirAll("/var/lib/acdc/dhcp", 0750); err != nil {
		fmt.Printf("mkdir error: %v\n", err)
		return
	}
	f, err := os.OpenFile(cfg.DHCPACDCConf, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Open error: %v\n", err)
		return
	}
	defer f.Close()
	if _, err := f.WriteString(block); err != nil {
		fmt.Printf("Write error: %v\n", err)
		return
	}
	dhcpEnsureInclude(cfg)
	restartService("isc-dhcp-server")
	fmt.Printf("Static lease: %s (%s) -> %s added.\n", hostname, mac, ip)
}

func dhcpRemoveStatic(reader *bufio.Reader, cfg *Config) {
	hostname := prompt(reader, "Hostname to remove")
	data, err := os.ReadFile(cfg.DHCPACDCConf)
	if err != nil {
		fmt.Printf("Read error: %v\n", err)
		return
	}
	lines := strings.Split(string(data), "\n")
	var out []string
	skip := false
	depth := 0
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "host "+hostname+" {") || trimmed == "host "+hostname+"{" {
			skip = true
			depth = 1
			continue
		}
		if skip {
			depth += strings.Count(line, "{") - strings.Count(line, "}")
			if depth <= 0 {
				skip = false
			}
			continue
		}
		out = append(out, line)
	}
	if err := os.WriteFile(cfg.DHCPACDCConf, []byte(strings.Join(out, "\n")), 0644); err != nil {
		fmt.Printf("Write error: %v\n", err)
		return
	}
	restartService("isc-dhcp-server")
	fmt.Printf("Static lease for %q removed.\n", hostname)
}

func dhcpStatus() {
	out, _ := runCapture("systemctl", "status", "isc-dhcp-server", "--no-pager")
	fmt.Println(out)
}

// dhcpEnsureInclude appends include line to dhcpd.conf if not already present.
func dhcpEnsureInclude(cfg *Config) {
	data, err := os.ReadFile(cfg.DHCPConfig)
	if err != nil {
		return
	}
	include := `include "` + cfg.DHCPACDCConf + `";`
	if strings.Contains(string(data), include) {
		return
	}
	f, err := os.OpenFile(cfg.DHCPConfig, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	_, _ = f.WriteString("\n" + include + "\n")
}
