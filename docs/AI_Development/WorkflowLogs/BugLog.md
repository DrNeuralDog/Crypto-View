# Bug Log - CryptoView

Формат записи: `[<Timestamp>] <Brief description of error/debug> - <Result>`

---
[2026-02-20 11:26:52] NotebookLM query failed: authentication expired (nlm login required) during Phase 2.5 planning check - Local implementation continued; waiting for re-auth to resume NLM queries.
[2026-02-20 11:27:31] NotebookLM log sync command returned authentication-expired traceback after execution - Sync status cannot be trusted without fresh nlm login.
[2026-02-20 19:03:51] /sync_logs failed: nlm source list returned authentication expired for notebook 7fa1e693-e081-4415-832c-69c96d1350bb - Blocked until nlm login re-authentication.
[2026-02-20 20:05:20] go test failed during UI theming integration: wrong widget.List.SetItemHeight signature and canvas.Circle.SetMinSize usage - Fixed by removing SetItemHeight call and wrapping status indicator with container.NewGridWrap.
[2026-02-20 20:08:15] Encoding/test mismatch during currency symbol update caused invalid UTF-8 and wrong RUB expectation in coin_list_test - Resolved by using Unicode escape sequences (\u20ac, \u20bd) and updating tests.
