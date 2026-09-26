# server

The Odyssey engine — the authoritative server that runs persistent, shared
worlds. It simulates the world, holds authoritative state, enforces the rules
of play, and persists worlds over time.

**In scope:** authoritative world state and simulation, game rules and systems,
player sessions and authentication, world persistence and migration, in-world
moderation and safety tooling, server-side enforcement of fairness.

Full detail: [`docs/scope.md`](./docs/scope.md).

**Out of scope:** presentation, rendering and input handling (see `client`), the
shared protocol contract (see `proto`), operator and host-facing tooling (see
`admin-tools`), cross-server identity and discovery (see `registry`).

**License:** AGPL-3.0 — the server is the system that makes a world run, and
network copyleft is what keeps the commons guarantee real for hosted software.
