# BUG_REPRO

## Bug 是什么
Store.New 没有初始化 subs map，首次 AddSubscription 向 nil map 写入。

## 如何触发
`go test ./...`（或 `go test -run TestAddGetSubscription ./internal/store/`）

## 错误信息
- panic: assignment to entry in nil map
