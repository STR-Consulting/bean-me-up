---
# bup-u5jy
title: Constrain type mapping to standard bean types
status: completed
type: feature
priority: normal
created_at: 2026-01-18T17:53:07Z
updated_at: 2026-02-08T23:11:25Z
extensions:
    clickup:
        synced_at: "2026-01-18T17:55:20Z"
        task_id: 868h4jued
---

Change the type mapping configuration so it only accepts the 5 standard bean types (milestone, epic, feature, bug, task) rather than arbitrary strings.

## Checklist
- [x] Define standard bean types in internal/beans/types.go
- [x] Add validation to config loading in internal/config/config.go
- [x] Update example config .beans.clickup.yml.example
- [x] Update beanup types output in cmd/types.go
- [x] Verify build, tests, and lint pass
