package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

func selinuxMenu(reader *bufio.Reader, cfg *Config) {
	items := []string{
		"1) List installed modules",
		"2) Build & install new policy",
		"3) Remove module",
		"4) SELinux status",
		"5) Set Enforcing",
		"6) Set Permissive",
		"0) Back",
	}
	for {
		clearScreen()
		drawBox("SELinux Policies", items)
		switch prompt(reader, "Select") {
		case "1":
			selinuxListModules()
		case "2":
			selinuxBuildPolicy(reader, cfg)
		case "3":
			selinuxRemoveModule(reader)
		case "4":
			selinuxStatus()
		case "5":
			selinuxSetMode("1")
		case "6":
			selinuxSetMode("0")
		case "0":
			return
		}
		pause(reader)
	}
}

func selinuxListModules() {
	out, err := runCapture("semodule", "-l")
	if err != nil {
		fmt.Printf("Error: %s\n", out)
		return
	}
	fmt.Println(out)
}

// selinuxBuildPolicy is the interactive policy constructor.
// It generates .te + .fc files, compiles them and installs the module.
func selinuxBuildPolicy(reader *bufio.Reader, cfg *Config) {
	name := prompt(reader, "Policy module name (e.g. myapp)")
	domain := prompt(reader, "Domain type (e.g. myapp_t)")
	netBind := prompt(reader, "Allow network bind? (yes/no)")
	fileRead := prompt(reader, "Path pattern to allow read (optional, e.g. /srv/app/*)")
	fileWrite := prompt(reader, "Path pattern to allow write (optional)")

	outDir := cfg.StateRoot + "/selinux/" + name
	if err := os.MkdirAll(outDir, 0750); err != nil {
		fmt.Printf("mkdir error: %v\n", err)
		return
	}

	te := fmt.Sprintf("policy_module(%s, 1.0)\n\nrequire {\n  type %s;\n}\n", name, domain)
	if netBind == "yes" {
		te += fmt.Sprintf("corenet_tcp_bind_generic_node(%s)\ncorenet_udp_bind_generic_node(%s)\n", domain, domain)
	}

	fc := ""
	if fileRead != "" {
		fc += fmt.Sprintf("%s  --  gen_context(system_u:object_r:var_t,s0)\n", fileRead)
		te += fmt.Sprintf("allow %s var_t:file { read open getattr };\n", domain)
	}
	if fileWrite != "" {
		fc += fmt.Sprintf("%s  --  gen_context(system_u:object_r:var_t,s0)\n", fileWrite)
		te += fmt.Sprintf("allow %s var_t:file { write append open getattr };\n", domain)
	}

	if err := os.WriteFile(outDir+"/"+name+".te", []byte(te), 0640); err != nil {
		fmt.Printf("Write .te error: %v\n", err)
		return
	}
	if err := os.WriteFile(outDir+"/"+name+".fc", []byte(fc), 0640); err != nil {
		fmt.Printf("Write .fc error: %v\n", err)
		return
	}

	fmt.Printf("Compiling policy module %q...\n", name)
	cmd := exec.Command("make", "-f", "/usr/share/selinux/devel/Makefile", "-C", outDir, name+".pp")
	out, err := cmd.CombinedOutput()
	fmt.Print(string(out))
	if err != nil {
		fmt.Printf("make error: %v\n", err)
		return
	}
	out2, err2 := runCapture("semodule", "-i", outDir+"/"+name+".pp")
	if err2 != nil {
		fmt.Printf("semodule error: %s\n", out2)
		return
	}
	fmt.Printf("Policy module %q installed.\n", name)
}

func selinuxRemoveModule(reader *bufio.Reader) {
	name := prompt(reader, "Module name to remove")
	out, err := runCapture("semodule", "-r", name)
	if err != nil {
		fmt.Printf("Error: %s\n", out)
		return
	}
	fmt.Printf("Module %q removed.\n", name)
}

func selinuxStatus() {
	out, _ := runCapture("sestatus")
	fmt.Println(out)
}

func selinuxSetMode(enforcing string) {
	out, err := runCapture("setenforce", enforcing)
	if err != nil {
		fmt.Printf("setenforce: %s\n", out)
		return
	}
	if enforcing == "1" {
		fmt.Println("SELinux set to Enforcing.")
	} else {
		fmt.Println("SELinux set to Permissive.")
	}
}
