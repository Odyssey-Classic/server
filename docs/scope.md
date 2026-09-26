# server — Scope of Responsibilities

Status: agreed 2026-09-25. This document defines what the `server` repository is
responsible for, what it deliberately is not, and the decisions that shape both.
It is the reference for routing work: if a task isn't claimed here, it belongs to
`client`, `proto`, `registry`, or `admin-tools`.

## What this is

The Odyssey engine — a **game-agnostic authoritative server** for persistent,
shared, tile-based worlds. The engine is the product; Odyssey Classic is the
first world built on it. The engine's capability envelope is set by what earlier
Odyssey servers could do, but none of their content ships with it.

Core systems live in Go. Content lives in data. Behaviour is overridable by
scripts. Anyone can run a world.

## Design targets

| Target | Value |
|---|---|
| Concurrent players per world | 100–500 |
| Simulation rate | 10 Hz default, configurable (validate under live load) |
| Hosting model | Many self-hosted worlds; discovery via `registry` |
| Movement authority | Server-authoritative tile steps; client predicts self, interpolates others |
| Transport | WebSocket + protobuf (message layer kept transport-agnostic) |
| Relational store | SQLite by default, behind an interface that leaves Postgres open |
| Distribution | Single static binary, no embedded UI |
| Author accessibility | Git literacy never required; available to those who want it |

## Responsibilities

### 1. Simulation and authority

Own the tick loop, world state, and every authoritative outcome. Resolve
movement, combat, and interaction; validate all client intent and reconcile
client mispredictions. The server is the only authority — no client assertion is
trusted, including from privileged clients.

### 2. Game systems

First-class engine primitives, not left to world authors:

- **Spatial and combat core** — maps and tiles, movement, line of sight, spawns, combat resolution, damage and death, inventory, items
- **Progression and content** — stats, skills, levels, experience, quests, dialogue, crafting, loot tables
- **Social** — chat and channels, parties, guilds, player trade, mail, friends and blocking
- **World and economy** — NPC shops, currency and sinks, housing, world events, day/night and weather, scheduled spawns

### 3. Scripting host

Embed **JavaScript/TypeScript via goja** so world authors can override engine
behaviour. Scripts are **semi-trusted**: operator-authored, so not assumed
hostile, but constrained — no filesystem or network access unless injected,
per-call execution budgets, and errors contained and logged rather than fatal.

Scripts may:

- **Hook and modify engine events** — alter or veto outcomes in flight
- **Replace whole systems** — swap combat formulas, loot rolls, progression curves, which requires engine systems be built behind script-visible seams
- **Create persistent content at runtime** — written to layer B (below), never into authored content
- **Read and write their own durable state** — namespaced, with script-provided migrations

Execution is **synchronous in-tick for gameplay hooks**, with deferred
scheduling for slow work that re-enters safely off-tick.

#### Scheduling and time

The engine provides **no per-concept timer systems.** Gated checks ("may this
player claim the daily reward?") are a stored timestamp compared against now,
evaluated lazily when asked — no timers, and correct across restarts.

Expiry needs a push, so the engine provides **one persisted scheduled-wakeup
queue**: `(fire_at, handler, payload)`, drained in order each tick, with
`schedule_in(duration)` as a helper. It replaces respawn timers, buff expiry,
growth timers and everything of that shape.

Wakeups missed while the server was down **fire on catch-up at boot**, in
`fire_at` order and rate-limited so boot is not a stampede. The handler receives
both `scheduled_for` and the actual `now`, because only the script can know
whether a late fire should still act — a boss respawn should, a stale daily reset
should not. Entries may carry a `drop_if_late_by` policy, and recurring entries
coalesce so a schedule that missed ten fires runs once.

Delivery is **at-least-once**: a crash between firing and committing effects
re-fires the wakeup. **Handlers must therefore be idempotent**, and the engine
passes a stable `fire_id` so authors can guard on it. This deliberately moves
responsibility to script authors, because a handler's effects may span stores the
engine cannot commit as one atomic unit — and it must be prominent in
author-facing documentation.

Scripts read time from their **invocation context** rather than an ambient clock,
giving one coherent `now` per tick and making scripts testable without clock
manipulation. Durations within a session use monotonic time; cross-restart
intervals use wall clock, which can jump.

### 4. Content model and lifecycle

Three layers, with different ownership and flow direction:

| Layer | Contents | Written by | Flows |
|---|---|---|---|
| **A — Authored content** | Maps, templates, scripts, dialogue | Humans, via API | Both directions |
| **B — Script-authored content** | Content minted at runtime by scripts | Scripts | Downstream only |
| **C — World and player state** | Characters, inventories, instances, live sim | Gameplay | Downstream only, opt-in |

**Layer A is versioned in immutable revisions,** Git-backed under the hood and
never exposed as Git through the UI. History is **linear per server — no
branching, no merging.** Booting at revision V is exactly reproducible, which is
why layer A is human-authored only: scripts write to B so that A stays a pure
function of what authors committed.

**Ingress is API-only** (validated bundle upload). The internal pipeline —
propose revision → validate → commit → publish — must stay transport-agnostic so
direct Git push can be added later without a second validation implementation.

**Cloning (live → dev)** is author-selectable: A, A+B, or A+B+anonymized
players. Full clones including credentials are not supported.

**Promotion (dev → live)** moves layer A only, as a **whole revision, atomically,
with a preview diff** showing what changes and what live state it affects.
Promotion is **fast-forward only**: if live took a hotfix meanwhile, dev must
re-sync and resolve first, keeping conflict resolution off the live path.

**Live content edits are permitted** but create a normal revision and are gated
by permission, so live is always at a known version.

**Engine upgrades auto-migrate content forward.** Migration is **non-destructive
by invariant** — it commits a new revision, leaving the prior one intact and
rollback-able, and is previewable before it applies.

### 5. Persistence

There is **no custom journal and no snapshot layer.** SQLite's WAL is the
durability mechanism and SQLite's recovery is the crash-recovery mechanism.
State divides two ways:

| Kind | Rule | Examples |
|---|---|---|
| **Valuable** | Committed in a transaction before the action is acknowledged | Item locations, currency, character records, quest and skill progress, permissions, housing, scheduled wakeups |
| **Not valuable** | Memory only, lost on restart by design | Sub-tile position and facing, combat timers, buffs, aggro, cooldowns, sessions, presence, rate limits |

Substrates by layer:

| Layer | Substrate |
|---|---|
| Authored content (A) | Git-backed revision store |
| Script-authored content (B) | SQLite, namespaced, above authored content |
| Accounts, characters, items, world state (C) | SQLite, behind an interface |
| Audit ledger | Transactional outbox in SQLite, shipped to a pluggable append-only sink |
| Live simulation | Memory |

#### Item identity

Every item instance has a **stable UUID and a single location field** —
`(container_kind, container_id, slot)`. Moving an item is one row update, so it
can never exist in two places or in none: duplication is structurally impossible
rather than something to be careful about. A unique constraint on
`(container, slot)` gives slot integrity. Ground drops are just a location, so
they survive restarts.

A **stack is one row with a quantity**, so moving 10 of 50 gold is a decrement
and an increment inside one transaction. **Currency is an item**, giving one
uniform model across coins, ammunition and collectibles.

**Split and merge are the only item operations that are not relocations** — split
mints a UUID, merge destroys one — and are therefore the only places a
duplication bug can originate. Both are single transactions.

The audit ledger consequently needs two entry shapes: UUID relocation for whole
items, and `(type, quantity, from, to)` for stack transfers.

The simulation holds items in memory as a **write-through cache; the database is
authoritative for location.**

#### Player position

The **anchor** — map and tile — is persistent; sub-tile offset and facing are
memory. The anchor is written **on transitions only** (map or zone change, death,
respawn, logout), which is affordable because worlds are built from small distinct
maps. A crash resumes a player at their last transition.

An anchor can dangle when a content change removes the map it referenced. A
**script hook gets first chance to remap it**, falling back to the world spawn
point, with the event reported to the operator. A bad anchor never blocks a login.

#### Restitution

Restitution reads the audit ledger and applies **compensating transactions** —
move the item back, restore the balance. Game logic is never re-executed, so
**no deterministic replay requirement exists anywhere in the engine.**

### 6. Identity and sessions

Each world owns **local accounts** and works fully standalone and offline.
Operators may opt into **registry SSO** so one account spans worlds. Session
lifecycle, authentication, and token issuance are the server's.

Pre-authentication work is rate-limited and cost-capped per IP: password hashing
is deliberately expensive, which makes unauthenticated login an amplifier if
left unbounded.

### 7. Permissions

**Fixed roles now** — owner, operator, content author, moderator, player — with
capability checks behind an indirection so granular capabilities can land later
without breaking callers. Gates live content editing, moderation, and every
control-plane operation.

### 8. Moderation and audit

The engine provides:

- **Player sanctions** — kick, temporary and permanent ban, mute, with reason and duration recorded
- **Live inspection and intervention** — locate and observe a player, inspect inventory, teleport, freeze an entity
- **Reporting and review queue** — player reports with relevant chat and action history attached
- **Rollback and restitution** — reverse a trade or item loss, restore a character to an earlier state

**Automated moderators** (including AI agents) are first-class API consumers,
**read-only to start**, with the permission model left open to granting them
actions later. Every audited action records its actor, automated ones included.

> **Dependency:** restitution requires the audit journal be *replay-complete*,
> which is a materially stronger requirement than logging for moderation review.
> The journal's completeness must be designed against restitution, not against
> the review queue.

### 9. Control-plane API

The server **owns the control plane** — world cloning, content promotion,
moderation, and operational actions are an authenticated API. `admin-tools` is
purely a client, as are scripted operators, CI, and automated moderators.

The control API contract is **defined and published by this repository**,
separate from the game protocol in `proto`, because the two change at different
rates. The server advertises its engine version and control-API semver; clients
declare a supported range and **refuse to operate on mismatch** with a clear
message.

**Style: REST + OpenAPI.** Chosen for consumer accessibility — `curl` works, any
language works, and an OpenAPI spec is directly consumable by third-party tools
and automated moderators with no protobuf toolchain required. Paths are versioned
(`/v1/...`) and **the spec is generated from code** so it cannot drift from the
implementation.

Connect (gRPC + gRPC-Web + JSON from one handler) must remain reachable later,
which imposes four constraints — cheap from the first handler, expensive to
retrofit:

- **Handlers are a thin adapter, never the logic.** The control plane is a Go service layer taking request structs and returning response structs; `net/http` types never appear in its signatures. A Connect handler becomes a second adapter over the same layer.
- **Operations are method-shaped, not deeply RESTful.** Prefer explicit actions (`POST /v1/worlds/{id}/promote`) over verb overloading or nested-resource cleverness, one request object per operation. HTTP semantics carry transport meaning only, never business meaning.
- **Errors are domain values.** Every failure has a stable machine-readable code, mapped to an HTTP status by the adapter. Handlers never invent status codes.
- **Progress is a typed event stream** in the service layer, delivered over SSE today; Connect server-streaming would be another adapter over it.

Content bundle upload is expected to remain a plain HTTP route regardless, since
large multipart and streaming bodies have no Connect equivalent worth adopting.

Expensive processing does **not** belong here. Clients do the heavy authoring-time
work — script type-checking, reference resolution, map baking, linting, reporting
— and send results. The server still performs cheap authoritative structural
validation on everything received: offloading is a performance and UX measure,
never a trust boundary.

**Exception — anonymization.** Redaction happens server-side at export time, off
the tick path, because doing it client-side would mean emitting raw credentials
and PII to operator machines. Export snapshots the database (`VACUUM INTO`) and
streams redacted output from the copy, avoiding contention with live traffic.

### 10. Operational facilities

- **Bounded background job runner** — context cancellation, progress reporting, concurrency limits, explicit throttling and bounded memory and disk. Shared by deferred script work, exports, clones, and migrations. A bare goroutine is insufficient: a saturated core will blow the tick budget on a small VPS, and allocation-heavy work raises GC cost process-wide.
- **Separate listeners** for game traffic and the control plane. The control listener binds to `127.0.0.1` by default; exposing it is an explicit operator choice.
- **First-run bootstrap** as CLI output — create the owner account, emit a token — since there is no embedded UI.
- Configuration, health checks, metrics, graceful shutdown, schema migration.

## Explicitly out of scope

| Concern | Owner |
|---|---|
| Rendering, presentation, input handling | `client` |
| Game protocol schema and generated bindings | `proto` |
| Operator UI, CLI tooling, and all heavy processing | `admin-tools` |
| Cross-server discovery and identity anchoring | `registry` |
| Game content, assets, and world design | World authors |

The engine also does not ship playable content, and does not implement
branching or merging of content revisions.

## Open questions

Decision history and rationale: [`decisions.md`](./decisions.md).

- **Tick rate** — 10 Hz is provisional pending live-load testing.
- **Interest management** — the algorithm for deciding what each client is told about.
- **Registry SSO** — token format and trust establishment.
- **Content manifest** — how content declares its required engine version range.
- **Deferred by choice** — capability-based permissions, direct Git push, automated moderation actions, Postgres, a second script runtime, non-WebSocket transports, Connect as a control-plane protocol.
