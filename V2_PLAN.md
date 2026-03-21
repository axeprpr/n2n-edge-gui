# n2nGUI 2.0 Plan

## Target

Build a maintainable `n2n` desktop client with a modern UI, clearer process
control, and a codebase that can evolve beyond Windows-only tray management.

## New Name

- Repository name: `n2nGUI`
- App name: `n2nGUI`
- Windows binary name: `n2nGUI.exe`

This removes the current name's platform lock-in and leaves room for Linux and
macOS clients later.

## Product Direction

The current app is only a thin wrapper around `edge.exe`. Version 2.0 should be
positioned as a desktop client with three responsibilities:

1. Manage local `n2n edge` configuration.
2. Start, stop, and observe the local edge process.
3. Present connection state and logs in a usable UI.

It should not attempt to reimplement the `n2n` protocol itself.

## UI Direction

The current `walk` tray UI is too limited and Windows-specific. For 2.0, use a
two-layer design:

1. Core service layer in Go.
2. Cross-platform desktop shell on top.

Recommended UI stack:

- `Wails` for the desktop shell
- Go backend for process/config logic
- Web UI frontend for the desktop experience

Why this route:

- Better cross-platform path than `walk`
- Easier to build a richer status page, logs view, onboarding flow, and settings
- Cleaner separation between UI and process logic

Fallback option:

- `Fyne` if the goal is a pure-Go desktop app with a simpler UI surface

## 2.0 UI Scope

Primary screens:

1. Dashboard
   - Running/stopped state
   - Connected community
   - Supernode target
   - Last start time
   - Last error

2. Node Configuration
   - Community name
   - Address mode: static IP or DHCP
   - Static address field
   - Supernode host and port
   - MTU
   - Advanced args

3. Runtime Logs
   - Process stdout/stderr
   - Copy/export logs
   - Clear log view

4. Diagnostics
   - Binary path
   - TAP/TUN availability
   - Current platform
   - Process PID
   - Config file path

5. Settings
   - Auto start
   - Minimize to tray
   - Start on login
   - Theme
   - Binary management

## Architecture

Suggested layout:

```text
cmd/
  n2nGUI/
internal/
  app/
  config/
  edge/
  diagnostics/
  platform/
  tray/
frontend/
  src/
  assets/
build/
docs/
```

Responsibilities:

- `internal/config`: load, validate, save config
- `internal/edge`: build args, start/stop process, stream logs
- `internal/platform`: OS-specific helpers
- `internal/diagnostics`: environment and adapter checks
- `internal/app`: application service layer exposed to the UI
- `internal/tray`: optional tray integration per platform

## Core Functional Changes

The current implementation should be replaced in these areas:

1. Process lifecycle
   - Stop using global `taskkill /im edge.exe /f`
   - Track the process started by the app
   - Capture stdout/stderr
   - Report start failures precisely

2. Config model
   - Move away from free-form `conf.ini` editing only
   - Use a typed config struct
   - Validate before save
   - Generate final edge args from the struct

3. Platform handling
   - Remove hardcoded `cmd`, `tasklist`, `taskkill`, and `.exe` assumptions
   - Introduce platform adapters for Windows, Linux, and later macOS

4. Packaging
   - Separate source, bundled runtime binaries, and installers
   - Keep release artifacts out of the main source tree where possible

## Migration Plan

### Phase 1: Core Extraction

- Add `go.mod`
- Extract config reading/writing out of `main.go`
- Extract edge process management into a dedicated package
- Add basic unit tests for config and arg generation

Deliverable:

- A CLI or service layer that can load config and start/stop `n2n edge`

### Phase 2: Rename and Repo Cleanup

- Rename project references from `n2nEdgeWindowsGui` to `n2nGUI`
- Rename build output to `n2nGUI.exe`
- Move installer assets and third-party binaries under a dedicated release or
  packaging directory
- Write a real README

Deliverable:

- Cleaner repository with explicit ownership of code vs packaged assets

### Phase 3: New Desktop UI

- Scaffold a `Wails` desktop app
- Build Dashboard, Configuration, and Logs screens
- Wire UI actions to the Go backend service
- Preserve tray behavior as a secondary feature, not the primary interface

Deliverable:

- Functional `n2nGUI 2.0` desktop client

### Phase 4: Cross-Platform Enablement

- Implement platform-specific binary/path/process adapters
- Support Linux packaging
- Assess macOS support separately due to networking and permissions constraints

Deliverable:

- Shared UI with platform-specific runtime handling

## Risks

- `n2n` virtual adapter setup is platform-specific and may require admin/root
  flows that differ sharply across OSes.
- A cross-platform GUI is straightforward; a cross-platform installer and
  virtual network setup is not.
- Bundling old `n2n` binaries inside the repo will keep maintenance expensive
  unless release management is cleaned up.

## Immediate Next Work

Recommended next implementation order:

1. Add `go.mod` and stabilize builds.
2. Split current `main.go` into config and process packages.
3. Replace global process kill behavior with tracked child process management.
4. Introduce a typed config model and validation.
5. Only then build the new UI shell.

## Versioning Suggestion

- Current codebase: `legacy` or `v1`
- New branch or milestone: `v2`
- First target release: `n2nGUI 2.0.0-alpha1`
