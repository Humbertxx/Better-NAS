# Run by renew-certs.timer
# also can run by hand

set -euo pipefail

HOST=${HOST}
CERT_DIR= ${CERT_DIR}
CADDY_CONTAINER=${CADDY_CONTAINER}

log() { printf '%s  %s\n' "$(date '+%Y-%m-%d %H:%M:%S')" "$*"; }

# --- preconditions -----------------------------------------------------------
mkdir -p "$CERT_DIR"

if ! command -v tailscale >/dev/null 2>&1; then
    log "ERROR: tailscale binary not found on PATH"; exit 1
fi

# --- issue / renew the cert --------------------------------------------------
# tailscale cert returns cached cert if valid, or reissues near expiry

log "Requesting cert for ${HOST}"
tailscale cert \
    --cert-file "${CERT_DIR}/${HOST}.crt" \
    --key-file  "${CERT_DIR}/${HOST}.key" \
    "${HOST}"

# Fail loudly if the files aren't actually there
if [[ ! -s "${CERT_DIR}/${HOST}.crt" || ! -s "${CERT_DIR}/${HOST}.key" ]]; then
    log "ERROR: cert or key missing/empty after issuance — not reloading Caddy"; exit 1
fi

# --- reload Caddy -------------------------------------
if docker ps --format '{{.Names}}' | grep -qx "$CADDY_CONTAINER"; then
    log "Reloading Caddy config"
    docker exec "$CADDY_CONTAINER" caddy reload --config /etc/caddy/Caddyfile
    log "Done — cert renewed and Caddy reloaded"
else
    log "WARNING: ${CADDY_CONTAINER} not running; cert written but no reload performed"
fi
