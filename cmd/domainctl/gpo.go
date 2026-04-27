package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	dconfDBDir    = "/etc/dconf/db/local.d"
	dconfLocksDir = "/etc/dconf/db/local.d/locks"
)

func gpoMenu(reader *bufio.Reader, _ *Config) {
	items := []string{
		"1) List all policies",
		"2) Apply dconf setting",
		"3) Lock key (prevent user override)",
		"4) Remove dconf setting",
		"5) Install packages via Ansible",
		"6) Run Ansible playbook",
		"7) Run dconf update",
		"0) Back",
	}
	for {
		clearScreen()
		drawBox("GNOME GPO Policies", items)
		switch prompt(reader, "Select") {
		case "1":
			gpoList()
		case "2":
			gpoSetDconf(reader)
		case "3":
			gpoLockDconf(reader)
		case "4":
			gpoRemoveSetting(reader)
		case "5":
			gpoInstallPackages(reader)
		case "6":
			gpoRunPlaybook(reader)
		case "7":
			gpoUpdateDconf()
		case "0":
			return
		}
		pause(reader)
	}
}

func gpoList() {
	entries, err := os.ReadDir(dconfDBDir)
	if err != nil {
		fmt.Printf("Read error: %v\n", err)
		return
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		fmt.Printf("\n=== %s ===\n", e.Name())
		data, rerr := os.ReadFile(dconfDBDir + "/" + e.Name())
		if rerr == nil {
			fmt.Println(string(data))
		}
	}
}

// gpoSetDconf adds or replaces a key in a named dconf policy file.
func gpoSetDconf(reader *bufio.Reader) {
	policyFile := prompt(reader, "Policy file name (e.g. 10-screensaver)")
	section := prompt(reader, "dconf section (e.g. org/gnome/desktop/screensaver)")
	key := prompt(reader, "Key (e.g. lock-enabled)")
	value := prompt(reader, "Value (e.g. true or uint32 900)")

	path := dconfDBDir + "/" + policyFile
	existing, _ := os.ReadFile(path)
	content := string(existing)
	sectionHeader := "[" + section + "]"
	entry := key + "=" + value + "\n"

	if strings.Contains(content, sectionHeader) {
		parts := strings.SplitN(content, sectionHeader, 2)
		content = parts[0] + sectionHeader + "\n" + entry + parts[1]
	} else {
		content += "\n" + sectionHeader + "\n" + entry
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		fmt.Printf("Write error: %v\n", err)
		return
	}
	gpoUpdateDconf()
	fmt.Printf("Set [%s] %s = %s\n", section, key, value)
}

// gpoLockDconf appends a key path to a locks file so users cannot override it.
func gpoLockDconf(reader *bufio.Reader) {
	locksFile := prompt(reader, "Locks file name (e.g. 10-locks)")
	keyPath := prompt(reader, "Full key path (e.g. /org/gnome/desktop/screensaver/lock-enabled)")

	if err := os.MkdirAll(dconfLocksDir, 0755); err != nil {
		fmt.Printf("mkdir error: %v\n", err)
		return
	}
	f, err := os.OpenFile(dconfLocksDir+"/"+locksFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Open error: %v\n", err)
		return
	}
	defer f.Close()
	_, _ = f.WriteString(keyPath + "\n")
	gpoUpdateDconf()
	fmt.Printf("Lock on %s added.\n", keyPath)
}

func gpoRemoveSetting(reader *bufio.Reader) {
	policyFile := prompt(reader, "Policy file name to edit")
	key := prompt(reader, "Key to remove")

	path := dconfDBDir + "/" + policyFile
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Read error: %v\n", err)
		return
	}
	lines := strings.Split(string(data), "\n")
	var out []string
	for _, l := range lines {
		if !strings.HasPrefix(strings.TrimSpace(l), key+"=") {
			out = append(out, l)
		}
	}
	if err := os.WriteFile(path, []byte(strings.Join(out, "\n")), 0644); err != nil {
		fmt.Printf("Write error: %v\n", err)
		return
	}
	gpoUpdateDconf()
	fmt.Printf("Key %q removed.\n", key)
}

func gpoInstallPackages(reader *bufio.Reader) {
	pkgs := prompt(reader, "Package names (space-separated)")
	args := []string{
		"all",
		"-i", "/opt/acdc/ansible/inventories/prod/hosts.ini",
		"-m", "apt",
		"-a", "name=" + pkgs + " state=present",
		"--become",
	}
	out, err := runCapture("ansible", args...)
	if err != nil {
		fmt.Printf("Ansible error: %s\n", out)
		return
	}
	fmt.Println(out)
}

func gpoRunPlaybook(reader *bufio.Reader) {
	playbook := prompt(reader, "Playbook path [/opt/acdc/ansible/site.yml]")
	if playbook == "" {
		playbook = "/opt/acdc/ansible/site.yml"
	}
	out, err := runCapture("ansible-playbook",
		"-i", "/opt/acdc/ansible/inventories/prod/hosts.ini", playbook)
	if err != nil {
		fmt.Printf("ansible-playbook error: %s\n", out)
		return
	}
	fmt.Println(out)
}

func gpoUpdateDconf() {
	out, err := runCapture("dconf", "update")
	if err != nil {
		fmt.Printf("dconf update error: %s\n", out)
		return
	}
	fmt.Println("dconf database updated.")
}
