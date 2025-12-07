# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

TE（网络通信引擎）是一个用Go语言编写的多协议网络通信引擎，用于学习和实践网络通信协议的设计与实现。项目采用客户端-服务器架构，支持TCP、Stream、HTTP、WebSocket等多种协议。

**项目状态**：学习型项目，已完成TCP协议的基础实现，其他协议正在设计中。

## 常用开发命令

### 构建项目
```bash
# 构建服务端
go build -o te_server server/main/main.go

# 构建客户端
go build -o te_client client/main/main.go

# 构建API网关
go build -o api_gateway apiGateway/main.go

# 一次性构建所有组件
go build ./...
```

### 运行项目
```bash
# 运行服务端（监听 0.0.0.0:9501）
./te_server

# 运行客户端（连接到 127.0.0.1:9501）
./te_client

# 运行API网关
./api_gateway
```

### 测试
```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./server/Server
go test ./client/Client

# 运行测试并显示覆盖率
go test -cover ./...
```

### 代码格式化和检查
```bash
# 格式化代码
go fmt ./...

# 静态分析
go vet ./...

# 代码检查
golangci-lint run
```

### 清理
```bash
# 清理构建产物
rm -f te_server te_client api_gateway

# 清理go mod缓存
go clean -modcache

# 清理所有生成的文件
go clean ./...
```

## 核心架构

### 1. 目录结构
```
te/
├── apiGateway/          # API网关实现
├── client/             # 客户端实现
│   ├── main/           # 客户端入口
│   └── Client/         # 客户端核心逻辑
├── server/             # 服务端实现
│   ├── main/           # 服务端入口
│   └── Server/         # 服务端核心逻辑
├── README.md           # 项目文档
└── go.mod              # Go模块定义
```

### 2. 核心组件

#### Server (server/Server/Server.go)
- **主要职责**：管理客户端连接、处理网络事件
- **关键结构**：`Server`结构体包含连接池、事件回调函数
- **事件处理**：通过`CallEventFunc`统一处理error、start、connect、close、receive事件
- **连接管理**：使用map存储客户端连接，支持最大连接数限制

#### Client (client/Client/Client.go)
- **主要职责**：连接到服务器、处理服务器响应
- **关键结构**：`Client`结构体包含连接实例和事件回调
- **生命周期管理**：使用WaitGroup管理协程

#### TcpConnection (TcpConnection.go)
- **主要职责**：封装底层网络连接、处理数据收发
- **协议支持**：TCP、Stream、HTTP、WebSocket
- **缓冲区管理**：固定1024字节缓冲区

### 3. 事件驱动架构

项目采用事件驱动模式，支持以下事件：
- `OnError`：错误处理
- `OnStart`：服务器启动
- `OnConnect`：客户端连接
- `OnClose`：连接关闭
- `OnReceive`：消息接收

### 4. 协议设计

虽然当前主要实现了TCP协议，但架构设计支持：
- **TCP**：已完全实现，支持连接管理和粘包处理
- **Stream**：预留接口，用于流式数据传输
- **HTTP**：预留接口，可扩展为HTTP服务器
- **WebSocket**：预留接口，用于实时通信

## 开发注意事项

### 1. 并发安全
- 当前实现中服务端的Clients map没有使用锁保护
- 在多协程环境下访问共享数据需要添加sync.Mutex或sync.RWMutex

### 2. 已知问题
- `server/Server/Response.go`文件为空，需要实现响应处理逻辑
- 客户端连接地址配置有误（120.0.0.1应为127.0.0.1）
- 固定大小的缓冲区可能不适合所有场景

### 3. 扩展建议
- 实现HTTP和WebSocket协议支持
- 添加动态缓冲区调整机制
- 实现连接超时和心跳机制
- 添加更完善的错误处理和日志记录

## 代码风格

- 遵循Go语言标准代码规范
- 使用驼峰命名法
- 结构体字段使用大写字母开头（public）
- 方法名根据接收者类型使用大小写
- 添加必要的注释和文档

## 模块依赖

项目目前依赖较少，主要使用Go标准库：
- `net`：网络通信
- `fmt`：格式化输出
- `sync`：同步原语（WaitGroup等）
- `os`：操作系统接口

## 调试技巧

1. 使用事件回调函数添加日志输出
2. 在关键位置添加断点调试
3. 使用`go run`配合`-race`参数检测竞态条件
4. 使用网络调试工具（如telnet、nc）测试连接