# AGENTS.md — merlin-cli

Guidance for AI coding agents working in this repo. Human contributors may find it useful too.

## What this is

The **Merlin command-line client**. Go module `github.com/Ne0nd0g/merlin-cli` (no `/v2`). An
interactive readline TUI that talks to the Merlin server over gRPC, enabling multi-user operation.

Target Go: **1.27** (revival target; drop the stale `toolchain` directive). Builds clean on Go 1.26.

## Build / run

```bash
go build ./...
go build -o merlinCLI .
./merlinCLI -password <pw>      # connect to a server (default 127.0.0.1:50051)
```

Connection flags: `-password`, `-secure`, `-tlsKey`, `-tlsCert`, `-tlsCA` (see `main.go`).

## gRPC client facts

- The generated stubs live in [`rpc/`](rpc/) (`rpc.pb.go`, `rpc_grpc.pb.go`) and are **generated from
  the server repo's `pkg/rpc/rpc.proto`**. This repo does not import the server as a Go module. After
  the server's proto changes, regenerate these stubs and **keep `google.golang.org/grpc` and
  `google.golang.org/protobuf` versions aligned with the server** (they are currently skewed — the
  server is older; align during the revival).
- The client **always dials TLS**; when not in secure mode it uses an `InsecureSkipVerify` config
  (the server serves a self-signed cert). Auth is an `authorization: <password>` metadata header
  added by a unary + stream interceptor (`services/rpc/rpc.go`).
- The TUI is interactive and not easily scriptable; the end-to-end smoke test in the server repo
  (`../merlin/test/smoke/`) drives the gRPC API directly instead of this CLI.

## Cross-repo relationship

`merlin-cli` ← (generated from) ← `merlin`'s `rpc.proto`. No Go import edge, but a **contract edge**:
regenerate + realign versions whenever the server's proto moves. Release order: server → cli.

## Conventions

- **Branches:** do all work on `dev` (or a feature branch). **Never commit to `main`.**
- **Commits:** the maintainer signs every commit with a YubiKey. **Do not run `git commit`** —
  stage changes and propose a commit message. Do **not** add a `Co-Authored-By` trailer.
- Match surrounding style; keep the GPLv3 license header on new Go files.
