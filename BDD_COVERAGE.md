# BDD Coverage Methodology

How we measure unit-test coverage on the coolmate backend. Line coverage alone doesn't prove behavior was verified; we track **two** metrics and combine them into a per-module health score.

---

## Current reality (2026-04-24)

- **BDD scenario coverage: 0%.** No test in the suite is named per the `TestUS_<CODE>_<NNN>_AC<n>_...` convention yet. Every AC in [USER_STORIES.md](USER_STORIES.md) is a gap under strict BDD accounting, even where a legacy testify test happens to cover the behavior.
- **Line coverage: unmeasured** in this session. [UNIT_TESTING_SUMMARY.md](UNIT_TESTING_SUMMARY.md) claims ~90% on implemented services (`go test ./internal/services/... -cover`), but we haven't verified it. Codecov CI target is only 20% per [codecov.yml](codecov.yml) — it catches regressions, not quality.
- **Legacy-mapped coverage: ~54%** — 48 of 88 ACs have at least one existing testify test listed in their `Tests:` footer in USER_STORIES.md. These are **behaviorally covered but not traceable** by tooling. Migrating them to the BDD name convention (or adding AC tags) converts this into real scenario coverage.

The three numbers differ on purpose. Scenario coverage is what a reviewer can verify by grepping test names; line coverage is what the tool reports; legacy mapping is what we *claim* exists.

---

## 1. Scenario coverage — from [USER_STORIES.md](USER_STORIES.md)

**Definition.** The percentage of documented acceptance criteria (ACs) that have a matching unit test named `TestUS_<CODE>_<NNN>_AC<n>_<shortScenario>`.

**Formula:**

```
scenarioCoverage = testedACs / totalACs
```

**Why the name convention matters.** A test named `TestApproveVendor_Success` proves a function runs; it does not tie to a business scenario by any machine-readable means. A test named `TestUS_VEN_002_AC1_Approve` is grep-able, 1:1 with the spec, and survives refactors.

**How to refresh.**
1. `go test -list '.*' ./... | grep '^TestUS_'` — list BDD-tagged tests.
2. Extract AC ids from `USER_STORIES.md` (pattern `US-[A-Z]+-\d+ · .*AC\d+`).
3. Diff: matched → tested; declared without match → gap; test without match → orphan (the story file is behind reality).
4. Update the per-module table below.

A ~50-line script under `scripts/bdd-coverage.sh` automates this and can gate CI.

---

## 2. Line coverage — from `go test -cover`

**Definition.** Percentage of Go statements (and branches, with detection on) executed by the suite. This is what Codecov tracks.

**Commands:**

```
go test ./internal/... -coverprofile=cover.out
go tool cover -func=cover.out    # per-function %
go tool cover -html=cover.out    # annotated source
```

**Current Codecov config** ([codecov.yml](codecov.yml)):
- Target: **20%** (floor), threshold 5%
- Flag: `unittests`, scoped to `internal/`
- Ignores: `tests/`, `**/*_test.go`, `**/mocks`
- Branch detection: on

The 20% floor is a safety net, not a quality bar. Raise per-module as stories reach ✅.

---

## 3. Why you need both

Line coverage says "this code ran in a test." Scenario coverage says "this documented behavior was verified." Different questions.

| Scenario cov | Line cov | Diagnosis |
|---|---|---|
| High | High | **Healthy** — documented behaviors verified + code paths exercised |
| High | Low | Happy-path only; edge branches silently untested |
| Low | High | Code runs under tests but nobody verified intent → either spec is incomplete or tests assert implementation detail |
| Low | Low | Don't ship |

**We are in the Low/High cell today** (scenario 0% strict, line claimed high). That is normal for a codebase that adopted BDD after the fact — the work is to rename/tag existing tests and close the 40 real gaps, not to rewrite anything.

---

## 4. Health score per module

**Formula:**

```
moduleHealth = scenarioCoverage × lineCoverage × statusWeight

statusWeight = 1.00   if all stories ✅
             = 0.50   if any 🟡 (implemented but untested)
             = 0.00   if any ❌ (not implemented)  — blocks ship
```

A module with ship-blocking gaps (`❌`) scores 0 regardless of line coverage on its other stories.

**Bands:** `≥ 0.80` shippable · `0.50–0.80` ship with known risk · `< 0.50` not ready.

### Current snapshot — honest

Scenario coverage is **0% under strict BDD** because no test follows the naming convention yet. Line coverage is unmeasured; `—` entries need `go test -cover`. The `Legacy map` column shows how many ACs have behaviorally-aligned legacy tests (from USER_STORIES.md), i.e. the work-in-hand to convert into real BDD coverage.

| Module | Stories | BDD scen. | Legacy map | Line cov | Status | Health | Next action |
|---|--:|--:|--:|:--:|:--:|--:|---|
| Vendor KYC | 4 | 0% | 11/11 | — | ✅ | 0.00 | rename 24 tests to `TestUS_VEN_*` → instant 100% |
| Product | 4 | 0% | 10/10 | — | ✅ | 0.00 | rename 22 tests to `TestUS_PRD_*` |
| Commission | 3 | 0% | 10/10 | — | ✅ | 0.00 | rename 15 tests to `TestUS_COM_*` |
| Variants & Stock | 1 | 0% | 2/2 | — | 🟡 | 0.00 | rename 2 tests |
| Checkout (split/promo) | 2 of 4 | 0% | 8/11 | — | 🟡 | 0.00 | rename 10 tests + write AC3/rollback/wallet |
| Request Validation | 1 | 0% | 7/7 | — | ✅ | 0.00 | rename 35+ validator tests |
| Auth | 3 | 0% | 0/7 | 0% | 🟡 | 0.00 | create `auth_service_test.go` — 7 new ACs |
| Order Lifecycle | 3 | 0% | 0/8 | — | ❌ | 0.00 | implement service, then test |
| Manual Payment | 1 | 0% | 0/3 | — | ❌ | 0.00 | implement service, then test |
| Returns & Refunds | 3 | 0% | 0/8 | — | ❌ | 0.00 | implement service, then test |
| Wallet & Settlement | 3 | 0% | 0/7 | — | ❌ | 0.00 | implement service, then test |
| Vendor Notifications | 2 | 0% | 0/4 | — | ❌ | 0.00 | implement service, then test |

**Overall MVP:** 0/88 ACs under strict BDD · 48/88 (54.5%) legacy-mapped · line cov unmeasured.

### The path to real health scores

1. **Measure line coverage** — run `go test -cover`, replace `—` in the Line cov column.
2. **Rename legacy tests** — convert the 48 mapped tests to `TestUS_<CODE>_<NNN>_AC<n>_...`. After this, 7 of 12 modules hit 100% BDD coverage overnight.
3. **Close the 7 Auth gaps** — create `auth_service_test.go`. Moves Auth from `0.00` to healthy in one PR.
4. **Implement ❌ modules** — Order Lifecycle, Returns, Wallet, Payment, Notifications. Each gets written with BDD names from day one.

---

## 5. Is unit testing enough? — Testing-pyramid evaluation

**Short answer: no.** For an e-commerce backend that moves money, split-orders, decrements stock, and settles vendors, unit tests alone — however well-covered — will not protect the domain. Unit tests verify logic in isolation with mocks; they don't catch:

- **Repository / SQL layer bugs** — constraints, migrations, transactional boundaries, race conditions on `UPDATE ... SET stock = stock - ?`
- **Handler wiring** — Gin middleware (JWT, CORS, rate limits), binding tags, error-to-HTTP mapping
- **Cross-service composition** — SplitOrder → stock decrement → wallet transaction rollback behavior under real failures (the `US-CHK-003` gap)
- **External I/O** — S3 uploads, Redis cache invalidation semantics, email delivery
- **Concurrency** — two customers racing the last unit of stock; two admin settlements racing on the same wallet
- **Contract stability** — breaking API changes that unit tests happily allow

### Recommended pyramid for this repo

| Level | Scope | Tooling | Priority |
|---|---|---|---|
| **Unit** | Services + validators with mocks | testify (in place) | **In progress** — close BDD gaps first |
| **Integration** | Repositories against real Postgres + Redis; transaction rollback; cache invalidation | `testcontainers-go` or existing [docker-compose.yml](docker-compose.yml); `_integration_test.go` + build tag | **Next (MVP blocker)** — the `US-CHK-003` rollback story cannot be proven with mocks |
| **Handler / API** | Gin routes end-to-end in-process | `httptest` + mocked services OR wired to integration stack | Before first external customer |
| **Contract** | OpenAPI schema drift | `swagger` gen + schema diff in CI | Before a frontend/mobile consumer exists |
| **E2E smoke** | One golden-path checkout against a running stack | `k8s` already configured; simple Go or HTTP test | Pre-launch |
| **Load / concurrency** | Stock decrement + commission math under 100+ concurrent checkouts | `k6` or `vegeta`; race detector on unit tests (`-race`) | Before high-traffic launch |
| **Security** | Input fuzzing on handlers, SQL-injection smoke, authz matrix | `go-fuzz` / handler-level negative tests | Pre-launch |

### What specifically to add next (in order)

1. **Integration suite for the checkout critical path.** Use the existing docker-compose services. Prove `SplitOrder → DecrementStock → RecordWalletTransaction` rolls back atomically on any failure. This is the highest-value test not currently possible with unit tests.
2. **Repository-layer tests.** Zero exist today. Every `*_repository.go` under `internal/repositories/` should have integration tests that exercise real SQL, including constraint violations and optimistic-lock paths.
3. **Handler tests via `httptest`.** Covers binding tags + middleware + service composition without rewriting to E2E scale.
4. **`go test -race` in CI.** Catches stock / wallet race conditions early at nearly zero cost.

Unit tests with strict BDD tagging are the **necessary foundation**, not the ceiling. The Checkout rollback, wallet ledger, and stock-race risks in this repo are pyramid-middle problems and need pyramid-middle tests to be honestly covered.

---

## 6. Updating this doc

- Refresh the **Current snapshot** table when `USER_STORIES.md` changes or a PR closes AC gaps.
- Replace `—` in Line cov with measured values from `go tool cover -func` after each significant test PR.
- When a module reaches ≥ 0.80 health, raise the Codecov `target` for its flag to lock in the improvement.
