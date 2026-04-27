#!/usr/bin/env bash
set -euo pipefail

REPO_DIR="/opt/acdc"
BRANCH="${ACDC_BRANCH:-main}"

if [[ ! -d "${REPO_DIR}/.git" ]]; then
  echo "Repository not found at ${REPO_DIR}"
  exit 1
fi

bash "${REPO_DIR}/scripts/backup.sh"

pushd "${REPO_DIR}" >/dev/null
git fetch --all --prune
git checkout "${BRANCH}"
git pull --ff-only origin "${BRANCH}"

cd ansible
ansible-playbook -i inventories/prod/hosts.ini site.yml
popd >/dev/null

echo "Self-update complete"
