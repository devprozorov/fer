package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func sambaMenu(reader *bufio.Reader, cfg *Config) {
	items := []string{
		"1) List ACDC-managed shares",
		"2) Add share",
		"3) Remove share",
		"4) Set user Samba password",
		"5) Delete Samba user",
		"6) Service status",
		"7) Restart Samba",
		"0) Back",
	}
	for {
		clearScreen()
		drawBox("Samba Shares", items)
		switch prompt(reader, "Select") {
		case "1":
			sambaListShares(cfg)
		case "2":
			sambaAddShare(reader, cfg)
		case "3":
			sambaRemoveShare(reader, cfg)
		case "4":
			sambaSetPassword(reader)
		case "5":
			sambaDeleteUser(reader)
		case "6":
			sambaStatus()
		case "7":
			restartService("smbd")
		case "0":
			return
		}
		pause(reader)
	}
}

func sambaListShares(cfg *Config) {
	data, err := os.ReadFile(cfg.SambaACDCConf)
	if err != nil {
		fmt.Println("No ACDC-managed shares found.")
		return
	}
	fmt.Println(string(data))
}

func sambaAddShare(reader *bufio.Reader, cfg *Config) {
	name := prompt(reader, "Share name")
	path := prompt(reader, "Path (e.g. /srv/samba/docs)")
	writable := prompt(reader, "Writable? (yes/no)")
	guestOk := prompt(reader, "Guest OK? (yes/no)")

	if err := os.MkdirAll(path, 0770); err != nil {
		fmt.Printf("mkdir error: %v\n", err)
		return
	}
	block := fmt.Sprintf("\n[%s]\n   path = %s\n   browsable = yes\n   writable = %s\n   guest ok = %s\n   create mask = 0660\n   directory mask = 0770\n\n",
		name, path, writable, guestOk)

	if err := os.MkdirAll("/var/lib/acdc/samba", 0750); err != nil {
		fmt.Printf("mkdir error: %v\n", err)
		return
	}
	f, err := os.OpenFile(cfg.SambaACDCConf, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Open error: %v\n", err)
		return
	}
	defer f.Close()
	if _, err := f.WriteString(block); err != nil {
		fmt.Printf("Write error: %v\n", err)
		return
	}
	sambaEnsureInclude(cfg)
	restartService("smbd")
	fmt.Printf("Share [%s] added.\n", name)
}

func sambaRemoveShare(reader *bufio.Reader, cfg *Config) {
	name := prompt(reader, "Share name to remove")
	data, err := os.ReadFile(cfg.SambaACDCConf)
	if err != nil {
		fmt.Printf("Read error: %v\n", err)
		return
	}
	lines := strings.Split(string(data), "\n")
	var out []string
	skip := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "["+name+"]" {
			skip = true
			continue
		}
		if skip && strings.HasPrefix(trimmed, "[") {
			skip = false
		}
		if !skip {
			out = append(out, line)
		}
	}
	if err := os.WriteFile(cfg.SambaACDCConf, []byte(strings.Join(out, "\n")), 0644); err != nil {
		fmt.Printf("Write error: %v\n", err)
		return
	}
	restartService("smbd")
	fmt.Printf("Share [%s] removed.\n", name)
}

func sambaSetPassword(reader *bufio.Reader) {
	user := prompt(reader, "UNIX username")
	pass := prompt(reader, "Samba password")
	// smbpasswd -s reads two password lines (new + confirm) from stdin
	cmd := exec.Command("smbpasswd", "-a", "-s", user)
	cmd.Stdin = strings.NewReader(pass + "\n" + pass + "\n")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error: %s\n", string(out))
		return
	}
	fmt.Printf("Samba password set for %q.\n", user)
}

func sambaDeleteUser(reader *bufio.Reader) {
	user := prompt(reader, "Username to remove")
	out, err := runCapture("smbpasswd", "-x", user)
	if err != nil {
		fmt.Printf("Error: %s\n", out)
		return
	}
	fmt.Printf("Samba credentials for %q removed.\n", user)
}

func sambaStatus() {
	out, _ := runCapture("systemctl", "status", "smbd", "--no-pager", "-l")
	fmt.Println(out)
}

// sambaEnsureInclude appends include directive to smb.conf if not present.
func sambaEnsureInclude(cfg *Config) {
	data, err := os.ReadFile(cfg.SambaConf)
	if err != nil {
		return
	}
	include := "include = " + cfg.SambaACDCConf
	if strings.Contains(string(data), include) {
		return
	}
	f, err := os.OpenFile(cfg.SambaConf, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	_, _ = f.WriteString("\n" + include + "\n")
}
