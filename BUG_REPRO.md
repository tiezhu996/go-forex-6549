# BUG_REPRO

## Bug 是什么
worker.Check 的生产者 goroutine 在 context 已取消时显式 close(ch)，同时 defer close(ch)，channel 被关闭两次。

## 如何触发
`go test -run TestCheckCancel ./internal/worker/`

## 错误信息
- panic: close of closed channel
