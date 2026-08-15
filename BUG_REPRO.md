# BUG_REPRO

## Bug 是什么
ValidSubscription 判断方向写反、Subscribe 丢掉了入参校验、MarkNotified 去掉了锁、worker 汇总去锁且 wg.Add 放进了 goroutine，导致通知重复触发、统计错误，并存在数据竞争；无效订阅也能被提交。

## 如何触发
`go test -race ./...`

## 错误信息
- TestValidSubscription / TestSubscribeRejectsInvalid / TestCheckSummary 失败。
- `go test -race` 报告 data race（WaitGroup 误用 + 汇总并发写）。
