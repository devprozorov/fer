package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	cfg := newConfig()
	reader := bufio.NewReader(os.Stdin)
	for {
		clearScreen()
		drawBox("FER DOMAIN CONTROLLER", []string{
			"1) Users & Groups  (LDAP)",
			"2) Samba Shares",
			"3) DNS Records  (BIND9)",
			"4) DHCP Leases / Static",
			"5) NTP / Chrony",
			"6) SELinux Policies",
			"7) GNOME GPO Policies",
			"8) Kerberos Principals",
			"9) Updates Server",
			"0) Exit",
		})
		switch prompt(reader, "Select") {
		case "1":
			ldapMenu(reader, cfg)
		case "2":
			sambaMenu(reader, cfg)
		case "3":
			dnsMenu(reader, cfg)
		case "4":
			dhcpMenu(reader, cfg)
		case "5":
			ntpMenu(reader, cfg)
		case "6":
			selinuxMenu(reader, cfg)
		case "7":
			gpoMenu(reader, cfg)
		case "8":
			kerberosMenu(reader, cfg)
		case "9":
			updatesMenu(reader, cfg)
		case "0":
			fmt.Println("Bye")
			return
		default:
			fmt.Print("Unknown option. Press Enter...")
			_, _ = reader.ReadString('\n')
		}
	}
}
