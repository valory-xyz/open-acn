# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

This repo is `libp2p_node`, the Go implementation of a node for the Agent Communication Network (ACN). ACN lets agents (see `open-aea`) discover each other and exchange messages addressed solely by wallet address. The binary is a single entry point (`libp2p_node.go`, `package main`, module `libp2p_node`) that wires together the AEA-side pipe and the libp2p-based DHT.

## Common commands

```bash
make build       # go build
make test        # full test suite with coverage (serial: -p 1, no inlining, no test cache)
make race_test   # same, with -race
make lint        # golines . -w && golangci-lint run
make install     # go get -v -t -d ./...
make clean       # remove binary + coverage.txt
```

Run a single test:
```bash
go test -gcflags=-l -count=1 -v ./dht/dhtpeer -run TestRoutingAllToAllConnectClientsSamePeer
```
`-gcflags=-l` disables inlining (required by `bou.ke/monkey`), and `-p 1` in `make test` forces serial package execution because tests bind fixed ports — keep those flags when running tests that touch networking.

`run_acn_node_standalone.py` launches the built binary standalone (outside an AEA), reading env-file config; `--config-from-env` on the binary reads config from environment variables instead of an AEA pipe. See README for bootstrap/entry node env-file examples.

## Architecture

The node can run in one of two modes, selected in `libp2p_node.go:main` by whether a public URI is set:

- **Peer mode** (`dht/dhtpeer`) — full libp2p host. Runs the DHT, optionally a relay service, a delegate service (TCP endpoint for `p2p_libp2p_client` connections), an optional mailbox service, and optional Prometheus monitoring. Constructed via functional options (`dht/dhtpeer/options.go`).
- **Client mode** (`dht/dhtclient`) — lightweight, bootstraps from entry peers, no public address. Used when `AEA_P2P_URI_PUBLIC` is unset.

Both modes implement the `dhtnode.DHTNode` interface (`dht/dhtnode/dhtnode.go`): `RouteEnvelope`, `ProcessEnvelope`, `MultiAddr`, `PeerID`, `Close`. `main` connects the AEA pipe to the node by (a) forwarding envelopes from `agent.Queue()` into `node.RouteEnvelope`, and (b) registering `agent.Put` as the node's inbound envelope handler.

Key packages:

- `aea/` — the AEA-facing API. `api.go` handles config (env or pipe-based), `envelope.proto`/`envelope.pb.go` is the wire format, `pipe.go` is the Unix-pipe transport to a co-located Python AEA.
- `dht/dhtnode/` — shared stream protocol handlers and the `DHTNode` interface. The ACN libp2p protocol IDs (`/aea/0.1.0`, `/aea-address/0.1.0`, `/aea-register/0.1.0`) are defined here.
- `dht/dhtpeer/` — peer implementation; `mailbox.go` is the store-and-forward mailbox service, `notifee.go` hooks libp2p connection events, `benchmarks_test.go` contains throughput benchmarks.
- `dht/dhtclient/` — client implementation.
- `dht/dhttests/` — shared test fixtures/harness (imported by other `dht/*` tests).
- `dht/common/`, `dht/monitoring/` — shared helpers and the Prometheus monitoring abstraction.
- `acn/`, `protocols/`, `utils/`, `common/` — ACN-level helpers, generated protocol messages, logging and crypto utilities.
- `mocks/` — gomock-generated mocks (see https://github.com/golang/mock for regeneration).

### Messaging patterns

ACN supports several delivery paths between an AEA `Connection` and a `Peer`, via optional `Delegate Client` and `Relay Peer` hops — see README.md "Messaging patterns" for the full matrix. ACN must guarantee total ordering of messages between any pair of agents regardless of which path is used.

## Audits

Security & correctness audits live under [`audits/`](./audits/). Most recent: [`audits/AUDIT-2026-04-15.md`](./audits/AUDIT-2026-04-15.md). When changing code in `aea/pipe.go`, `dht/common/handlers.go`, `dht/dhtpeer/{dhtpeer,mailbox}.go`, `dht/monitoring/`, or `utils/utils.go`, cross-check the audit's open findings before merging — many of the Critical/High items live in those files.

## Go / tooling notes

- Go module: `libp2p_node` (go 1.17). Internal imports use the module path, e.g. `libp2p_node/dht/dhtpeer`.
- Dependencies are pinned to older libp2p (`go-libp2p v0.8.3`, `go-libp2p-core v0.5.3`, `go-libp2p-kad-dht v0.7.11`) — do not casually bump these; the DHT protocol and stream APIs differ substantially in newer versions.
- `golines` reformats long lines as part of `make lint`; run it before committing Go changes.
