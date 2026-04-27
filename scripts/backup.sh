#!/usr/bin/env bash
set -euo pipefail

STAMP="$(date +%Y%m%d-%H%M%S)"
DEST="/var/backups/acdc/${STAMP}"
mkdir -p "${DEST}"

slapcat -l "${DEST}/ldap.ldif" || true
cp -a /etc/krb5.conf "${DEST}/krb5.conf" 2>/dev/null || true
cp -a /etc/samba/smb.conf "${DEST}/smb.conf" 2>/dev/null || true
cp -a /etc/bind "${DEST}/bind" 2>/dev/null || true
cp -a /etc/dhcp/dhcpd.conf "${DEST}/dhcpd.conf" 2>/dev/null || true
cp -a /etc/chrony/chrony.conf "${DEST}/chrony.conf" 2>/dev/null || true

echo "Backup saved to ${DEST}"
