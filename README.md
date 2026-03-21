# n2nGUI

`n2nGUI` is a 2.0-oriented control plane for `n2n edge`.

The repository still contains the old Windows tray utility, but it now also
includes a new prototype with:

- a typed config model
- tracked edge process lifecycle
- live diagnostics and buffered logs
- a polished web UI embedded in a Go server

## Prototype Run

```bash
go build ./cmd/n2nGUI
./n2nGUI -listen 127.0.0.1:8787
```

Open `http://127.0.0.1:8787`.

The server expects the `n2n` runtime binaries under `./n2n/`.

## Desktop Shell

A Wails desktop wrapper now lives in
[desktop/main.go](/root/github-fork-review/repos/n2nEdgeWindowsGui/desktop/main.go).

It reuses the same API and UI instead of adding a second application stack.

## Current Layout

- `cmd/n2nGUI`: entrypoint and embedded UI
- `internal/config`: typed config model, validation, legacy `conf.ini` import
- `internal/edge`: child process management and log capture
- `internal/app`: HTTP API for config, status, logs, and diagnostics

## Notes

- The new prototype is cross-platform at the Go service layer.
- Actual adapter installation and OS packaging are still platform-specific.
- See [V2_PLAN.md](/root/github-fork-review/repos/n2nEdgeWindowsGui/V2_PLAN.md)
  for the broader migration plan.
