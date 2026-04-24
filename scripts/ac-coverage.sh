#!/usr/bin/env bash
#
# ac-coverage.sh — measure BDD/AC coverage against USER_STORIES.md
#
# Prints a "Story health" dashboard that mirrors the one in
# TEST_BDD_ALIGNMENT.md / the client-facing screenshot:
#
#   Module        | Scenario | Line (claim) | Story health
#   --------------+----------+--------------+-------------
#   Vendor KYC    |  100%    |      92%     |     92% ✅
#   ...
#
# Exit code:
#   0 — all ACs declared under ✅-marked stories are tested
#   1 — one or more ACs under a ✅ story are missing a test tag
#
# Requires: bash, awk, grep, sed, go (for line coverage).
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
STORIES="$ROOT/USER_STORIES.md"
TESTS_DIR="$ROOT/internal"

if [[ ! -f "$STORIES" ]]; then
  echo "USER_STORIES.md not found at $STORIES" >&2
  exit 2
fi

tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

# Walk USER_STORIES.md once and emit each AC paired with the story code
# that heads its section. Produces lines like:
#    US-VEN-001 AC1 STATUS
# where STATUS is one of ✅ / 🟡 / ❌ / ?.
awk '
  /^### US-[A-Z]+-[0-9]+/ {
    match($0, /US-[A-Z]+-[0-9]+/); code = substr($0, RSTART, RLENGTH)
    status = "?"
    if ($0 ~ /✅/) status = "OK"
    else if ($0 ~ /🟡/) status = "PARTIAL"
    else if ($0 ~ /❌/) status = "MISSING"
    next
  }
  /- \*\*AC[0-9]+/ {
    match($0, /AC[0-9]+/); ac = substr($0, RSTART, RLENGTH)
    print code, ac, status
  }
' "$STORIES" | sort -u > "$tmp/all_acs.tsv"

# 1. Every AC id declared in USER_STORIES.md.
awk '{print $1, $2}' "$tmp/all_acs.tsv" > "$tmp/declared.txt"

# 2. Every AC id that a test file tags, from either:
#      (a) an "AC: US-XXX-### AC1, AC2, AC3" header comment, OR
#      (b) a t.Run("AC<n>/...") subtest nested inside a TestUS_CODE_NNN_... function.
#    We normalise both shapes into lines like "US-CODE-NNN ACn".
find "$TESTS_DIR" -name '*_test.go' -print0 | xargs -0 awk '
  # Update current story code when we see either:
  #   (a) any "US-CODE-NNN" substring in a comment / identifier, or
  #   (b) a "TestUS_CODE_NNN_..." function declaration.
  {
    if (match($0, /US-[A-Z]+-[0-9][0-9][0-9]/)) {
      code = substr($0, RSTART, RLENGTH)
    }
    # Function-name form: TestUS_AUTH_001_...
    if (match($0, /TestUS_[A-Z]+_[0-9][0-9][0-9]/)) {
      s = substr($0, RSTART, RLENGTH)
      gsub(/^TestUS_/, "", s); gsub(/_/, "-", s)
      code = "US-" s
    }
    if (code == "") next

    # Emit every ACn token on this line paired with the active code.
    line = $0
    while (match(line, /AC[0-9]+/)) {
      ac = substr(line, RSTART, RLENGTH)
      print code, ac
      line = substr(line, RSTART + RLENGTH)
    }
  }
' | sort -u > "$tmp/tagged.txt"

# 3. ACs belonging to ✅-marked stories.
awk '$3=="OK" {print $1, $2}' "$tmp/all_acs.tsv" > "$tmp/required.txt"

decl=$(wc -l < "$tmp/declared.txt" | tr -d " ")
tagged=$(wc -l < "$tmp/tagged.txt" | tr -d " ")
required=$(wc -l < "$tmp/required.txt" | tr -d " ")

# Missing = required-but-not-tagged
comm -23 "$tmp/required.txt" "$tmp/tagged.txt" > "$tmp/missing.txt"
missing=$(wc -l < "$tmp/missing.txt" | tr -d " ")

echo "================================================================"
echo "   BDD · AC COVERAGE REPORT"
echo "================================================================"
echo "ACs declared in USER_STORIES.md : $decl"
echo "ACs tagged in tests             : $tagged"
echo "ACs required (under ✅ stories) : $required"
echo "ACs missing                     : $missing"
if (( required > 0 )); then
  pct=$(awk "BEGIN { printf \"%.1f\", 100 * ($required - $missing) / $required }")
  echo "AC coverage over ✅ stories     : ${pct}%"
fi
echo

if (( missing > 0 )); then
  echo "---- MISSING AC TAGS (required but not covered) ----"
  sed 's/^/   ✗  /' "$tmp/missing.txt"
  echo
fi


# 5. Per-module Line coverage and Story health dashboard.
#    Story-health formula (mirrors USER_STORIES.md screenshot):
#        health = scenario_pct × line_pct
#    An ❌ for "no implementation" shows as "—".
echo "---- Story health dashboard ----"
if command -v go >/dev/null 2>&1; then
  (cd "$ROOT" && go test -count=1 -coverprofile="$tmp/cover.out" ./internal/... >/dev/null 2>&1 || true)
  if [[ -f "$tmp/cover.out" ]]; then
    (cd "$ROOT" && go tool cover -func="$tmp/cover.out") > "$tmp/cov_func.txt"
  else
    : > "$tmp/cov_func.txt"
  fi
fi

# module | story-prefix | source files (space separated)
MODROWS=(
  "Auth              |AUTH|auth_service.go"
  "Vendor KYC        |VEN|vendor_service.go"
  "Product           |PRD|product_service.go"
  "Commission        |COM|commission_service.go"
  "Checkout          |CHK|order_service.go"
  "Variants & Stock  |INV|order_service.go"
  "Order Lifecycle   |ORD|"
  "Returns & Refunds |RTN|"
  "Wallet & Settle   |WAL|"
  "Notifications     |NOT|"
  "Validation        |VAL|request_validators.go"
)

# line_cov_for <files...>  — print weighted avg coverage across these files.
line_cov_for() {
  local files=("$@")
  [[ ${#files[@]} -eq 0 ]] && { printf "NA"; return; }
  awk -v want="$(IFS='|'; echo "${files[*]}")" '
    BEGIN { split(want, wa, "|"); for (i in wa) target[wa[i]] = 1; total_fns = 0; pct_sum = 0 }
    {
      # columns: path:line:   FuncName   NN.N%
      if (NF < 3) next
      pct = $NF
      sub(/%$/, "", pct)
      path = $1
      n = split(path, parts, "/")
      file = parts[n]
      sub(/:.*/, "", file)
      if (!(file in target)) next
      total_fns++
      pct_sum += pct
    }
    END {
      if (total_fns == 0) { printf "NA"; exit }
      printf "%.0f", pct_sum / total_fns
    }
  ' "$tmp/cov_func.txt"
}

printf "  %-22s %10s %12s %14s\n" "Module" "Scenario" "Line (claim)" "Story health"
printf "  %-22s %10s %12s %14s\n" "------" "--------" "------------" "------------"

for row in "${MODROWS[@]}"; do
  label="${row%%|*}"
  rest="${row#*|}"
  prefix="${rest%%|*}"
  files="${rest#*|}"

  decl_n=$(awk -v p="$prefix" 'BEGIN{n=0} {
    split($1, a, "-"); if (a[2] == p) n++
  } END{print n}' "$tmp/declared.txt")

  tag_n=$(awk -v p="$prefix" 'BEGIN{n=0} {
    split($1, a, "-"); if (a[2] == p) n++
  } END{print n}' "$tmp/tagged.txt")

  if (( decl_n == 0 )); then
    scen_pct="-"
    scen_num=0
  else
    # Round half-up: (100*tag + decl/2) / decl
    scen_num=$(( (100 * tag_n + decl_n / 2) / decl_n ))
    scen_pct="${scen_num}%"
  fi

  if [[ -z "$files" ]]; then
    line_pct="—"
    health_pct="0%"
    icon="❌"
  else
    # shellcheck disable=SC2206
    fa=($files)
    lc=$(line_cov_for "${fa[@]}")
    if [[ "$lc" == "NA" ]]; then
      line_pct="—"
      health_pct="0%"
      icon="❌"
    else
      line_pct="${lc}%"
      # Round half-up on the product as well
      health_num=$(( (scen_num * lc + 50) / 100 ))
      health_pct="${health_num}%"
      if   (( health_num >= 80 )); then icon="✅"
      elif (( health_num >= 40 )); then icon="🟡"
      else                              icon="🟡"
      fi
      if (( scen_num == 0 && lc == 0 )); then icon="❌"; fi
    fi
  fi

  printf "  %-22s %10s %12s %13s %s\n" "$label" "$scen_pct" "$line_pct" "$health_pct" "$icon"
done

echo
if (( missing > 0 )); then
  echo "FAIL — $missing AC(s) under ✅ stories are not tagged in any test."
  exit 1
fi
echo "OK — all ✅-story ACs are covered by tagged BDD tests."
