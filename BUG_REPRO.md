# BUG_REPRO

## Bug 是什么
store 层用 fmt.Errorf 替换了哨兵错误、service 层把 %w 改成 %v、MergeSummary 漏掉 Failed 累加、worker 未累计 Failed，导致 errors.Is 判断失效、失败数统计丢失。

## 如何触发
`go test ./...`

## 错误信息
- TestMergeSummary / TestConvertChain / TestSubscribeRejectsInvalid / TestPutGetRate / TestMarkNotified 失败。
