#!/usr/bin/env bash
set -euo pipefail

read -rp "Policy name: " POLICY
read -rp "Process type (e.g. myapp_t): " DOMAIN_TYPE
read -rp "Path to allow read (optional): " READ_PATH
read -rp "Path to allow write (optional): " WRITE_PATH

OUT_DIR="/var/lib/acdc/selinux/${POLICY}"
mkdir -p "${OUT_DIR}"

cat > "${OUT_DIR}/${POLICY}.te" <<EOF
policy_module(${POLICY}, 1.0)

type ${DOMAIN_TYPE};
EOF

if [[ -n "${READ_PATH}" ]]; then
  cat >> "${OUT_DIR}/${POLICY}.te" <<EOF
allow ${DOMAIN_TYPE} var_t:file { read open getattr };
EOF
fi

if [[ -n "${WRITE_PATH}" ]]; then
  cat >> "${OUT_DIR}/${POLICY}.te" <<EOF
allow ${DOMAIN_TYPE} var_t:file { write append open getattr };
EOF
fi

cat > "${OUT_DIR}/${POLICY}.fc" <<EOF
${READ_PATH}    --      gen_context(system_u:object_r:var_t,s0)
${WRITE_PATH}   --      gen_context(system_u:object_r:var_t,s0)
EOF

make -f /usr/share/selinux/devel/Makefile -C "${OUT_DIR}" "${POLICY}.pp"
semodule -i "${OUT_DIR}/${POLICY}.pp"

echo "Installed SELinux module: ${POLICY}"
