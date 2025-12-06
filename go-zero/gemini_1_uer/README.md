# Go-Zero 用户API服务

这是一个基于 go-zero 框架开发的用户API服务，实现了用户登录功能。

## 项目结构

```
gemini_1_uer/
├── user.api          # API定义文件
├── user.go           # 主入口文件
├── user.sql          # 数据库表结构
├── init_db.sql       # 数据库初始化脚本
├── etc/
│   └── user-api.yaml # 配置文件
└── internal/
    ├── config/       # 配置结构体
    ├── handler/      # HTTP处理器
    ├── logic/        # 业务逻辑
    ├── model/        # 数据模型
    ├── svc/          # 服务上下文
    └── types/        # 请求/响应类型
```

## 环境要求

- Go 1.18+
- MySQL 5.7+
- Redis 6.0+

## 快速开始

### 1. 初始化数据库

```bash
mysql -u root -p < init_db.sql
```

### 2. 启动Redis服务

确保Redis服务在 127.0.0.1:6379 运行

### 3. 修改配置（如果需要）

编辑 `etc/user-api.yaml` 文件，修改数据库和Redis连接信息：

```yaml
Name: user-api
Host: 0.0.0.0
Port: 8888

Mysql:
  DataSource: root:root@tcp(127.0.0.1:3306)/gozero?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai

Cache:
  - Host: 127.0.0.1:6379
    Pass: ""
```

### 4. 启动服务

```bash
go run user.go -f etc/user-api.yaml
```

或者编译后运行：

```bash
go build -o user-api user.go
./user-api -f etc/user-api.yaml
```

服务将在 http://localhost:8888 启动

## API接口

### 用户登录

**请求：**
```bash
curl -i -X POST http://localhost:8888/v1/user/login \
  -H "Content-Type: application/json" \
  -d '{"username": "root", "password": "root"}'
```

**成功响应：**
```json
{
  "id": 1,
  "name": "root",
  "token": "mock_token_root",
  "expire_at": "2025-12-31"
}
```

**错误响应：**
```json
用户不存在
```
或
```json
密码错误
```

## 测试账号

系统预置了以下测试账号：

| 用户名 | 密码 | 手机号 | 昵称 |
|--------|------|--------|------|
| root | root | 13800138000 | 管理员 |
| test | test123 | 13800138001 | 测试用户 |

## 故障排除

### 1. 数据库连接失败

检查MySQL服务是否启动，连接信息是否正确

### 2. Redis连接失败

检查Redis服务是否启动，连接信息是否正确

### 3. 404错误

确保使用正确的API路径：`/v1/user/login`（注意是单数user，不是复数users）

### 4. 编译错误

运行 `go mod tidy` 更新依赖

## 项目特点

1. **分层架构**：清晰的handler->logic->model分层
2. **代码生成**：通过API定义自动生成路由、处理器等代码
3. **内置缓存**：自动集成Redis缓存机制
4. **配置管理**：统一的配置文件和结构体映射
5. **错误处理**：统一的错误处理机制

## 学习要点

作为Go语言学习者，这个项目展示了：

1. Go模块结构和包管理
2. 接口设计和API定义
3. 数据库操作和ORM模式
4. 缓存策略和Redis集成
5. 错误处理和日志记录
6. 微服务架构设计