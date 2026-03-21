# n2nGUI Desktop Shell

This directory contains the Wails desktop wrapper for the `n2nGUI 2.0`
prototype.

It reuses the same API and visual layer as the embedded web build:

- static assets are served from `desktop/frontend`
- runtime API calls are forwarded through the Wails asset handler
- `n2n` binaries and config are resolved from the repository root

## Why this shape

The web prototype was already split into:

- typed config
- process manager
- HTTP API
- polished frontend

Wails is added here only as the desktop container, not as a second backend.

## Expected next step

Install Wails and its Go dependency, then run the desktop shell from this
directory.

The Go entrypoint is [main.go](/root/github-fork-review/repos/n2nEdgeWindowsGui/desktop/main.go).
