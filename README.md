# proxy_server — an HTTP reverse proxy in Go (from scratch)

A reverse/forwarding HTTP proxy written in Go using **only the standard library** — no proxy framework. Built up in stages (`part1` → `part5`) to explore Go networking and systems programming from first principles.

## What it does

- **Reverse proxying** — forwards incoming HTTP requests to backend servers.
- **Connection pooling** — reuses backend TCP connections via a buffered channel, dialing new ones on demand.
- **Traffic stats** — tracks per-backend request bytes behind a `sync.Mutex` for safe concurrent access.
- **RPC stats server** — exposes a `GetStats` endpoint over Go's `net/rpc` so the proxy's metrics can be queried by a separate client (`part5/rpc_client.go`).

## Layout

```
webserver.go        RPC stats server (net/rpc over HTTP)
part1..part5/       progressive build-up of the proxy + an RPC client in part5
```

## Stack

**Go** standard library only — `net`, `net/http`, `net/rpc`, `bufio`, `sync`.

## Running

```bash
go run part5/proxy_server.go    # or any stage you want to inspect
```

> A learning project focused on the mechanics of proxying, connection pooling, and inter-process stats over RPC — no third-party networking libraries.
