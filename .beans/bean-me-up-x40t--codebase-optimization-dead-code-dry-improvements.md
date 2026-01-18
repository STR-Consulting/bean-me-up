---
# bean-me-up-x40t
title: Codebase optimization - dead code & DRY improvements
status: completed
type: task
priority: normal
created_at: 2026-01-18T01:37:59Z
updated_at: 2026-01-18T01:43:49Z
---

Remove dead code, inline over-abstractions, and consolidate repeated patterns.

## Checklist

### Phase A - Dead Code Removal
- [x] Delete `internal/frontmatter/` directory entirely
- [x] Remove unused methods from `internal/beans/types.go` (GetClickUpTaskID, GetClickUpSyncedAt, NeedsClickUpSync)
- [x] `GetPriorityMapping()` is NOT unused (called in check.go) - skipped

### Phase B - Over-Abstraction Cleanup
- [x] Inline `Default()` into `Load()` in config.go
- [x] `FindConfig()` has external caller (check.go) - kept as-is
- [x] Inline `FindBeansConfig()` into `LoadBeansPath()` in config.go
- [x] Remove `GetStatusMapping()` wrapper from clickup/sync.go

### Phase C - DRY Consolidation
- [x] Add `outputJSON()` helper to cmd/root.go
- [x] Add `toTaskInfo()` method to clickup/types.go (on taskResponse)
- [x] Add `requireListID()` helper to cmd/root.go
- [x] Add `loadSyncStore()` helper to cmd/root.go

### Verification
- [x] Build compiles successfully
- [x] All tests pass
- [x] golangci-lint passes