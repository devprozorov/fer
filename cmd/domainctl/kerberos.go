package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
)

func kerberosMenu(reader *bufio.Reader, _ *Config) {
	items := []string{
		"1) List principals",
		"2) Add principal",
		"3) Delete principal",
		"4) Change principal password",
		"5) Show realm info",
		"6) Service status",
		"0) Back",
	}
	for {
		clearScreen()
		drawBox("Kerberos Principals", items)
		switch prompt(reader, "Select") {
		case "1":
			kerbListPrincipals()
		case "2":
			kerbAddPrincipal(reader)
		case "3":
			kerbDeletePrincipal(reader)
		case "4":
			kerbChangePassword(reader)
		case "5":
			kerbRealmInfo()
		case "6":
			kerbStatus()
		case "0":
			return
		}
		pause(reader)
	}
}

func kadminLocal(query string) (string, error) {
	cmd := exec.Command("kadmin.local", "-q", query)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func kerbListPrincipals() {
	out, err := kadminLocal("listprincs")
	if err != nil {
		fmt.Printf("Error: %s\n", out)
		return
	}
	fmt.Println(out)
}

func kerbAddPrincipal(reader *bufio.Reader) {
	principal := prompt(reader, "Principal name (e.g. alice or alice@CORP.LOCAL)")
	password := prompt(reader, "Initial password")
	out, err := kadminLocal(fmt.Sprintf("addprinc -pw %s %s", password, principal))
	if err != nil {
		fmt.Printf("Error: %s\n", out)
		return
	}
	fmt.Printf("Principal %q added.\n", principal)
}

func kerbDeletePrincipal(reader *bufio.Reader) {
	principal := prompt(reader, "Principal to delete")
	out, err := kadminLocal("delprinc -force " + principal)
	if err != nil {
		fmt.Printf("Error: %s\n", out)
		return
	}
	fmt.Printf("Principal %q deleted.\n", principal)
}

func kerbChangePassword(reader *bufio.Reader) {
	principal := prompt(reader, "Principal name")
	password := prompt(reader, "New password")
	out, err := kadminLocal(fmt.Sprintf("cpw -pw %s %s", password, principal))
	if err != nil {
		fmt.Printf("Error: %s\n", out)
		return
	}
	fmt.Printf("Password for %q changed.\n", principal)
}

func kerbRealmInfo() {
	cmd := exec.Command("kadmin.local")
	cmd.Stdin = bytes.NewBufferString("getdefaultprinc\nquit\n")
	out, _ := cmd.CombinedOutput()
	fmt.Println(string(out))
}

func kerbStatus() {
	out, _ := runCapture("systemctl", "status", "krb5-kdc", "--no-pager")
	fmt.Println(out)
	out2, _ := runCapture("systemctl", "status", "krb5-admin-server", "--no-pager")
	fmt.Println(out2)
}
