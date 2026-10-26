#!/usr/bin/env bash
# Запускает инцидент: ветка level-NN от main, поломка отдельным коммитом.
set -euo pipefail
source "$(dirname "$0")/lib.sh"

n="${1:-}"
[ -n "$n" ] || die "Укажи номер уровня: make level-1"
n=$((10#$n))
dir=$(level_dir "$n") || die "Уровня $n нет. Список: make status"

if [ -n "$(git status --porcelain --untracked-files=no)" ]; then
  printf '%s\n' "${YEL}В рабочей копии есть незакоммиченные правки:${R}"
  git status --short --untracked-files=no
  printf 'Они пропадут. Начать уровень %s заново? [y/N] ' "$n"
  read -r answer
  case "$answer" in y|Y|д|Д) ;; *) die "Отменено. Сохрани правки через git commit или git stash." ;; esac
fi

branch=$(printf 'level-%02d' "$n")
git checkout -q -f -B "$branch" main
git apply "$dir/break.patch"
git -c user.name="Дежурный по проду" -c user.email="oncall@gopher.local" \
  commit -q -am "incident: $(level_title "$dir")"

mkdir -p "$STATE"
echo "$n" > "$STATE/current"
echo 0 > "$STATE/hints-$n"
date +%s > "$STATE/started-$n"

echo
render_md "$dir/ticket.md"
echo
printf '%s\n' "${DIM}Ты на ветке $branch. Проверить: make check · подсказка: make hint · тикет ещё раз: make ticket${R}"
