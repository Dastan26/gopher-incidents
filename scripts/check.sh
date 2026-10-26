#!/usr/bin/env bash
# Проверяет текущий уровень: сборка, vet, приёмочные тесты уровня,
# а на некоторых уровнях ещё и мутанты (проверка твоих собственных тестов).
set -uo pipefail
source "$(dirname "$0")/lib.sh"
need_go

n=$(current_level)
dir=$(level_dir "$n") || die "Не нашёл каталог уровня $n"
printf '%s\n\n' "${B}Проверка уровня $n: $(level_title "$dir")${R}"

copied=()
backup=""
target=""
cleanup() {
  local f
  for f in ${copied[@]+"${copied[@]}"}; do rm -f "$f"; done
  if [ -n "$backup" ] && [ -f "$backup" ]; then cp "$backup" "$target"; rm -f "$backup"; fi
}
trap cleanup EXIT INT TERM

fail() {
  printf '\n%s\n' "${RED}${B}Красный.${R} $*"
  printf '%s\n' "${DIM}Застрял? make hint. Перечитать тикет: make ticket${R}"
  exit 1
}

step() {
  local name="$1"; shift
  local out
  printf '  %-28s' "$name"
  if out=$("$@" 2>&1); then
    printf '%s\n' "${GRN}ok${R}"
  else
    printf '%s\n\n' "${RED}FAIL${R}"
    printf '%s\n' "$out" | sed 's/^/    /' | head -n 60
    fail "$name не прошёл."
  fi
}

step "Сборка" go build ./...
step "go vet" go vet ./...

# Приёмочные тесты уровня лежат в levels/NN/acceptance/<путь пакета>/*.go.txt
# и подкладываются в пакеты только на время проверки.
if [ -d "$dir/acceptance" ]; then
  while IFS= read -r src; do
    rel="${src#"$dir/acceptance/"}"
    dst="${rel%.txt}"
    cp "$src" "$dst"
    copied+=("$dst")
  done < <(find "$dir/acceptance" -name '*.go.txt' -type f)
fi
step "Тесты (твои + приёмочные)" go test -count=1 ./...
cleanup
copied=()

# Мутанты: подменяем файл заведомо сломанной версией. Твои тесты обязаны
# упасть на каждой. Если какой-то мутант «выжил», тесты его не ловят.
if [ -d "$dir/mutants" ]; then
  target=$(cat "$dir/mutants/target")
  pkg="./$(dirname "$target")/"
  backup=$(mktemp)
  cp "$target" "$backup"
  survived=()
  printf '  %s\n' "Мутанты (ломаем $target, твои тесты должны упасть):"
  for m in "$dir"/mutants/*.go.txt; do
    name=$(basename "$m" .go.txt)
    desc=$(sed -n 's#^// MUTANT: ##p' "$m" | head -n 1)
    cp "$m" "$target"
    printf '    %-26s' "$name"
    if go test -count=1 "$pkg" >/dev/null 2>&1; then
      printf '%s\n' "${RED}выжил${R}  ${DIM}${desc}${R}"
      survived+=("$name")
    else
      printf '%s\n' "${GRN}пойман${R}"
    fi
    cp "$backup" "$target"
  done
  rm -f "$backup"; backup=""
  if [ "${#survived[@]}" -gt 0 ]; then
    fail "Твои тесты пропускают ${#survived[@]} из сломанных версий. Добавь случаи, которые их поймают."
  fi
fi

# Победа: время и подсказки.
started=$(cat "$STATE/started-$n" 2>/dev/null || date +%s)
minutes=$(( ($(date +%s) - started + 59) / 60 ))
hints=$(cat "$STATE/hints-$n" 2>/dev/null || echo 0)
if [ ! -f "$STATE/done-$n" ]; then
  echo "$minutes $hints" > "$STATE/done-$n"
fi

printf '\n%s\n' "${GRN}${B}Зелёный. Инцидент закрыт.${R} ${DIM}${minutes} мин, подсказок: ${hints}${R}"
printf '\n%s\n' "Дальше:"
printf '  %s\n' "1. Закоммить фикс:   git commit -am \"fix: <что починил>\""
printf '  %s\n' "2. Ревью у тимлида:  git show | pbcopy   и вставь Claude"
nextn=$((n + 1))
if level_dir "$nextn" >/dev/null; then
  printf '  %s\n' "3. Следующий:        make level-$nextn"
else
  printf '  %s\n' "3. Это последний уровень на сегодня. Новые добавит Claude."
fi
