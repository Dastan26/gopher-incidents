# Общие функции для скриптов игры. Подключается через source.
# Пишем под bash 3.2: он стоит на macOS по умолчанию.

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
STATE="$ROOT/.gopher"
cd "$ROOT"

if [ -t 1 ]; then
  B=$'\033[1m'; DIM=$'\033[2m'; RED=$'\033[31m'; GRN=$'\033[32m'; YEL=$'\033[33m'; CYN=$'\033[36m'; R=$'\033[0m'
else
  B=; DIM=; RED=; GRN=; YEL=; CYN=; R=
fi

die() { printf '%s\n' "${RED}$*${R}" >&2; exit 1; }

need_go() {
  command -v go >/dev/null 2>&1 || die "Go не установлен. Поставь: brew install go"
}

# level_dir 3 -> levels/03-status-200
level_dir() {
  local d
  d=$(ls -d "levels/$(printf '%02d' "$1")-"* 2>/dev/null | head -n 1)
  [ -n "$d" ] && [ -d "$d" ] || return 1
  printf '%s' "$d"
}

level_title() { head -n 1 "$1/title"; }

current_level() {
  [ -f "$STATE/current" ] || die "Уровень не выбран. Начни с: make level-1"
  cat "$STATE/current"
}

all_levels() { ls -d levels/[0-9][0-9]-* 2>/dev/null | sed -E 's#levels/0*([0-9]+)-.*#\1#'; }

# Печатает markdown почти как есть, но заголовки и код выделяет цветом.
render_md() {
  local in_code=0
  while IFS= read -r line || [ -n "$line" ]; do
    case "$line" in
      '```'*) in_code=$((1 - in_code)); continue ;;
    esac
    if [ "$in_code" = 1 ]; then
      printf '    %s\n' "${CYN}${line}${R}"
    else
      case "$line" in
        '# '*) printf '%s\n' "${B}${line#\# }${R}" ;;
        '## '*) printf '\n%s\n' "${B}${YEL}${line#\#\# }${R}" ;;
        *) printf '%s\n' "$line" ;;
      esac
    fi
  done < "$1"
}
