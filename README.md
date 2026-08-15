# forex

一个用 Go 写的内存汇率服务，演示分层、批量查询、链式换算、订阅通知与并发 worker 池。

## 功能
- 批量查询汇率、支持货币列表、直接/链式货币换算
- 汇率订阅：目标价触发后生成通知
- 并发订阅检查 worker 池，支持 context 取消

## 目录结构
```
cmd/forex/          程序入口
internal/config/    环境配置
internal/model/     模型与纯工具函数
internal/store/     内存存储（汇率 + 订阅 + 锁）
internal/service/   业务逻辑
internal/worker/    订阅检查 worker 池
```

## 运行与测试
```bash
go build ./...
go test ./...
```
