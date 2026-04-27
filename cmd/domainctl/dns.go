package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func dnsMenu(reader *bufio.Reader, cfg *Config) {
	items := []string{
		"1) List zone records",
		"2) Add A record",
		"3) Add CNAME record",
		"4) Delete record",
		"5) Reload BIND9",
		"6) Check BIND9 config",
		"7) Service status",
		"0) Back",
	}
	for {
		clearScreen()
		drawBox("DNS: BIND9", items)
		switch prompt(reader, "Select") {
		case "1":
			dnsListRecords(cfg)
		case "2":
			dnsAddA(reader, cfg)
		case "3":
			dnsAddCNAME(reader, cfg)
		case "4":
			dnsDeleteRecord(reader, cfg)
		case "5":
			dnsReload()
		case "6":
			dnsCheck()
		case "7":
			dnsServiceStatus()
		case "0":
			return
		}
		pause(reader)
	}
}

// dnsNsupdate sends nsupdate commands to the local DNS server.
// Requires BIND9 zone configured with: allow-update { 127.0.0.1; };
func dnsNsupdate(input string) error {
	cmd := exec.Command("nsupdate")
	cmd.Stdin = bytes.NewBufferString(input)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("nsupdate: %s", string(out))
	}
	return nil
}

func dnsListRecords(cfg *Config) {
	out, err := runCapture("dig", "@127.0.0.1", cfg.Domain, "AXFR")
	if err != nil {
		out2, _ := runCapture("host", "-l", cfg.Domain, "127.0.0.1")
		fmt.Println(out2)
		return
	}
	fmt.Println(out)
}

func dnsAddA(reader *bufio.Reader, cfg *Config) {
	name := prompt(reader, "Hostname (short, e.g. ws1)")
	ip := prompt(reader, "IPv4 address")
	ttl := prompt(reader, "TTL [300]")
	if ttl == "" {
		ttl = "300"
	}
	fqdn := name + "." + cfg.Domain + "."
	input := fmt.Sprintf("zone %s\nupdate add %s %s A %s\nsend\n", cfg.Domain, fqdn, ttl, ip)
	if err := dnsNsupdate(input); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("A record %s -> %s added.\n", fqdn, ip)
}

func dnsAddCNAME(reader *bufio.Reader, cfg *Config) {
	alias := prompt(reader, "Alias (short, e.g. www)")
	target := prompt(reader, "Target FQDN (with trailing dot, e.g. dc1.corp.local.)")
	ttl := prompt(reader, "TTL [300]")
	if ttl == "" {
		ttl = "300"
	}
	fqdn := alias + "." + cfg.Domain + "."
	input := fmt.Sprintf("zone %s\nupdate add %s %s CNAME %s\nsend\n", cfg.Domain, fqdn, ttl, target)
	if err := dnsNsupdate(input); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("CNAME %s -> %s added.\n", fqdn, target)
}

func dnsDeleteRecord(reader *bufio.Reader, cfg *Config) {
	name := prompt(reader, "Hostname (short)")
	rtype := prompt(reader, "Record type (A/CNAME/MX/...)")
	fqdn := name + "." + cfg.Domain + "."
	input := fmt.Sprintf("zone %s\nupdate delete %s %s\nsend\n",
		cfg.Domain, fqdn, strings.ToUpper(rtype))
	if err := dnsNsupdate(input); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Record %s %s deleted.\n", fqdn, strings.ToUpper(rtype))
}

func dnsReload() {
	out, err := runCapture("rndc", "reload")
	if err != nil {
		fmt.Printf("Error: %s\n", out)
		return
	}
	fmt.Println("BIND9 reloaded.")
}

func dnsCheck() {
	out, err := runCapture("named-checkconf", "-z")
	if err != nil {
		fmt.Printf("Config errors:\n%s\n", out)
		return
	}
	fmt.Println(out)
	fmt.Println("Config OK.")
}

func dnsServiceStatus() {
	out, _ := runCapture("systemctl", "status", "bind9", "--no-pager")
	fmt.Println(out)
}
