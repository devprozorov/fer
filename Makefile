ANSIBLE_INVENTORY ?= ansible/inventories/prod/hosts.ini

.PHONY: bootstrap deploy build-fer install-fer build-domainctl install-domainctl self-update

bootstrap:
	sudo bash scripts/bootstrap.sh

deploy:
	cd ansible && ansible-playbook -i inventories/prod/hosts.ini site.yml

build-fer:
	go build -o bin/fer ./cmd/domainctl

install-fer: build-fer
	sudo install -m 0755 bin/fer /usr/local/bin/fer

# Backward-compatible aliases
build-domainctl: build-fer

install-domainctl: install-fer

self-update:
	sudo bash scripts/self-update.sh
