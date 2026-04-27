# examples

这里放 go-zero 的独立学习示例。

每个示例尽量保持独立 Go module，方便单独运行、单独理解。

## 示例列表

| 示例 | 类型 | 说明 |
| --- | --- | --- |
| `greet` | API | 最基础的 HTTP API 服务 |
| `user` | RPC | 基于 proto 生成的 user RPC 服务 |
| `user-api` | API | 调用 `user` RPC 的 HTTP API 服务 |

## 推荐学习顺序

1. `greet`

   先理解 go-zero API 服务的基本结构：

   ```text
   .api -> handler -> logic -> response
   ```

2. `user`

   再理解 RPC 服务如何由 proto 生成：

   ```text
   .proto -> pb.go -> server -> logic -> response
   ```

3. `user-api`

   最后理解 API 服务如何调用 RPC 服务：

   ```text
   HTTP -> user-api -> userclient -> user-rpc
   ```

## 运行 user-api + user-rpc

终端 1：

```bash
cd examples/user
go run user.go
```

终端 2：

```bash
cd examples/user-api
go run user.go
```

测试：

```bash
curl http://127.0.0.1:8888/user/7
```
