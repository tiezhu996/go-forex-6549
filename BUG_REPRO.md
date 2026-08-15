# BUG_REPRO

## Bug 是什么
BuildSubscriptionBatches 返回共享底层数组的子切片、SubscriptionBatches 手写共享切片、RatePairs 直接返回内部 order、worker 把每个 batch 截掉最后一个元素，导致返回列表被外部改动污染、并发检查时订阅被漏掉。

## 如何触发
`go test ./...`

## 错误信息
- TestBuildSubscriptionBatchesFresh / TestRatePairsFresh / TestCheckSummary 失败。
