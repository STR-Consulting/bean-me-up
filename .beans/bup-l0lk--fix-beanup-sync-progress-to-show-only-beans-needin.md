---
# bup-l0lk
title: Fix beanup sync progress to show only beans needing sync
status: completed
type: task
priority: normal
created_at: 2026-01-18T21:00:36Z
updated_at: 2026-02-08T23:56:53Z
extensions:
    clickup:
        synced_at: "2026-02-08T23:56:53Z"
        task_id: 868h4kb92
---

Pre-compute which beans need syncing by comparing timestamps upfront, then only process and show progress for those beans.
