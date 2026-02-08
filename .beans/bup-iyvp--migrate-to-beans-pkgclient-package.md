---
# bup-iyvp
title: Migrate to beans pkg/client package
status: completed
type: task
priority: normal
created_at: 2026-02-08T23:28:55Z
updated_at: 2026-02-08T23:56:53Z
extensions:
    clickup:
        synced_at: "2026-02-08T23:56:52Z"
        task_id: 868hdrpbb
---

Replace the hand-rolled GraphQL client in internal/beans/client.go with the new github.com/hmans/beans/pkg/client package. This eliminates the raw GraphQL string construction that caused the quoted-keys bug.

## Checklist

- [x] Add github.com/hmans/beans dependency
- [x] Update internal/beans/client.go to wrap or replace with pkg/client
- [x] Update call sites (migrate command, external_sync.go)
- [x] Remove toGraphQLLiteral helper (no longer needed)
- [x] Verify build and tests pass
- [x] Run beanup migrate successfully
