# Engine Decisions

Decisions taken for the `server` engine, with the reasoning behind them. This is
the record to consult before revisiting a choice, and the source material for
author- and operator-facing documentation.

Scope of responsibilities: [`scope.md`](./scope.md).

Format: **D<N> (date) — decision.** Why, and what it costs.

---

## Shape of the engine

**D1 (2026-09-25) — The engine is the product; Odyssey Classic is its first world.**
Its capability envelope is set by what earlier Odyssey servers could do, but none
of their content ships with it.

**D2 (2026-09-25) — Hybrid: core systems in Go, content as data, behaviour overridable by scripts.**
Rules in code are faster to build and reason about than a fully data-driven
engine; scripts supply the extensibility that pure-code rules would deny world
authors.

**D3 (2026-09-25) — Many self-hosted worlds, discovered via `registry`.**
Drives the AGPL choice, first-class operator concerns, and zero-ops defaults.

**D4 (2026-09-25) — Target 100–500 concurrent players per world.**
Single process, no sharding, no distributed state. Enough to make simple choices
affordable and cheap enough to keep them honest.

**D5 (2026-09-25) — Tile-grid authority with client-side interpolation; clients predict self, interpolate others.**
The genre-standard compromise and what players expect. Server remains the only
authority.

**D6 (2026-09-25) — WebSocket + protobuf, with the message layer kept transport-agnostic.**
Works everywhere a browser does. Head-of-line blocking is acceptable at
tile-movement rates; datagrams stay possible later.

**D7 (2026-09-25) — 10 Hz simulation, configurable.**
Ample for tile movement with interpolation, and leaves script budget inside the
frame. Provisional until validated under live load; systems must be written
rate-correct rather than assuming 10 Hz.

## Identity and permissions

**D8 (2026-09-25) — Each world owns local accounts; registry SSO is opt-in.**
A world must work standalone and offline. Cost: two auth paths to build and
secure.

**D9 (2026-09-25) — Fixed roles now (owner, operator, content author, moderator, player), behind a capability indirection.**
Covers real worlds immediately; granular capabilities can land later without
breaking callers.

**D10 (2026-09-25) — Pre-authentication work is rate-limited and cost-capped per IP.**
Password hashing is deliberately expensive, which makes unauthenticated login an
amplifier if left unbounded.

## Scripting

**D11 (2026-09-25) — JavaScript/TypeScript via goja.**
The widest pool of authors, real types and editor support, hot reload, honest
stack traces. WASM was rejected: its ABI design and compile step cost exactly the
iteration speed the project exists to gain, buying isolation that D12 says is not
needed.

**D12 (2026-09-25) — Scripts are semi-trusted.**
Operator-authored, so not assumed hostile, but bounded: per-call execution
budgets, no filesystem or network unless injected, errors contained and logged
rather than fatal. goja's lack of ambient I/O makes this the default rather than
an additional mechanism.

**D13 (2026-09-25) — Scripts may hook and modify events, replace whole systems, create runtime content into layer B, and keep their own durable namespaced state.**
"Replace whole systems" requires engine systems be built behind script-visible
seams from the start — retrofitting seams is expensive.

**D14 (2026-09-25) — Gameplay hooks run synchronously in-tick; slow work is deferred off-tick.**
Keeps world state consistent and reasoning simple while denying a slow script the
ability to degrade the tick for everyone.

**D15 (2026-09-25) — No per-concept timer systems. One persisted scheduled-wakeup queue.**
Gated checks ("may this be claimed yet?") need no timer at all — a stored
timestamp compared against now, evaluated lazily. Expiry needs a push, and one
`(fire_at, handler, payload)` queue serves every case: respawns, buff expiry,
growth. Replaces what would otherwise be a timer subsystem per game concept.

**D16 (2026-09-25) — Missed wakeups fire on catch-up at boot, and the handler receives both `scheduled_for` and `now`.**
Only the script can know whether a late fire should still act: a boss respawn
should, a stale daily reset should not. Overdue entries drain in `fire_at` order,
rate-limited so boot is not a stampede. Entries may set `drop_if_late_by`, and
recurring entries coalesce rather than firing once per missed interval.

**D17 (2026-09-25) — Wakeup delivery is at-least-once; handlers must be idempotent.**
**Superseded by D47 (2026-09-26.)** Reasoning at the time: a handler's effects may
span stores the engine cannot commit atomically, so exactly-once was judged
undeliverable in general, and the burden was placed on authors via a stable
`fire_id`. This overstated the limit — exactly-once *is* deliverable for the
subset of effects the database can roll back, which is most wakeup handlers.

**D18 (2026-09-25) — Scripts read time from their invocation context, not an ambient clock.**
One coherent `now` per tick, and scripts testable without clock manipulation.
Durations within a session use monotonic time; cross-restart intervals use wall
clock, which can jump under NTP correction or VM migration. Note: this is *not*
required for replay correctness — see D26.

## Content model and lifecycle

**D19 (2026-09-25) — Three content layers: A authored, B script-authored, C world and player state.**
A flows both ways, B only downstream, C only downstream and only by choice. B
exists because script-created content derives from one server's history and can
never be authored upstream. Keeping scripts out of A is what makes "boot at
revision V" exactly reproducible.

**D20 (2026-09-25) — Layer A is Git-backed, linear per server, with no branching or merging, and Git never surfaces in the UI.**
Git's data model gives exact reproducibility, diffs and rollback; Git's interface
would exclude the non-technical authors the project is trying to serve. Refusing
branch and merge is what keeps this from becoming a version-control project.

**D21 (2026-09-25) — Content ingress is API-only, over a transport-agnostic pipeline.**
Propose revision → validate → commit → publish, with no HTTP assumptions, so
direct Git push can be added later without a second validation implementation.
Two ingress paths with two validators is how content becomes valid through one
door and not the other.

**D22 (2026-09-25) — Clone scope is author-selectable: A, A+B, or A+B+anonymized players. No full clones.**
Each has real use; none justifies putting live credentials on a developer's
laptop.

**D23 (2026-09-25) — Promotion moves layer A only, as a whole revision, atomically, with a preview diff, fast-forward only.**
Atomic whole-revision promotion needs no merge logic and keeps live at a known
version. Fast-forward-only forces a dev re-sync when live has taken a hotfix,
which deliberately puts conflict resolution on the dev server where no outage is
at stake.

**D24 (2026-09-25) — Live content edits are allowed, create a normal revision, and are permission-gated.**
A one-word typo fix should not require a pipeline trip, and live stays at a known
version either way.

**D25 (2026-09-25) — Engine upgrades auto-migrate content forward, non-destructively by invariant.**
Smoothest path for non-technical operators. The risk that a migration bug edits
content is bounded by committing a *new* revision — the prior one stays intact and
rollback-able — and by previewing before applying.

## Persistence

**D26 (2026-09-25) — No custom journal and no snapshot layer. SQLite's WAL is the durability mechanism.**
Once every valuable mutation is a committed transaction, a journal-plus-snapshot
layer only reimplements the database. This also removes any deterministic-replay
requirement from the engine: restitution reads history and applies compensating
transactions rather than re-executing game logic.

**D27 (2026-09-25) — State is either valuable (transactional) or not (memory only). No middle tier.**
The discarded three-tier model created seams where value crossed durability
boundaries, and every such seam is a duplication or vanishing bug.

**D28 (2026-09-25) — Items are a stable UUID plus a single location field.**
`(container_kind, container_id, slot)`. A move is one row update, so an item can
never exist twice or nowhere — duplication becomes structurally impossible rather
than a discipline. A unique constraint on `(container, slot)` gives slot
integrity. Ground drops are just a location and therefore survive restarts. The
simulation caches items in memory write-through; the database is authoritative.

**D29 (2026-09-25) — A stack is one row with a quantity, and currency is an item.**
One uniform model for coins, ammunition and collectibles. Moving 10 of 50 gold is
a decrement and an increment in one transaction, so stacks are not a hot-write
problem. Consequence: the audit ledger needs two entry shapes — UUID relocation
for whole items, `(type, quantity, from, to)` for stack transfers.

**D30 (2026-09-25) — Split and merge are the only non-relocation item operations, and each is a single transaction.**
They mint and destroy UUIDs respectively, making them the only remaining origin
for a duplication bug.

**D31 (2026-09-25) — Death, including item loss and PvP loot transfer, commits as one transaction.**
Prevents the asymmetry where a crash leaves a player alive while a killer already
holds their loot.

**D32 (2026-09-25) — The player anchor (map and tile) persists on transitions only; sub-tile position and facing are memory.**
Worlds are built from small distinct maps, so transitions are frequent enough to
bound how far back a crash sets a player. No heartbeat, no per-step writes.

**D33 (2026-09-25) — A dangling anchor is remapped by a script hook, falling back to world spawn, and reported.**
Content changes can remove the map a player stood on. A bad anchor must never
block a login, and a world usually knows a better answer than world spawn.

**D34 (2026-09-25) — SQLite by default, behind an interface.**
Zero-ops for self-hosters and ample for 500 players; Postgres stays possible
without being built now.

## Audit and moderation

**D35 (2026-09-25) — Audit events are written to a transactional outbox, then shipped to a pluggable append-only sink.**
Committing the audit row with the change it describes means the record can never
disagree with reality; the outbox consumer still allows external sinks and
operator-chosen retention. Cost: one table and one background consumer.

**D36 (2026-09-25) — The engine provides sanctions, live inspection and intervention, a report and review queue, and rollback with restitution.**
Restitution requires the ledger record enough to construct a compensating action
for any value movement — a stronger requirement than logging for human review.

**D37 (2026-09-25) — Automated moderators, including AI agents, are first-class API consumers, read-only to start.**
Actions stay open as a later grant. Every audited action records its actor,
automated ones included.

## Control plane and operations

**D38 (2026-09-25) — The server owns the control plane; `admin-tools` is purely a client.**
Cloning, promotion, moderation and operational actions are an authenticated API,
so scripted operators, CI and automated moderators are first-class too, and
state-mutating logic lives where the state is.

**D39 (2026-09-25) — Control API is REST + OpenAPI, with the spec generated from code.**
Chosen for consumer accessibility: `curl` works, any language works, and an
OpenAPI spec is directly consumable by third-party tools and agent moderators
without protobuf tooling. Generated rather than hand-maintained, because a spec
that drifts is worse than none when clients refuse on version mismatch.

**D40 (2026-09-25) — Connect stays reachable later, under four constraints.**
Handlers are a thin adapter with no `net/http` types in service signatures;
operations are method-shaped rather than deeply RESTful; errors are domain values
with stable codes; progress is a typed event stream, over SSE today. Cheap from
the first handler, expensive to retrofit. Bundle upload is expected to remain a
plain HTTP route, as large multipart bodies have no Connect equivalent worth
adopting.

**D41 (2026-09-25) — The control API contract is defined and published by `server`, not `proto`.**
The control API and the game protocol change at different rates.

**D42 (2026-09-25) — Server advertises engine and control-API versions; clients declare a supported range and refuse on mismatch.**
Predictable failure, and it catches skew before players see it.

**D43 (2026-09-25) — Heavy processing belongs in clients; the server still validates everything it receives.**
Moving script type-checking, reference resolution, map baking and reporting out of
the engine protects the tick better than scheduling them politely inside it. But a
privileged client is still a client: offloading is a performance and UX measure,
never a trust boundary.

**D44 (2026-09-25) — Anonymization is the exception: redaction happens server-side at export.**
Client-side redaction would require emitting raw credentials and PII to operator
machines. Export snapshots the database (`VACUUM INTO`) and streams redacted
output from the copy, avoiding contention with live traffic. A continuously
maintained redacted shadow was rejected: dual writes fail silently, which is the
worst failure shape for a privacy control.

**D45 (2026-09-25) — Single static binary, no embedded UI; control listener binds to localhost by default.**
`admin-tools` is a separate suite of UI and binaries. First-run bootstrap emits
an owner account and token as CLI output. Separate listeners for game traffic and
control plane, and exposing the control listener is an explicit operator choice.

**D46 (2026-09-25) — One bounded background job runner, shared by deferred script work, exports, clones and migrations.**
Context cancellation, progress reporting, concurrency limits, explicit throttling,
bounded memory and disk. A bare goroutine is insufficient: a saturated core blows
the tick budget on a small VPS, and allocation-heavy work raises GC cost
process-wide.

## Wakeup delivery (revision)

**D47 (2026-09-26) — Three delivery modes, defaulting to `exactly_once`.**
`exactly_once` wraps the handler and its dequeue in one transaction;
`at_least_once` may re-run and guards on `fire_id`; `at_most_once` dequeues before
invoking. The discriminator is **whether every effect is rollback-able by the
database**, not handler complexity — a long reward pipeline touching only
persistent state is safe, a three-line chat broadcast is not.

`exactly_once` is the default because it makes the safe thing free and the risky
thing explicit: an author writing a reward handler gets correctness without
knowing the concept exists, while an author doing non-transactional work is told
at the keyboard rather than discovering a duplicated-reward bug in production.
Consistent with D28 (duplication structurally impossible) and D42 (refuse rather
than limp). Cost: `exactly_once` handlers must be short, since SQLite has a single
writer and the handler holds the write lock, must avoid non-transactional side
effects, and cannot use the D14 deferred-work hatch.

**D48 (2026-09-26) — Delivery mode is declared on handler registration, not per schedule call.**
It is a property of what the handler does, and the handler body is fixed. The same
handler scheduled from three call sites with three different guarantees would be a
bug, not a feature. Registration-site declaration also lets the engine validate a
handler's calls against its declared mode consistently. Options ride in a named
bag (`{ delivery }`) rather than a positional flag, alongside the D16 entry
options.

**D49 (2026-09-26) — The `exactly_once` constraint is a hard error from the first release.**
Static detection is not achievable — JavaScript dispatch is too dynamic — so a
non-transactional API called inside an `exactly_once` handler throws at the call
site and its transaction rolls back. No warning period: consistent behaviour from
day one is worth more than gentleness toward early adopters, who are the most
sophisticated users the project will ever have, and a warning that never becomes an
error is a permanent soft failure. Consistent with D42 (refuse rather than limp).
Violations are additionally recorded as content health findings surfaced in
`admin-tools`.

**D50 (2026-09-26) — A throwing wakeup handler is retried a bounded number of times with backoff, then parked.**
A hard error (D49) rolls the transaction back, so the entry stays pending and would
otherwise re-fire on every drain. Bounded retries recover transient failures such
as a lock timeout; the cap prevents a permanently broken handler from looping
forever. Parked entries are dead-lettered, reported in `admin-tools`, and manually
re-triggerable once the script is fixed, so nothing is silently lost and no backlog
grows unseen. Applies to any handler left pending after a failed fire, not only to
constraint violations — an ordinary null dereference has the same effect.

**D51 (2026-09-26) — Interest scope is the map: clients are told about everything on their map and nothing beyond it.**
Sufficient because worlds are built from small distinct maps (the same premise as
D32), putting perhaps 20–40 players on a map rather than 500 in one space. Map
entry sends a full snapshot, then deltas, and map exit discards the set wholesale
— so scope changes only on map transition, which is already a meaningful persisted
event, and the per-entity enter/exit bookkeeping that produces ghost entities in
radius- or line-of-sight schemes never exists.

Two limits are accepted deliberately. **Map occupancy, not total population, is the
scaling limit** — within-map cost grows with the square of occupancy, so crowd hubs
such as markets, capitals and world events are what break this, and occupancy needs
a watched threshold. **Hidden information within a map is impossible** — anything
sent to a client is known to it regardless of what the UI draws, so stealth and
invisibility cannot be implemented by omission and require a deliberate per-entity
exception to scope.

Also the anti-cheat boundary: a map is approximately what a player could legitimately
observe, so per-map scope denies cross-map maphacks by construction while allowing
within-map wallhacks by design.
