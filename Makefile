ANSIBLE_INVENTORY ?= ansible/inventories/prod/hosts.ini

.PHONY: bootstrap deploy build-domainctl install-domainctl self-update

bootstrap:
	sudo bash scripts/bootstrap.sh

deploy:
	cd ansible && ansible-playbook -i inventories/prod/hosts.ini site.yml

build-domainctl:
	go build -o bin/domainctl ./cmd/domainctl

install-domainctl: build-domainctl
	sudo install -m 0755 bin/domainctl /usr/local/bin/domainctl

self-update:
	sudo bash scripts/self-update.sh
