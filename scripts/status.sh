#!/usr/bin/env bash
# Доска инцидентов: что закрыто, что открыто, сколько времени ушло.
set -euo pipefail
source "$(dirname "$0")/lib.sh"

cur=$(cat "$STATE/current" 2>/dev/null || echo "")
closed=0
total=0
printf '\n%s\n\n' "${B}Доска инцидентов${R}"
for n in $(all_levels); do
  total=$((total + 1))
  dir=$(level_dir "$n")
  title=$(level_title "$dir")
  if [ -f "$STATE/done-$n" ]; then
    read -r minutes hints < "$STATE/done-$n"
    closed=$((closed + 1))
    printf '  %s    %-3s ' "${GRN}закрыт${R}" "$n"; pad "$title" 62; printf '  %s\n' "${DIM}${minutes} мин, подсказок ${hints}${R}"
  elif [ "$n" = "$cur" ]; then
    printf '  %s  %-3s %s\n' "${YEL}в работе${R}" "$n" "$title"
  else
    printf '  %s    %-3s %s\n' "${DIM}открыт${R}" "$n" "$title"
  fi
done
printf '\n  Закрыто %s из %s.' "$closed" "$total"
if [ -n "$cur" ] && [ ! -f "$STATE/done-$cur" ]; then
  printf ' Сейчас в работе уровень %s: make ticket\n\n' "$cur"
else
  next=""
  for n in $(all_levels); do [ -f "$STATE/done-$n" ] || { next=$n; break; }; done
  if [ -n "$next" ]; then printf ' Начать следующий: make level-%s\n\n' "$next"; else printf '\n\n'; fi
fi
