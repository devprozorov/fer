package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
)

func ldapMenu(reader *bufio.Reader, cfg *Config) {
	items := []string{
		"1) Create user",
		"2) Edit user (displayName)",
		"3) Delete user",
		"4) List users",
		"5) Create group",
		"6) Add member to group",
		"7) Remove member from group",
		"8) Delete group",
		"9) List groups",
		"0) Back",
	}
	for {
		clearScreen()
		drawBox("LDAP: Users & Groups", items)
		switch prompt(reader, "Select") {
		case "1":
			ldapCreateUser(reader, cfg)
		case "2":
			ldapEditUser(reader, cfg)
		case "3":
			ldapDeleteUser(reader, cfg)
		case "4":
			ldapListUsers(cfg)
		case "5":
			ldapCreateGroup(reader, cfg)
		case "6":
			ldapAddMember(reader, cfg)
		case "7":
			ldapRemoveMember(reader, cfg)
		case "8":
			ldapDeleteGroup(reader, cfg)
		case "9":
			ldapListGroups(cfg)
		case "0":
			return
		}
		pause(reader)
	}
}

// ldapExec runs ldapadd or ldapmodify piping ldif to stdin.
func ldapExec(toolName string, cfg *Config, ldif string) error {
	cmd := exec.Command(toolName, "-x", "-D", cfg.LDAPAdminDN, "-w", cfg.LDAPPassword)
	cmd.Stdin = bytes.NewBufferString(ldif)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", toolName, string(out))
	}
	return nil
}

func ldapDel(cfg *Config, dn string) error {
	cmd := exec.Command("ldapdelete", "-x", "-D", cfg.LDAPAdminDN, "-w", cfg.LDAPPassword, dn)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ldapdelete: %s", string(out))
	}
	return nil
}

func ldapSearch(cfg *Config, base, filter string, attrs ...string) (string, error) {
	args := []string{"-x", "-D", cfg.LDAPAdminDN, "-w", cfg.LDAPPassword, "-b", base, filter}
	args = append(args, attrs...)
	cmd := exec.Command("ldapsearch", args...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func ldapCreateUser(reader *bufio.Reader, cfg *Config) {
	uid := prompt(reader, "UID (login)")
	cn := prompt(reader, "Common Name")
	sn := prompt(reader, "Surname")
	given := prompt(reader, "Given Name")
	uidNum := prompt(reader, "UID Number (e.g. 20001)")
	gidNum := prompt(reader, "GID Number (e.g. 20000)")
	ldif := fmt.Sprintf("dn: uid=%s,ou=People,%s\nobjectClass: inetOrgPerson\nobjectClass: posixAccount\nobjectClass: shadowAccount\nuid: %s\ncn: %s\nsn: %s\ngivenName: %s\nuidNumber: %s\ngidNumber: %s\nhomeDirectory: /home/%s\nloginShell: /bin/bash\n",
		uid, cfg.LDAPBaseDN, uid, cn, sn, given, uidNum, gidNum, uid)
	if err := ldapExec("ldapadd", cfg, ldif); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("User %q created.\n", uid)
}

func ldapEditUser(reader *bufio.Reader, cfg *Config) {
	uid := prompt(reader, "UID")
	displayName := prompt(reader, "New displayName")
	ldif := fmt.Sprintf("dn: uid=%s,ou=People,%s\nchangetype: modify\nreplace: displayName\ndisplayName: %s\n",
		uid, cfg.LDAPBaseDN, displayName)
	if err := ldapExec("ldapmodify", cfg, ldif); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println("User updated.")
}

func ldapDeleteUser(reader *bufio.Reader, cfg *Config) {
	uid := prompt(reader, "UID to delete")
	if err := ldapDel(cfg, fmt.Sprintf("uid=%s,ou=People,%s", uid, cfg.LDAPBaseDN)); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("User %q deleted.\n", uid)
}

func ldapListUsers(cfg *Config) {
	out, err := ldapSearch(cfg, "ou=People,"+cfg.LDAPBaseDN, "(objectClass=inetOrgPerson)", "uid", "cn")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println(out)
}

func ldapCreateGroup(reader *bufio.Reader, cfg *Config) {
	cn := prompt(reader, "Group CN")
	gid := prompt(reader, "GID Number")
	ldif := fmt.Sprintf("dn: cn=%s,ou=Groups,%s\nobjectClass: posixGroup\ncn: %s\ngidNumber: %s\n",
		cn, cfg.LDAPBaseDN, cn, gid)
	if err := ldapExec("ldapadd", cfg, ldif); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Group %q created.\n", cn)
}

func ldapAddMember(reader *bufio.Reader, cfg *Config) {
	group := prompt(reader, "Group CN")
	uid := prompt(reader, "UID to add")
	ldif := fmt.Sprintf("dn: cn=%s,ou=Groups,%s\nchangetype: modify\nadd: memberUid\nmemberUid: %s\n",
		group, cfg.LDAPBaseDN, uid)
	if err := ldapExec("ldapmodify", cfg, ldif); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Member %q added to group %q.\n", uid, group)
}

func ldapRemoveMember(reader *bufio.Reader, cfg *Config) {
	group := prompt(reader, "Group CN")
	uid := prompt(reader, "UID to remove")
	ldif := fmt.Sprintf("dn: cn=%s,ou=Groups,%s\nchangetype: modify\ndelete: memberUid\nmemberUid: %s\n",
		group, cfg.LDAPBaseDN, uid)
	if err := ldapExec("ldapmodify", cfg, ldif); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Member %q removed from group %q.\n", uid, group)
}

func ldapDeleteGroup(reader *bufio.Reader, cfg *Config) {
	cn := prompt(reader, "Group CN to delete")
	if err := ldapDel(cfg, fmt.Sprintf("cn=%s,ou=Groups,%s", cn, cfg.LDAPBaseDN)); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Group %q deleted.\n", cn)
}

func ldapListGroups(cfg *Config) {
	out, err := ldapSearch(cfg, "ou=Groups,"+cfg.LDAPBaseDN, "(objectClass=posixGroup)",
		"cn", "gidNumber", "memberUid")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println(out)
}
