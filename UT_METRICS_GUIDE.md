# Unit Test — Run & Evaluate Metrics

Operational runbook for this repo's unit-test suite. For the **methodology** (why we track what we track), see [`BDD_COVERAGE.md`](BDD_COVERAGE.md). This guide is the **how**: one-liner commands, what each number means, how to add a new BDD test, how to debug a failing dashboard.

---

## TL;DR

```bash
# 1. Install once
make deps              # go mod tidy + download

# 2. Run everything — line coverage + AC (BDD) coverage + dashboard
bash scripts/ac-coverage.sh
```

Expected output (one page, `exit 0` when healthy):

```
================================================================
   BDD · AC COVERAGE REPORT
================================================================
ACs declared in USER_STORIES.md : 89
ACs tagged in tests             : 56
ACs required (under ✅ stories) : 47
ACs missing                     : 0
AC coverage over ✅ stories     : 100.0%

---- Story health dashboard ----
  Module                  Scenario Line (claim)   Story health
  ------                  -------- ------------   ------------
  Auth                        100%          89%           89% ✅
  Vendor KYC                  100%          95%           95% ✅
  Product                     100%          73%           73% 🟡
  Commission                  100%          96%           96% ✅
  Checkout                     73%          80%           58% 🟡
  Variants & Stock            100%          80%           80% ✅
  Order Lifecycle               0%          —             0% ❌
  Returns & Refunds             0%          —             0% ❌
  Wallet & Settle               0%          —             0% ❌
  Notifications                 0%          —             0% ❌
  Validation                  100%         100%          100% ✅

OK — all ✅-story ACs are covered by tagged BDD tests.
```

---

## 1. Prerequisites

| Tool | Why | Install (macOS) |
|---|---|---|
| Go ≥ 1.22 | build + `go test -cover` | `brew install go` |
| bash + awk + grep + sed | runs the dashboard script | built-in |
| Docker + Docker Compose | only for integration / end-to-end tests | `brew install --cask docker` |

The BDD dashboard does **not** require the database or Redis to be running — unit tests use mocks.

---

## 2. Commands

### 2.1 Dashboard + AC gate

```bash
bash scripts/ac-coverage.sh          # runs from repo root; exits 1 if a ✅-story AC is untagged
# or
make ac-coverage                     # identical, wrapped in Make
```

### 2.2 Run only BDD-tagged tests (fast)

```bash
go test -v -run "^TestUS_" ./internal/...
# or
make test-bdd
```

The `^TestUS_` filter selects only tests named per our convention (see §4). Existing legacy tests are skipped — useful when validating that AC tagging itself works.

### 2.3 Run the full suite + line coverage

```bash
go test -count=1 -coverprofile=cover.out ./internal/...
go tool cover -func=cover.out | tail -1      # package total
go tool cover -html=cover.out -o cov.html    # annotated source, open in browser
# or
make test-coverage                            # writes coverage_report.html
```

### 2.4 Bundle (what CI runs)

```bash
make bdd-report                      # BDD tests first, then the dashboard, then line cov
```

---

## 3. Reading the dashboard

Three columns; each answers a different question.

| Column | What it measures | Tool |
|---|---|---|
| **Scenario** | % of Acceptance Criteria (ACs) in [`USER_STORIES.md`](USER_STORIES.md) that have at least one tagged test | `grep` of `AC:` comments and `TestUS_*` function names |
| **Line (claim)** | Function-average line coverage across the Go files that implement the module | `go tool cover -func` |
| **Story health** | `scenario × line`, rounded — a combined health bar | computed |

Status icons:

| Icon | Meaning | Health band |
|---|---|---|
| ✅ | Healthy | ≥ 80 % |
| 🟡 | Gaps remain | 1–79 % |
| ❌ | Story has no implementation **or** story is marked ❌ in [`USER_STORIES.md`](USER_STORIES.md) | 0 % |

### Why "Line (claim)" can be — instead of a number

Modules marked ❌ in [`USER_STORIES.md`](USER_STORIES.md) have no code to cover (the service methods don't exist). The script deliberately reports `—` instead of `0%` because `0%` coverage of an empty function set would look like a bug.

### Why Product shows 73% line coverage

Two service methods (`GetProductBySlug`, `ListPendingApproval`) are reachable from handlers but not under any user story, so they sit at 0 %. That pulls the module's function-average down. Options:
1. Add tests + a story for them (raises the module to ~90 %).
2. Mark them `// Infrastructure / non-AC` and exclude from the module count (cosmetic fix).

### Scenario vs Line — what each diagnosis means

See [`BDD_COVERAGE.md`](BDD_COVERAGE.md#3-why-you-need-both) for the 2×2 matrix. Short version:

| scenario | line | diagnosis |
|---|---|---|
| high | high | healthy |
| high | low | legit behaviors covered but branches untested |
| low | high | code runs but nobody verified intent |
| low | low | don't ship |

---

## 4. BDD test conventions — how to write one

Every behavioral unit test follows this shape. Keep it boring. No DSL.

```go
// =============================================================================
// US-VEN-002 · Admin drives vendor status machine
//
// AC: US-VEN-002 AC1, AC2, AC3
// =============================================================================

func TestUS_VEN_002_VendorStatusTransitions(t *testing.T) {
    // AC1
    // Given  vendor status = pending
    // When   ApproveVendor(id) is called
    // Then   status becomes approved
    t.Run("AC1/approve_sets_status_approved", func(t *testing.T) {
        // --- Given ---
        svc, vr, _ := newVendorSvcBDD(t)
        vr.On("UpdateStatus", uint(10), "approved").Return(nil)

        // --- When ---
        err := svc.ApproveVendor(10)

        // --- Then ---
        require.NoError(t, err)
        vr.AssertCalled(t, "UpdateStatus", uint(10), "approved")
    })
}
```

Mandatory ingredients (enforced by the coverage script):

1. **Function name:** `TestUS_<CODE>_<NNN>_<PascalCaseStory>` — e.g. `TestUS_VEN_002_VendorStatusTransitions`. The `CODE` / `NNN` must match the story in [`USER_STORIES.md`](USER_STORIES.md).
2. **File-header comment**: `AC: US-<CODE>-<NNN> AC1, AC2, AC3` listing every AC the function covers. The coverage script greps this pattern.
3. **Sub-test per AC**: `t.Run("AC<n>/<scenario>", ...)`. One AC = one sub-test, or one AC = one row in a table-driven sub-test. The `AC<n>` prefix is what the script matches.
4. **Given / When / Then comment block** before the function, describing business state (not mock setup).
5. **Inline markers**: `// --- Given ---`, `// --- When ---`, `// --- Then ---` inside the body.
6. **Named assertions**: `require.NoError(t, err, "why")` and `assert.Equal(t, want, got, "why")`. Test failures should read like a business sentence.

Table-driven form (preferred when an AC has many row-level scenarios):

```go
// AC: US-VEN-001 AC2, AC3, AC4
func TestUS_VEN_001_ValidateVendor(t *testing.T) {
    cases := []struct {
        ac   string            // "AC2", "AC3", "AC4"
        name string
        mut  func(*models.Vendor)
        want string            // substring of expected error, "" = nil
    }{
        {"AC2", "unknown_model",  func(v *models.Vendor){ v.CommissionModel = "flat" }, "commission model"},
        {"AC3", "rate_negative",  func(v *models.Vendor){ v.CommissionRate  = -0.1 },   "rate"},
        {"AC4", "missing_userID", func(v *models.Vendor){ v.UserID          = 0 },      "user"},
    }
    for _, tc := range cases {
        t.Run(tc.ac+"/"+tc.name, func(t *testing.T) {
            v := validVendor(); tc.mut(v)
            err := svc.ValidateVendor(v)
            if tc.want == "" { require.NoError(t, err) } else {
                require.ErrorContains(t, err, tc.want)
            }
        })
    }
}
```

### Naming rules — strict

- **DO** write `TestUS_<CODE>_<NNN>_...` for new BDD-layer tests.
- **DO** keep existing AAA-style tests (`TestValidateVendor_Valid`, etc.) untouched — they provide line coverage.
- **DO NOT** reuse a story code across files. One story's ACs live in one `TestUS_` function.
- **DO NOT** skip the `AC:` header comment — the coverage script has two parsers (header-based and function-name-based) but the header is the primary source.

### Where to put the file

- Service-level BDD tests → `internal/services/<domain>_bdd_test.go`.
- Handler/validation BDD tests → `internal/handlers/<domain>_bdd_test.go`.
- Reuse existing mocks from the same package (see top of the existing `*_bdd_test.go` files for examples).

---

## 5. Adding a new AC to USER_STORIES.md

1. Find the story section in [`USER_STORIES.md`](USER_STORIES.md), e.g. `### US-AUTH-002 🟡 · Login issues a token pair`.
2. Append the AC as a markdown bullet:
   ```markdown
   - **AC3** — **Given** a user whose status is `inactive`; **When** `Login` is called; **Then** a 401 is returned and nothing is cached.
   ```
3. Bump the per-story **Gaps** line so reviewers see it's unfinished.
4. Run `bash scripts/ac-coverage.sh` — expect AC3 to appear in the "MISSING" list.
5. Write the test. Rerun the script. Dashboard should update.
6. If the story is currently `✅`, the script will **fail CI** until the new AC is tagged. This is intentional.

---

## 6. Interpreting the Story health formula

```
health = scenario × line           (percent × percent, rounded half-up)
```

Why `×` and not `min`?
- **Multiplication punishes both dimensions at once** — 80 % × 80 % = 64 % feels off, but it correctly reflects that each dimension is independently incomplete.
- `min` hides asymmetry. A 100 % scenario / 60 % line module looks the same as a 60 % / 100 % module under `min`. Under multiplication they both score 60 %, but we can see which dimension is lagging in the adjacent columns.
- Matches the original [`USER_STORIES.md`](USER_STORIES.md) screenshot used by the team.

A ship-blocking `❌` status forces `health = 0 %` regardless of line coverage — see [`BDD_COVERAGE.md` §4](BDD_COVERAGE.md).

### Rounding

- Scenario: `(100·tagged + declared/2) / declared` → half-up integer percent.
- Line: function-average from `go tool cover -func`, rounded to integer.
- Health: `(scenario·line + 50) / 100` → half-up integer percent.

The script uses integer math throughout for reproducibility across shells.

---

## 7. CI integration

Minimal GitHub Actions job:

```yaml
- name: BDD + AC coverage
  run: |
    bash scripts/ac-coverage.sh
  # Exit 1 ⇒ CI fails ⇒ PR cannot merge.
```

Optional — gate line coverage separately via [`codecov.yml`](codecov.yml). The dashboard does **not** enforce a line-coverage threshold by itself; it reports it. Raise the Codecov `target` per module as stories reach ✅ (see [`BDD_COVERAGE.md` §6](BDD_COVERAGE.md)).

---

## 8. Troubleshooting

| Symptom | Cause | Fix |
|---|---|---|
| `bash: scripts/ac-coverage.sh: No such file or directory` | You're in a subdirectory | `cd` to repo root first, or run `make ac-coverage` |
| `ACs declared: 0` | [`USER_STORIES.md`](USER_STORIES.md) missing or heading format changed | Make sure story headings match `### US-<CODE>-<NNN> …` and AC bullets match `- **AC<n>** — …` |
| `ACs tagged: 0` but you wrote tests | Test function name is missing the `US-CODE-NNN` pattern | Rename the function to `TestUS_<CODE>_<NNN>_...` or add an `AC:` header comment |
| "Missing" AC is actually tested | Sub-test name not prefixed with `AC<n>/` | Rename the `t.Run` to `AC<n>/<scenario>` |
| Line coverage drops to 0 % for a module after refactor | Old AAA tests still live, new code paths are untouched | Run `go tool cover -func=cover.out | grep <file>` to find the uncovered functions |
| `Xcode license agreements` error on `make` | macOS toolchain lock | Run `sudo xcodebuild -license accept`, or invoke `bash scripts/ac-coverage.sh` directly |
| Script fails silently with exit 1 | `set -e` on a grep-with-no-match | Already wrapped in `|| true`; if you added new greps to the script, wrap them the same way |
| Module table shows `—` for a module with code | Source file not listed in `MODROWS` of the script | Add the `.go` file name to the module's row in [`scripts/ac-coverage.sh`](scripts/ac-coverage.sh) |

---

## 9. What each layer does NOT measure

Unit tests are the foundation, not the ceiling. For the money-critical paths in this repo (commission, checkout rollback, wallet ledger, stock race), see [`BDD_COVERAGE.md` §5](BDD_COVERAGE.md) for the rest of the testing pyramid:

- **Integration tests** against real Postgres + Redis — proves transactional rollback.
- **Handler / API tests** via `httptest` — proves Gin wiring + middleware.
- **Concurrency tests** with `-race` — proves stock and wallet aren't over-sold.
- **Contract / OpenAPI diff** — proves we don't silently break consumers.

Unit tests should stay fast (`<5 s` full run), mock-driven, and grounded in ACs. Anything slower or I/O-bound belongs in the integration tier.

---

## 10. Keeping this doc honest

When any of the following changes, update the TL;DR output snippet at the top and the troubleshooting table:

- New module added to [`USER_STORIES.md`](USER_STORIES.md) → update `MODROWS` in the script + re-verify the dashboard.
- Rounding rules changed → update §6.
- CI pipeline changed → update §7.
- Test discovery regex in the script changed → update §4 and §8.

No dashboards in this doc are auto-generated; the **authoritative** dashboard is what `bash scripts/ac-coverage.sh` prints right now.
