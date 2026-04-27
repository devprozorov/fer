# ACDC: Debian Domain Controller (GNOME)

ACDC is a Windows Server-like management stack for Debian with:

- OpenLDAP for directory
- Kerberos (MIT) for auth
- Samba shares
- BIND9 DNS
- ISC DHCP
- Chrony NTP
- SELinux policy constructor
- GNOME "GPO-like" policies via dconf + Ansible packages
- Centralized updates server (APT cache)
- Git-based self-update without losing state
- Terminal pseudo-graphics admin UI (`domainctl`) in Go

## Project Layout

- `ansible/` - infrastructure and service configuration
- `scripts/` - bootstrap, self-update, backups, SELinux constructor
- `cmd/domainctl/` - terminal pseudo-graphics domain management tool
- `systemd/` - timer/service for automatic self-updates
- `policies/` - sample GNOME and SELinux policy definitions

## Quick Start

1. Prepare Debian host and clone repository.
2. Run bootstrap:

```bash
sudo bash scripts/bootstrap.sh
```

3. Edit inventory and variables:

- `ansible/inventories/prod/hosts.ini`
- `ansible/group_vars/all.yml`

4. Apply full stack:

```bash
cd ansible
ansible-playbook -i inventories/prod/hosts.ini site.yml
```

5. Build and install terminal UI:

```bash
go build -o /usr/local/bin/domainctl ./cmd/domainctl
chmod +x /usr/local/bin/domainctl
```

6. Run pseudo-graphics management UI:

```bash
sudo domainctl
```

## Git Self-Update Without Data Loss

`scripts/self-update.sh` performs:

- pre-update backup (`/var/backups/acdc`)
- git fetch + fast-forward pull to tracked branch
- apply Ansible playbook
- rollback path for state files

Install timer:

```bash
sudo cp systemd/acdc-self-update.service /etc/systemd/system/
sudo cp systemd/acdc-self-update.timer /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable --now acdc-self-update.timer
```

## Security Notes

- Change all example passwords/secrets in `ansible/group_vars/all.yml`.
- Restrict SSH and management networks before production use.
- Validate SELinux policy modules before loading to production.
