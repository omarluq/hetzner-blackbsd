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
- The golangci-lint c...

### Prompt 2

<teammate-message teammate_id="team-lead" summary="Modular structure and CI verification guidelines">
IMPORTANT GUIDELINES FROM TEAM LEAD:

1. MODULAR FILE STRUCTURE: Do NOT make huge files. Break code into separate files by concern. E.g., for hcloud: put power actions in a separate file (power.go), SSH key ops in sshkeys.go, wait helpers in wait.go. For SSH: put SFTP in sftp.go, PTY in pty.go.

2. REUSABLE MODULES: Design code so it can be reused. Use interfaces, not concrete types.

3. UTILS: ...

### Prompt 3

This session is being continued from a previous conversation that ran out of context. The summary below covers the earlier portion of the conversation.

Analysis:
Let me chronologically analyze the conversation:

1. **Initial Setup**: I'm working as a teammate agent (`customize-module`) in a team, assigned by a team lead to create an `internal/customize` package for the hetzner-blackbsd Go project.

2. **Team Lead's First Message**: Detailed task to create `internal/customize` package with:
   -...

