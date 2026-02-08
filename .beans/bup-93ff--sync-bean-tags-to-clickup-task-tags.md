---
# bup-93ff
title: Sync bean tags to ClickUp task tags
status: in-progress
type: feature
priority: normal
created_at: 2026-01-30T18:55:15Z
updated_at: 2026-02-08T23:11:25Z
extensions:
    clickup:
        synced_at: "2026-02-08T22:25:52Z"
        task_id: 868ha7twn
---

Bean tags are now defined and used in pacer/core but aren't synced to ClickUp. Add tag sync support.

## Context

- ClickUp tags are space-level, managed via separate endpoints (not part of create/update task body)
- `POST /api/v2/task/{task_id}/tag/{tag_name}` — add tag (one at a time)
- `DELETE /api/v2/task/{task_id}/tag/{tag_name}` — remove tag (one at a time)
- Tags auto-create in ClickUp space when first applied to a task
- Bean `Tags` field is already parsed in `internal/beans/types.go` but ignored during sync

## Checklist

- [x] Add `Tags` field to `TaskInfo` and `taskResponse` in `internal/clickup/types.go` (ClickUp returns tags on GET task)
- [x] Add `AddTagToTask(ctx, taskID, tagName) error` method to `internal/clickup/client.go`
- [x] Add `RemoveTagFromTask(ctx, taskID, tagName) error` method to `internal/clickup/client.go`
- [x] Add tag sync logic in `internal/clickup/sync.go` after task create/update:
  - Compare bean tags with current ClickUp task tags
  - Add missing tags, remove extra tags
  - Best-effort (log but don't fail sync on tag errors)
- [x] Add tests for tag sync logic
- [ ] Test manually: `beanup sync --force` on a bean with tags, verify tags appear in ClickUp

## Sync Logic

```
currentTags := set of tags on ClickUp task (from GET response)
desiredTags := set of tags on bean

toAdd := desiredTags - currentTags
toRemove := currentTags - desiredTags

for each tag in toAdd: POST /task/{id}/tag/{tag}
for each tag in toRemove: DELETE /task/{id}/tag/{tag}
```

## No config changes needed

Tag names pass through directly (bean tag name = ClickUp tag name). No mapping config required.
