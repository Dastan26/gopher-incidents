#!/usr/bin/env bash
# Выдаёт следующую подсказку по текущему уровню. Подсказки идут лестницей:
# сначала куда смотреть, потом где именно, потом объяснение.
set -euo pipefail
source "$(dirname "$0")/lib.sh"

n=$(current_level)
dir=$(level_dir "$n") || die "Не нашёл каталог уровня $n"
used=$(cat "$STATE/hints-$n" 2>/dev/null || echo 0)
next=$((used + 1))
file="$dir/hints/$next.md"

if [ ! -f "$file" ]; then
  total=$(ls "$dir"/hints/*.md | wc -l | tr -d ' ')
  printf '%s\n' "Подсказки кончились: ты открыл все $total. Перечитать: cat $dir/hints/*.md"
  printf '%s\n' "Если совсем застрял, покажи Claude вывод make check и свой git diff."
  exit 0
fi

echo "$next" > "$STATE/hints-$n"
echo
render_md "$file"
echo
total=$(ls "$dir"/hints/*.md | wc -l | tr -d ' ')
printf '%s\n' "${DIM}Подсказка $next из $total.${R}"
