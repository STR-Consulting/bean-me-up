---
# bup-kcnt
title: Fix GraphQL syntax error in SetExtensionData
status: completed
type: bug
priority: normal
created_at: 2026-02-08T23:10:28Z
updated_at: 2026-02-08T23:56:53Z
extensions:
    clickup:
        synced_at: "2026-02-08T23:56:52Z"
        task_id: 868hdrpb9
---

The SetExtensionData and SetExtensionDataBatch methods in internal/beans/client.go use json.Marshal to format the data parameter, producing JSON with quoted keys (e.g. {"task_id": "123"}). GraphQL input object literals require unquoted keys (e.g. {task_id: "123"}). This causes a 'Expected Name, found String' parse error when calling the setExtensionData mutation.

Fix: replace JSON marshaling with a helper that produces GraphQL input object literal syntax.
