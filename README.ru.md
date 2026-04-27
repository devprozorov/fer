# fer: Контроллер Домена Debian (GNOME)

fer — это стек управления для Debian, похожий по возможностям на Windows Server, с использованием:

- OpenLDAP для каталога
- Kerberos (MIT) для аутентификации
- Samba для общих папок
- BIND9 для DNS
- ISC DHCP
- Chrony для NTP
- Конструктора политик SELinux
- GNOME-политик в стиле "GPO" через dconf + пакеты Ansible
- Централизованного сервера обновлений (APT cache)
- Самообновления из Git без потери состояния
- Терминального псевдографического UI администратора (`fer`) на Go

## Структура Проекта

- `ansible/` — инфраструктура и конфигурация сервисов
- `scripts/` — bootstrap, самообновление, бэкапы, конструктор SELinux
- `cmd/domainctl/` — терминальная утилита управления доменом с псевдографикой
- `systemd/` — таймер/сервис для автоматических самообновлений
- `policies/` — примеры GNOME- и SELinux-политик

## Быстрый Старт

1. Подготовьте Debian-хост и клонируйте репозиторий.
2. Запустите bootstrap:

```bash
sudo bash scripts/bootstrap.sh
```

3. Отредактируйте inventory и переменные:

- `ansible/inventories/prod/hosts.ini`
- `ansible/group_vars/all.yml`

4. Примените полный стек:

```bash
cd ansible
ansible-playbook -i inventories/prod/hosts.ini site.yml
```

5. Соберите и установите терминальный UI:

```bash
go build -o /usr/local/bin/fer ./cmd/domainctl
chmod +x /usr/local/bin/fer
```

6. Запустите псевдографический интерфейс управления:

```bash
sudo fer
```

## Самообновление Из Git Без Потери Данных

Скрипт `scripts/self-update.sh` выполняет:

- бэкап перед обновлением (`/var/backups/acdc`)
- `git fetch` + `fast-forward pull` для отслеживаемой ветки
- применение Ansible playbook
- сценарий отката для state-файлов

Установка таймера:

```bash
sudo cp systemd/acdc-self-update.service /etc/systemd/system/
sudo cp systemd/acdc-self-update.timer /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable --now acdc-self-update.timer
```

## Заметки По Безопасности

- Замените все примерные пароли/секреты в `ansible/group_vars/all.yml`.
- Ограничьте SSH и управляющие сети перед запуском в production.
- Проверяйте SELinux-модули на staging перед загрузкой в production.
