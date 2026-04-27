#!/usr/bin/env bash
set -euo pipefail

if [[ "${EUID}" -ne 0 ]]; then
  echo "Run as root"
  exit 1
fi

export DEBIAN_FRONTEND=noninteractive
apt-get update
apt-get install -y git ansible golang-go make rsync ldap-utils

mkdir -p /opt/acdc /var/lib/acdc /var/backups/acdc
echo "Bootstrap complete"
