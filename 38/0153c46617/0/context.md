# Session Context

## User Prompts

### Prompt 1

<teammate-message teammate_id="team-lead">
You are working on the hetzner-blackbsd Go project at /home/omar/sandbox/blackBSD/bsd-hcloud.

## CRITICAL RULES
- DO NOT modify .golangci.yml
- DO NOT add //nolint comments or any magic comments  
- Write clean, idiomatic, DRY Go code
- No dead code, no smells
- Use external test packages (package foo_test)
- ALL linters enabled (exhaustruct, varnamelen, dupl, paralleltest, errcheck, etc.)
- Run `task ci` after ALL changes to verify

## YOUR TASKS

###...

### Prompt 2

<teammate-message teammate_id="team-lead" summary="Modular structure and CI verification guidelines">
IMPORTANT GUIDELINES FROM TEAM LEAD:

1. MODULAR FILE STRUCTURE: Do NOT make huge files. Break code into separate files by concern. E.g., for hcloud: put power actions in a separate file (power.go), SSH key ops in sshkeys.go, wait helpers in wait.go. For SSH: put SFTP in sftp.go, PTY in pty.go.

2. REUSABLE MODULES: Design code so it can be reused. Use interfaces, not concrete types.

3. UTILS: ...

### Prompt 3

<teammate-message teammate_id="team-lead" summary="New task: create NetBSD installer module">
Great work on SFTP and PTY! Here's your next task:

## Task: Create internal/netbsd package - NetBSD installer module

### Files to Create
- `internal/netbsd/installer.go` - Main installer struct and methods
- `internal/netbsd/installer_test.go` - Tests

### Context
This module automates NetBSD installation via QEMU inside Hetzner rescue mode. It downloads the NetBSD ISO, runs QEMU with it, and automate...

