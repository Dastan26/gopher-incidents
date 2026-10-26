#!/usr/bin/env bash
set -euo pipefail
source "$(dirname "$0")/lib.sh"

n=$(current_level)
dir=$(level_dir "$n") || die "Не нашёл каталог уровня $n"
echo
render_md "$dir/ticket.md"
echo
