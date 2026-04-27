#!/usr/bin/env bash
set -euo pipefail

read -rp "Policy name: " POLICY
read -rp "Domain type (example myapp_t): " DOMAIN
read -rp "Allow network bind? (yes/no): " NET_BIND
read -rp "Output directory [/var/lib/acdc/selinux]: " OUT_DIR

OUT_DIR="${OUT_DIR:-/var/lib/acdc/selinux}"
mkdir -p "${OUT_DIR}/${POLICY}"

cat > "${OUT_DIR}/${POLICY}/${POLICY}.te" <<EOF
policy_module(${POLICY}, 1.0)

type ${DOMAIN};
EOF

if [[ "${NET_BIND}" == "yes" ]]; then
  cat >> "${OUT_DIR}/${POLICY}/${POLICY}.te" <<EOF
corenet_tcp_bind_generic_node(${DOMAIN})
corenet_udp_bind_generic_node(${DOMAIN})
EOF
fi

touch "${OUT_DIR}/${POLICY}/${POLICY}.fc"
make -f /usr/share/selinux/devel/Makefile -C "${OUT_DIR}/${POLICY}" "${POLICY}.pp"
semodule -i "${OUT_DIR}/${POLICY}/${POLICY}.pp"

echo "Policy installed: ${POLICY}"
