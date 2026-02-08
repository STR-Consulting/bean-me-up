---
# bup-bhc6
title: Investigate beanup sync syncing too many beans
status: completed
type: bug
priority: normal
created_at: 2026-01-18T20:51:39Z
updated_at: 2026-02-08T23:11:25Z
extensions:
    clickup:
        synced_at: "2026-01-19T19:05:53Z"
        task_id: 868h4kb90
---

When running `beanup sync` after changing one bean, it syncs all 467 beans instead of just the changed one. Need to investigate why change detection isn't working properly.
