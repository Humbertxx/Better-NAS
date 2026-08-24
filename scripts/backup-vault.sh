# Run by backup-vault.timer
# also can run by hand
set -euo pipefail

STACK_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

if [[ -z "${VAULT_PATH:-}" && -f "${STACK_ROOT}/ai-brain/.env" ]]; then
    VAULT_PATH="$(grep -E '^VAULT_PATH=' "${STACK_ROOT}/ai-brain/.env" | tail -n1 | cut -d= -f2- | tr -d \"\')"
fi
VAULT_PATH="${VAULT_PATH:-/home/humbertx/vault}"

GIT_AUTHOR_NAME="${GIT_AUTHOR_NAME:-vault-backup}"
GIT_AUTHOR_EMAIL="${GIT_AUTHOR_EMAIL:-vault-backup@nas.local}"
GIT_COMMITTER_NAME="${GIT_COMMITTER_NAME:-${GIT_AUTHOR_NAME}}"
GIT_COMMITTER_EMAIL="${GIT_COMMITTER_EMAIL:-${GIT_AUTHOR_EMAIL}}"
REMOTE="${BACKUP_GIT_REMOTE:-origin}"

export GIT_AUTHOR_NAME GIT_AUTHOR_EMAIL GIT_COMMITTER_NAME GIT_COMMITTER_EMAIL

# remember to add git add and git commit to the script