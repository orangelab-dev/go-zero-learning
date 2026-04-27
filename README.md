# go-zero-learning

这个仓库用于系统学习 [go-zero](https://go-zero.dev/) 微服务框架。

当前目标不是一次性做完整项目，而是先通过小例子理解 go-zero 的核心开发方式：

- API 服务如何接收 HTTP 请求
- RPC 服务如何通过 proto 定义接口
- API 服务如何调用 RPC 服务
- handler、logic、svc、config 各自负责什么
- 后续逐步接入 MySQL、Redis、JWT、服务发现等能力

## 目录结构

```text
go-zero-learning/
├── docs/                  # 学习笔记和路线
├── examples/              # 独立示例，每个示例通常是独立 Go module
│   ├── greet/             # 最基础的 API 服务
│   ├── user/              # user RPC 服务
│   └── user-api/          # 调用 user RPC 的 HTTP API 服务
├── projects/              # 后续放完整微服务项目
└── scripts/               # 后续放辅助脚本
```

## 当前示例

### examples/greet

最基础的 go-zero API 示例。

它的作用是学习：

- `.api` 文件如何定义 HTTP 接口
- `handler -> logic -> svc` 的基本调用链
- `logic` 中如何返回业务响应

运行方式：

```bash
cd examples/greet
go run greet.go
```

### examples/user

一个 go-zero RPC 服务示例。

核心协议在：

```text
examples/user/user.proto
```

当前定义了一个 RPC 方法：

```proto
service User {
  rpc GetUser(GetUserRequest) returns(GetUserResponse);
}
```

调用 `GetUser` 时，服务会返回一个模拟用户：

```json
{
  "id": 传入的 id,
  "name": "alice"
}
```

运行方式：

```bash
cd examples/user
go run user.go
```

默认监听：

```text
0.0.0.0:8080
```

当前这个示例没有启用 etcd，适合本地学习直连调用。

### examples/user-api

一个 HTTP API 服务示例，用来演示 API 服务如何调用 RPC 服务。

当前接口：

```http
POST /user/login
GET  /user/:id
```

其中 `GET /user/:id` 会调用 `examples/user` 里的 `GetUser` RPC。

运行方式：

```bash
cd examples/user-api
go run user.go
```

默认监听：

```text
0.0.0.0:8888
```

测试请求：

```bash
curl http://127.0.0.1:8888/user/7
```

预期返回：

```json
{"id":7,"username":"alice","email":"alice@example.com"}
```

## API 调用 RPC 的完整链路

本仓库当前最重要的一条学习链路是：

```text
HTTP Client
  -> user-api: GET /user/:id
  -> internal/handler/getuserinfohandler.go
  -> internal/logic/getuserinfologic.go
  -> userclient.User.GetUser(...)
  -> user RPC: GetUser(...)
  -> internal/logic/getuserlogic.go
  -> 返回用户数据
```

对应关系：

| 层级 | 目录/文件 | 职责 |
| --- | --- | --- |
| API 协议 | `examples/user-api/user.api` | 定义 HTTP 路由、请求、响应 |
| API handler | `examples/user-api/internal/handler/` | 解析 HTTP 请求，调用 logic |
| API logic | `examples/user-api/internal/logic/` | 组织业务流程，调用 RPC |
| API svc | `examples/user-api/internal/svc/` | 初始化依赖，比如 RPC client |
| RPC 协议 | `examples/user/user.proto` | 定义 RPC 方法和消息结构 |
| RPC server | `examples/user/internal/server/` | 接收 gRPC 调用，转给 logic |
| RPC logic | `examples/user/internal/logic/` | 实现 RPC 业务逻辑 |
| RPC client | `examples/user/userclient/` | 给其他服务调用 RPC 使用 |

## go-zero 分层理解

### handler

贴近传输层。

在 API 服务里，handler 负责：

- 接收 HTTP 请求
- 解析 path/query/body 参数
- 调用 logic
- 输出 JSON 响应

原则：handler 保持轻薄，不写复杂业务逻辑。

### logic

业务逻辑层。

这里负责：

- 参数的业务校验
- 调用数据库、Redis、RPC
- 组织业务流程
- 返回业务结果

原则：大部分真正的业务代码写在 logic。

### svc

服务上下文，也可以理解成依赖容器。

这里放：

- 配置
- 数据库连接
- Redis client
- RPC client
- 其他公共依赖

原则：svc 负责准备依赖，不负责写业务流程。

### config

配置结构定义。

例如 `user-api` 中：

```go
type Config struct {
    rest.RestConf
    UserRpc zrpc.RpcClientConf
}
```

它对应 `etc/user-api.yaml` 中的：

```yaml
UserRpc:
  Endpoints:
  - 127.0.0.1:8080
```

## proto 生成代码后有什么

从 `user.proto` 生成后，会得到：

- `user/user.pb.go`：请求、响应结构体
- `user/user_grpc.pb.go`：gRPC client/server 接口
- `internal/server/userserver.go`：go-zero RPC server 适配层
- `internal/logic/getuserlogic.go`：业务逻辑入口
- `userclient/user.go`：给其他服务调用的 client 封装
- `user.go`：RPC 服务启动入口

平时最常改：

- `user.proto`
- `internal/logic/*.go`
- `internal/svc/*.go`
- `internal/config/*.go`
- `etc/*.yaml`

一般不手改：

- `*.pb.go`
- `*_grpc.pb.go`
- `internal/server/*.go`
- `userclient/*.go`

这些通常由 `goctl` 重新生成。

## 重新生成 RPC 代码

如果修改了 `examples/user/user.proto`，可以在 `examples/user` 下执行：

```bash
PATH=$PATH:/Users/orange/go/bin /Users/orange/go/bin/goctl rpc protoc user.proto --go_out=. --go-grpc_out=. --zrpc_out=.
go mod tidy
```

注意：

- `goctl` 负责生成 go-zero RPC 骨架
- `protoc-gen-go` 负责生成 protobuf Go 类型
- `protoc-gen-go-grpc` 负责生成 gRPC 接口
- 如果命令提示找不到插件，确认 `/Users/orange/go/bin` 是否在 `PATH` 中

## etcd 和直连模式

go-zero RPC 不强制要求 etcd。

本仓库当前使用直连模式：

```yaml
UserRpc:
  Endpoints:
  - 127.0.0.1:8080
```

这种方式适合本地学习，因为只要 RPC 服务运行在固定端口，API 服务就能直接连接。

如果使用 etcd，配置会类似：

```yaml
UserRpc:
  Etcd:
    Hosts:
      - 127.0.0.1:2379
    Key: user.rpc
```

etcd 的作用是服务注册与发现。适合多个服务实例、动态地址、负载均衡等更完整的微服务场景。

## 本地运行顺序

先启动 RPC：

```bash
cd examples/user
go run user.go
```

再启动 API：

```bash
cd examples/user-api
go run user.go
```

最后请求 API：

```bash
curl http://127.0.0.1:8888/user/7
```

如果遇到端口占用：

```bash
lsof -nP -iTCP:8080 -sTCP:LISTEN
lsof -nP -iTCP:8888 -sTCP:LISTEN
```

然后结束对应进程，或修改配置端口。

## 当前学习进度

已经完成：

- Git 仓库初始化和基础提交
- `.gitignore` 忽略 `.DS_Store`
- 基础 API 服务 `greet`
- RPC 服务 `user`
- HTTP API 服务 `user-api`
- `user-api -> user-rpc` 的直连调用

下一步建议：

1. 给 `user-rpc` 接入真实数据存储
2. 让 `user-api` 的登录接口调用 RPC
3. 增加 JWT 鉴权
4. 学习 etcd 服务注册发现
5. 开始设计完整项目，比如 mini-mall
