# Gsystes Backend

中后台管理系统后端服务骨架，采用分层架构设计，基于 Go 语言开发。

## 技术栈

| 组件          | 技术                          |
| ------------- | ----------------------------- |
| 语言          | Go 1.26                       |
| Web 框架      | Gin                           |
| ORM           | GORM + MySQL 驱动             |
| 数据库        | MySQL 8.0+                    |
| 缓存          | Redis                         |
| 认证          | JWT (golang-jwt)              |
| 配置管理      | Viper                         |
| 日志          | Zap + Lumberjack（文件轮转）    |
| 密码加密      | bcrypt                        |
| 容器化        | Docker                        |

## 架构

采用分层架构，遵循**依赖倒置原则**：

```
┌──────────────────────────────────────┐
│          Communication               │  ← 通信层：HTTP 路由、Handler、中间件、DTO
│   (handler / middleware / router)    │
└──────────────┬───────────────────────┘
               │ 依赖
┌──────────────▼───────────────────────┐
│         Orchestration                │  ← 业务编排层：组合领域服务，编排业务流程
│   (service / request / response)     │
└──────────────┬───────────────────────┘
               │ 依赖
┌──────────────▼───────────────────────┐
│             Domain                    │  ← 领域层：核心业务逻辑与接口定义
│  (entity / repository interface /    │
│         domain service)              │
└──────────────┬───────────────────────┘
               │ 实现
┌──────────────▼───────────────────────┐
│              Data                     │  ← 数据层：仓储实现、ORM 模型、迁移
│   (repository impl / model / seed)   │
└──────────────┬───────────────────────┘
               │ 依赖
┌──────────────▼───────────────────────┐
│         Infrastructure                │  ← 基础设施层：配置、日志、数据库、缓存
│  (config / logger / database / auth) │
└──────────────────────────────────────┘
```

## 目录结构

```
├── cmd/server/main.go          # 应用入口，依赖注入与启动
├── config/
│   ├── config.yaml             # 生产配置
│   └── config.dev.yaml         # 开发配置
├── internal/
│   ├── communication/          # 通信层
│   │   ├── dto/                # 数据传输对象
│   │   ├── handler/            # HTTP 处理器
│   │   ├── middleware/         # Gin 中间件
│   │   └── router/            # 路由注册
│   ├── orchestration/          # 业务编排层
│   │   └── service/           # 编排服务
│   ├── domain/                 # 领域层
│   │   ├── entity/            # 领域实体
│   │   ├── repository/        # 仓储接口
│   │   └── service/           # 领域服务
│   ├── data/                   # 数据层
│   │   ├── model/             # GORM 数据模型
│   │   ├── repository/        # 仓储实现
│   │   ├── migration/         # 数据库迁移
│   │   └── seed/              # 种子数据
│   └── infrastructure/        # 基础设施层
│       ├── auth/              # JWT 认证
│       ├── cache/             # Redis 缓存
│       ├── config/            # 配置管理
│       ├── database/          # 数据库连接
│       ├── logger/            # 日志
│       ├── middleware/        # 跨层中间件工具
│       └── utils/             # 通用工具
├── bin/                       # 编译产物
├── logs/                      # 日志文件
├── Makefile                   # 构建命令
├── Dockerfile                 # Docker 镜像构建
└── README.md
```

## 快速开始

### 前置条件

- Go 1.26+
- MySQL 8.0+
- Redis（可选，无 Redis 时自动降级）

### 1. 创建数据库

```sql
CREATE DATABASE IF NOT EXISTS gsystes DEFAULT CHARSET utf8mb4;
```

### 2. 修改配置

编辑 `config/config.dev.yaml`，修改数据库连接信息：

```yaml
database:
  host: 127.0.0.1
  port: 3306
  username: root
  password: your_password
  database: gsystes
```

### 3. 启动服务

```bash
# 开发模式（热启动）
make dev

# 或手动运行
go run cmd/server/main.go -config config/config.dev.yaml
```

首次启动会自动完成：
1. 创建 `sys_users`、`sys_roles`、`sys_permissions` 三张表
2. 初始化种子数据（管理员账户、角色、权限）

### 4. 验证

```bash
# 登录获取 Token
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

## 默认账户

| 用户名  | 密码      | 角色         |
| ------- | --------- | ------------ |
| admin   | admin123  | 超级管理员   |

> 首次登录后请及时修改密码。

## API 文档

### 认证

| 方法   | 路径                  | 说明     |
| ------ | --------------------- | -------- |
| `POST` | `/api/v1/auth/login`  | 用户登录 |

**登录请求：**

```json
{
  "username": "admin",
  "password": "admin123"
}
```

**登录响应：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": 1,
      "username": "admin",
      "nickname": "超级管理员",
      "avatar": "",
      "role_id": 1
    }
  }
}
```

### 用户管理（需认证）

在请求头中添加 `Authorization: Bearer <token>`。

| 方法     | 路径                    | 说明     |
| -------- | ----------------------- | -------- |
| `POST`   | `/api/v1/users`         | 创建用户 |
| `DELETE` | `/api/v1/users/:id`     | 删除用户 |
| `PUT`    | `/api/v1/users/:id`     | 编辑用户 |
| `GET`    | `/api/v1/users/:id`     | 用户详情 |
| `GET`    | `/api/v1/users`         | 用户列表 |
| `PUT`    | `/api/v1/users/password`| 修改密码 |

### 通用响应格式

```json
{
  "code": 0,        // 0 成功，非 0 失败
  "message": "success",
  "data": {}        // 响应数据
}
```

分页响应：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [],
    "total": 0,
    "page": 1,
    "page_size": 10
  }
}
```

## 配置说明

| 配置项            | 默认值 | 说明 |
| ----------------- | ------ | ---- |
| `server.port`     | 8080   | 服务端口 |
| `server.mode`     | release | Gin 运行模式 |
| `database.driver` | mysql  | 数据库驱动 |
| `jwt.secret`      | -      | JWT 签名密钥 |
| `jwt.expire_hours`| 24     | Token 过期时间 |
| `log.level`       | info   | 日志级别 |

## 测试

```bash
# 运行所有测试
make test

# 查看测试覆盖率
go test ./... -v -cover
```

## 构建与部署

### 本地构建

```bash
# 编译
make build

# 运行编译产物
./bin/gsystes-server -config config/config.yaml
```

### Docker 部署

```bash
# 构建镜像
docker build -t gsystes-server .

# 运行容器
docker run -d \
  --name gsystes-server \
  -p 8080:8080 \
  -v $(pwd)/config/config.yaml:/app/config/config.yaml \
  gsystes-server
```

### Makefile 命令

| 命令          | 说明 |
| ------------- | ---- |
| `make build`  | 编译项目 |
| `make run`    | 编译并运行（使用开发配置） |
| `make dev`    | 直接开发运行 |
| `make test`   | 运行测试 |
| `make clean`  | 清理构建产物和日志 |
| `make lint`   | 代码检查 |

## 许可

MIT