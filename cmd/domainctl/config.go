package main

import (
	"os"
	"strings"
)

// Config holds runtime configuration sourced from environment variables.
// All default values match ansible/group_vars/all.yml.
type Config struct {
	LDAPAdminDN   string
	LDAPPassword  string
	LDAPBaseDN    string
	Domain        string
	SambaConf     string
	SambaACDCConf string
	DHCPConfig    string
	DHCPACDCConf  string
	StateRoot     string
}

func newConfig() *Config {
	return &Config{
		LDAPAdminDN:   getenv("LDAP_ADMIN_DN", "cn=admin,dc=corp,dc=local"),
		LDAPPassword:  getenv("LDAP_ADMIN_PASSWORD", "ChangeMeLDAP"),
		LDAPBaseDN:    getenv("LDAP_BASE_DN", "dc=corp,dc=local"),
		Domain:        getenv("DOMAIN_NAME", "corp.local"),
		SambaConf:     "/etc/samba/smb.conf",
		SambaACDCConf: "/var/lib/acdc/samba/shares.conf",
		DHCPConfig:    "/etc/dhcp/dhcpd.conf",
		DHCPACDCConf:  "/var/lib/acdc/dhcp/static.conf",
		StateRoot:     "/var/lib/acdc",
	}
}

func getenv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && strings.TrimSpace(v) != "" {
		return v
	}
	return fallback
}
