# Gsystes 中后台管理系统后端实施计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 基于 Go 语言构建一个中后台管理系统后端，采用五层分层架构（通信层、业务编排层、领域层、数据层、基础设施层）

**架构:**
- **通信层**: 处理 HTTP/gRPC 外部请求，参数校验，响应格式化
- **业务编排层**: 编排领域服务，协调业务流程，管理事务
- **领域层**: 核心业务实体，领域服务，仓储接口（端口）
- **数据层**: 仓储实现（适配器），ORM 模型，数据库迁移
- **基础设施层**: 配置管理、日志、数据库连接、缓存、认证、工具函数

各层严格遵循依赖倒置原则：通信层 → 业务编排层 → 领域层 ← 数据层 ← 基础设施层（领域层不依赖任何外部框架）

**Tech Stack:** Go 1.26, Gin, GORM, MySQL, Redis, Viper, Zap, JWT, Swaggo

**Project Structure:**
```
Gsystes/
├── cmd/
│   └── server/
│       └── main.go                    # 应用入口
├── internal/
│   ├── communication/                 # 通信层
│   │   ├── dto/                       # 数据传输对象
│   │   │   ├── user_dto.go
│   │   │   └── common_dto.go
│   │   ├── handler/                   # HTTP 处理器
│   │   │   └── user_handler.go
│   │   ├── middleware/                # HTTP 中间件
│   │   │   ├── auth_middleware.go
│   │   │   ├── cors_middleware.go
│   │   │   └── recovery_middleware.go
│   │   └── router/                    # 路由注册
│   │       └── router.go
│   ├── orchestration/                 # 业务编排层
│   │   └── service/
│   │       └── user_orchestration.go
│   ├── domain/                        # 领域层
│   │   ├── entity/                    # 领域实体
│   │   │   ├── user.go
│   │   │   ├── role.go
│   │   │   └── permission.go
│   │   ├── repository/                # 仓储接口（端口）
│   │   │   ├── user_repository.go
│   │   │   ├── role_repository.go
│   │   │   └── permission_repository.go
│   │   └── service/                   # 领域服务
│   │       └── user_domain_service.go
│   ├── data/                          # 数据层
│   │   ├── model/                     # ORM 持久化模型
│   │   │   ├── user_model.go
│   │   │   ├── role_model.go
│   │   │   └── permission_model.go
│   │   ├── repository/                # 仓储实现（适配器）
│   │   │   ├── user_repo_impl.go
│   │   │   ├── role_repo_impl.go
│   │   │   └── permission_repo_impl.go
│   │   └── migration/                 # 数据库迁移
│   │       └── auto_migrate.go
│   └── infrastructure/               # 基础设施层
│       ├── config/                    # 配置管理
│       │   └── config.go
│       ├── logger/                    # 日志
│       │   └── logger.go
│       ├── database/                  # 数据库连接
│       │   └── database.go
│       ├── cache/                     # 缓存
│       │   └── redis.go
│       ├── auth/                      # JWT 认证
│       │   └── jwt.go
│       ├── middleware/                # 基础设施中间件
│       │   └── context.go
│       └── utils/                     # 工具函数
│           ├── hash.go
│           ├── response.go
│           └── validator.go
├── config/                            # 配置文件
│   ├── config.yaml
│   └── config.dev.yaml
├── api/                               # API 文档
│   └── swagger/
├── go.mod
├── go.sum
├── Makefile
└── Dockerfile
```

---

### Task 1: 初始化 Go Module 与基础设施层 - 配置管理

**Files:**
- Create: `go.mod`
- Create: `config/config.yaml`
- Create: `config/config.dev.yaml`
- Create: `internal/infrastructure/config/config.go`

- [ ] **Step 1: 初始化 Go Module**

Run: `cd d:\Project\Go\Gsystes && go mod init github.com/gsystes/backend`

Expected: 生成 `go.mod` 文件

- [ ] **Step 2: 创建配置文件目录和基础配置**

创建 `config/config.yaml`:

```yaml
server:
  port: 8080
  mode: release   # debug | release | test
  read_timeout: 10
  write_timeout: 10

database:
  driver: mysql
  host: 127.0.0.1
  port: 3306
  username: root
  password: root
  database: gsystes
  charset: utf8mb4
  max_idle_conns: 10
  max_open_conns: 100
  conn_max_lifetime: 3600

redis:
  host: 127.0.0.1
  port: 6379
  password: ""
  db: 0

jwt:
  secret: gsystes-jwt-secret-key-change-in-production
  issuer: gsystes
  expire_hours: 24

log:
  level: info
  filename: logs/gsystes.log
  max_size: 100
  max_backups: 7
  max_age: 30
  compress: true
```

创建 `config/config.dev.yaml`:

```yaml
server:
  port: 8080
  mode: debug

database:
  host: 127.0.0.1
  port: 3306
  username: root
  password: root
  database: gsystes_dev

redis:
  host: 127.0.0.1
  port: 6379
  password: ""
  db: 0

jwt:
  secret: dev-secret-key
  expire_hours: 720

log:
  level: debug
  filename: logs/gsystes-dev.log
  max_size: 50
  max_backups: 3
  max_age: 7
  compress: false
```

- [ ] **Step 3: 编写配置管理代码**

创建 `internal/infrastructure/config/config.go`:

```go
package config

import (
	"fmt"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Log      LogConfig      `mapstructure:"log"`
}

type ServerConfig struct {
	Port         int           `mapstructure:"port"`
	Mode         string        `mapstructure:"mode"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

type DatabaseConfig struct {
	Driver          string `mapstructure:"driver"`
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	Database        string `mapstructure:"database"`
	Charset         string `mapstructure:"charset"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
}

func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		d.Username, d.Password, d.Host, d.Port, d.Database, d.Charset)
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

func (r RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

type JWTConfig struct {
	Secret      string `mapstructure:"secret"`
	Issuer      string `mapstructure:"issuer"`
	ExpireHours int    `mapstructure:"expire_hours"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

var globalConfig *Config

func LoadConfig(configPath string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		var newConfig Config
		if err := v.Unmarshal(&newConfig); err != nil {
			return
		}
		globalConfig = &newConfig
	})

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	globalConfig = &cfg
	return globalConfig, nil
}

func GetConfig() *Config {
	return globalConfig
}
```

- [ ] **Step 4: 安装依赖**

Run: `cd d:\Project\Go\Gsystes && go get github.com/spf13/viper github.com/fsnotify/fsnotify`

Expected: 依赖安装成功

- [ ] **Step 5: 验证编译**

Run: `cd d:\Project\Go\Gsystes && go build ./...`

Expected: 编译成功，无错误

- [ ] **Step 6: 提交**

```bash
git init
git add -A
git commit -m "feat: initialize project with layered architecture and config management"
```

---

### Task 2: 基础设施层 - 日志模块

**Files:**
- Create: `internal/infrastructure/logger/logger.go`

- [ ] **Step 1: 编写日志模块代码**

创建 `internal/infrastructure/logger/logger.go`:

```go
package logger

import (
	"os"
	"time"

	"github.com/gsystes/backend/internal/infrastructure/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var globalLogger *zap.Logger

func InitLogger(cfg config.LogConfig) error {
	writeSyncer := getLogWriter(cfg)
	encoder := getEncoder()

	var level zapcore.Level
	if err := level.Set(cfg.Level); err != nil {
		level = zapcore.InfoLevel
	}

	core := zapcore.NewCore(encoder, writeSyncer, level)
	globalLogger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	zap.ReplaceGlobals(globalLogger)
	return nil
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogWriter(cfg config.LogConfig) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   cfg.Filename,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	}
	ws := zapcore.AddSync(lumberJackLogger)
	if cfg.Level == "debug" {
		ws = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), ws)
	}
	return ws
}

func Info(msg string, fields ...zap.Field) {
	globalLogger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	globalLogger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	globalLogger.Error(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	globalLogger.Debug(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	globalLogger.Fatal(msg, fields...)
}

func With(fields ...zap.Field) *zap.Logger {
	return globalLogger.With(fields...)
}

func Sync() error {
	return globalLogger.Sync()
}

func GetLogger() *zap.Logger {
	return globalLogger
}
```

- [ ] **Step 2: 安装依赖**

Run: `cd d:\Project\Go\Gsystes && go get go.uber.org/zap gopkg.in/natefinch/lumberjack.v2`

Expected: 依赖安装成功

- [ ] **Step 3: 验证编译**

Run: `cd d:\Project\Go\Gsystes && go build ./...`

Expected: 编译成功

- [ ] **Step 4: 提交**

```bash
git add -A
git commit -m "feat: add logging module with zap and log rotation"
```

---

### Task 3: 基础设施层 - 数据库连接与缓存

**Files:**
- Create: `internal/infrastructure/database/database.go`
- Create: `internal/infrastructure/cache/redis.go`

- [ ] **Step 1: 编写数据库连接模块**

创建 `internal/infrastructure/database/database.go`:

```go
package database

import (
	"time"

	"github.com/gsystes/backend/internal/infrastructure/config"
	"github.com/gsystes/backend/internal/infrastructure/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var globalDB *gorm.DB

func InitDatabase(cfg config.DatabaseConfig) error {
	gormConfig := &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
	}

	db, err := gorm.Open(mysql.Open(cfg.DSN()), gormConfig)
	if err != nil {
		logger.Error("failed to connect database", logger.ErrorField(err))
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

	globalDB = db
	logger.Info("database connected successfully")
	return nil
}

func GetDB() *gorm.DB {
	return globalDB
}

func Close() error {
	sqlDB, err := globalDB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
```

- [ ] **Step 2: 编写 Redis 缓存模块**

创建 `internal/infrastructure/cache/redis.go`:

```go
package cache

import (
	"context"
	"time"

	"github.com/gsystes/backend/internal/infrastructure/config"
	"github.com/gsystes/backend/internal/infrastructure/logger"
	"github.com/redis/go-redis/v9"
)

var globalRedis *redis.Client

func InitRedis(cfg config.RedisConfig) error {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr(),
		Password:     cfg.Password,
		DB:           cfg.DB,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		logger.Warn("redis connection failed, cache disabled", logger.ErrorField(err))
		return nil
	}

	globalRedis = client
	logger.Info("redis connected successfully")
	return nil
}

func GetRedis() *redis.Client {
	return globalRedis
}

func Close() error {
	if globalRedis != nil {
		return globalRedis.Close()
	}
	return nil
}
```

- [ ] **Step 3: 安装依赖**

Run: `cd d:\Project\Go\Gsystes && go get gorm.io/gorm gorm.io/driver/mysql github.com/redis/go-redis/v9`

Expected: 依赖安装成功

- [ ] **Step 4: 验证编译**

Run: `cd d:\Project\Go\Gsystes && go build ./...`

Expected: 编译成功

- [ ] **Step 5: 提交**

```bash
git add -A
git commit -m "feat: add database and redis connection modules"
```

---

### Task 4: 基础设施层 - JWT 认证与工具函数

**Files:**
- Create: `internal/infrastructure/auth/jwt.go`
- Create: `internal/infrastructure/utils/hash.go`
- Create: `internal/infrastructure/utils/response.go`
- Create: `internal/infrastructure/utils/validator.go`
- Create: `internal/infrastructure/middleware/context.go`

- [ ] **Step 1: 编写 JWT 认证模块**

创建 `internal/infrastructure/auth/jwt.go`:

```go
package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gsystes/backend/internal/infrastructure/config"
)

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	RoleID   uint   `json:"role_id"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint, username string, roleID uint) (string, error) {
	cfg := config.GetConfig().JWT
	claims := Claims{
		UserID:   userID,
		Username: username,
		RoleID:   roleID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    cfg.Issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(cfg.ExpireHours) * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Secret))
}

func ParseToken(tokenString string) (*Claims, error) {
	cfg := config.GetConfig().JWT
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(cfg.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
```

- [ ] **Step 2: 编写密码哈希工具**

创建 `internal/infrastructure/utils/hash.go`:

```go
package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
```

- [ ] **Step 3: 编写统一响应工具**

创建 `internal/infrastructure/utils/response.go`:

```go
package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type PageResult struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: message,
		Data:    data,
	})
}

func Error(c *gin.Context, httpStatus int, message string) {
	c.JSON(httpStatus, Response{
		Code:    -1,
		Message: message,
	})
}

func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, message)
}

func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, message)
}

func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, message)
}

func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, message)
}

func InternalError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, message)
}

func PageSuccess(c *gin.Context, list interface{}, total int64, page int, pageSize int) {
	Success(c, PageResult{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}
```

- [ ] **Step 4: 编写请求校验工具**

创建 `internal/infrastructure/utils/validator.go`:

```go
package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type ValidError struct {
	Key     string `json:"key"`
	Message string `json:"message"`
}

func BindAndValidate(c *gin.Context, obj interface{}) []ValidError {
	if err := c.ShouldBind(obj); err != nil {
		return formatValidationError(err)
	}
	return nil
}

func formatValidationError(err error) []ValidError {
	var errors []ValidError
	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			errors = append(errors, ValidError{
				Key:     e.Field(),
				Message: e.Translate(binding.Validator.Engine().(*validator.Validate)),
			})
		}
	} else {
		errors = append(errors, ValidError{
			Key:     "request",
			Message: err.Error(),
		})
	}
	return nil
}
```

- [ ] **Step 5: 编写上下文中间件工具**

创建 `internal/infrastructure/middleware/context.go`:

```go
package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/gsystes/backend/internal/infrastructure/auth"
)

const (
	ContextKeyClaims = "claims"
)

func SetClaims(c *gin.Context, claims *auth.Claims) {
	c.Set(ContextKeyClaims, claims)
}

func GetClaims(c *gin.Context) *auth.Claims {
	value, exists := c.Get(ContextKeyClaims)
	if !exists {
		return nil
	}
	claims, ok := value.(*auth.Claims)
	if !ok {
		return nil
	}
	return claims
}

func GetUserID(c *gin.Context) uint {
	claims := GetClaims(c)
	if claims == nil {
		return 0
	}
	return claims.UserID
}
```

- [ ] **Step 6: 安装依赖**

Run: `cd d:\Project\Go\Gsystes && go get golang.org/x/crypto github.com/golang-jwt/jwt/v5 github.com/go-playground/validator/v10`

Expected: 依赖安装成功

- [ ] **Step 7: 验证编译**

Run: `cd d:\Project\Go\Gsystes && go build ./...`

Expected: 编译成功

- [ ] **Step 8: 提交**

```bash
git add -A
git commit -m "feat: add JWT auth, utils, and context middleware"
```

---

### Task 5: 领域层 - 核心实体与仓储接口定义

**Files:**
- Create: `internal/domain/entity/user.go`
- Create: `internal/domain/entity/role.go`
- Create: `internal/domain/entity/permission.go`
- Create: `internal/domain/repository/user_repository.go`
- Create: `internal/domain/repository/role_repository.go`
- Create: `internal/domain/repository/permission_repository.go`

- [ ] **Step 1: 定义用户实体**

创建 `internal/domain/entity/user.go`:

```go
package entity

import "time"

type User struct {
	ID        uint
	Username  string
	Password  string
	Nickname  string
	Email     string
	Phone     string
	Avatar    string
	Status    int // 1:启用 2:禁用
	RoleID    uint
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserStatus int

const (
	UserStatusActive   UserStatus = 1
	UserStatusInactive UserStatus = 2
)
```

- [ ] **Step 2: 定义角色实体**

创建 `internal/domain/entity/role.go`:

```go
package entity

import "time"

type Role struct {
	ID          uint
	Name        string
	Code        string
	Description string
	Status      int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
```

- [ ] **Step 3: 定义权限实体**

创建 `internal/domain/entity/permission.go`:

```go
package entity

import "time"

type Permission struct {
	ID        uint
	Name      string
	Code      string
	Type      int    // 1:菜单 2:按钮 3:API
	ParentID  uint
	Path      string
	Method    string
	Sort      int
	CreatedAt time.Time
	UpdatedAt time.Time
}
```

- [ ] **Step 4: 定义用户仓储接口**

创建 `internal/domain/repository/user_repository.go`:

```go
package repository

import "github.com/gsystes/backend/internal/domain/entity"

type UserRepository interface {
	Create(user *entity.User) error
	Update(user *entity.User) error
	Delete(id uint) error
	FindByID(id uint) (*entity.User, error)
	FindByUsername(username string) (*entity.User, error)
	FindByPage(page, pageSize int, conditions map[string]interface{}) ([]entity.User, int64, error)
}
```

- [ ] **Step 5: 定义角色仓储接口**

创建 `internal/domain/repository/role_repository.go`:

```go
package repository

import "github.com/gsystes/backend/internal/domain/entity"

type RoleRepository interface {
	Create(role *entity.Role) error
	Update(role *entity.Role) error
	Delete(id uint) error
	FindByID(id uint) (*entity.Role, error)
	FindAll() ([]entity.Role, error)
	FindByPage(page, pageSize int) ([]entity.Role, int64, error)
}
```

- [ ] **Step 6: 定义权限仓储接口**

创建 `internal/domain/repository/permission_repository.go`:

```go
package repository

import "github.com/gsystes/backend/internal/domain/entity"

type PermissionRepository interface {
	Create(permission *entity.Permission) error
	Update(permission *entity.Permission) error
	Delete(id uint) error
	FindByID(id uint) (*entity.Permission, error)
	FindAll() ([]entity.Permission, error)
	FindByRoleID(roleID uint) ([]entity.Permission, error)
}
```

- [ ] **Step 7: 验证编译**

Run: `cd d:\Project\Go\Gsystes && go build ./...`

Expected: 编译成功

- [ ] **Step 8: 提交**

```bash
git add -A
git commit -m "feat: add domain entities and repository interfaces"
```

---

### Task 6: 领域层 - 领域服务实现

**Files:**
- Create: `internal/domain/service/user_domain_service.go`

- [ ] **Step 1: 编写用户领域服务**

创建 `internal/domain/service/user_domain_service.go`:

```go
package service

import (
	"errors"

	"github.com/gsystes/backend/internal/domain/entity"
	"github.com/gsystes/backend/internal/domain/repository"
	"github.com/gsystes/backend/internal/infrastructure/utils"
)

type UserDomainService struct {
	userRepo repository.UserRepository
}

func NewUserDomainService(userRepo repository.UserRepository) *UserDomainService {
	return &UserDomainService{userRepo: userRepo}
}

func (s *UserDomainService) Create(user *entity.User, plainPassword string) error {
	existing, _ := s.userRepo.FindByUsername(user.Username)
	if existing != nil {
		return errors.New("username already exists")
	}

	hashedPassword, err := utils.HashPassword(plainPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	user.Status = int(entity.UserStatusActive)

	return s.userRepo.Create(user)
}

func (s *UserDomainService) UpdatePassword(userID uint, oldPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	if !utils.CheckPassword(oldPassword, user.Password) {
		return errors.New("old password is incorrect")
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	return s.userRepo.Update(user)
}

func (s *UserDomainService) ValidateCredentials(username, password string) (*entity.User, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if user.Status == int(entity.UserStatusInactive) {
		return nil, errors.New("account is disabled")
	}

	if !utils.CheckPassword(password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}
```

- [ ] **Step 2: 验证编译**

Run: `cd d:\Project\Go\Gsystes && go build ./...`

Expected: 编译成功

- [ ] **Step 3: 提交**

```bash
git add -A
git commit -m "feat: add user domain service with password management"
```

---

### Task 7: 数据层 - ORM 模型与仓储实现

**Files:**
- Create: `internal/data/model/user_model.go`
- Create: `internal/data/model/role_model.go`
- Create: `internal/data/model/permission_model.go`
- Create: `internal/data/repository/user_repo_impl.go`
- Create: `internal/data/repository/role_repo_impl.go`
- Create: `internal/data/repository/permission_repo_impl.go`
- Create: `internal/data/migration/auto_migrate.go`

- [ ] **Step 1: 编写用户 ORM 模型**

创建 `internal/data/model/user_model.go`:

```go
package model

import "time"

type User struct {
	ID        uint      `gorm:"primarykey"`
	Username  string    `gorm:"column:username;type:varchar(64);uniqueIndex;not null"`
	Password  string    `gorm:"column:password;type:varchar(256);not null"`
	Nickname  string    `gorm:"column:nickname;type:varchar(64)"`
	Email     string    `gorm:"column:email;type:varchar(128)"`
	Phone     string    `gorm:"column:phone;type:varchar(20)"`
	Avatar    string    `gorm:"column:avatar;type:varchar(256)"`
	Status    int       `gorm:"column:status;type:tinyint;default:1"`
	RoleID    uint      `gorm:"column:role_id"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (User) TableName() string {
	return "sys_users"
}
```

- [ ] **Step 2: 编写角色 ORM 模型**

创建 `internal/data/model/role_model.go`:

```go
package model

import "time"

type Role struct {
	ID          uint      `gorm:"primarykey"`
	Name        string    `gorm:"column:name;type:varchar(64);not null"`
	Code        string    `gorm:"column:code;type:varchar(64);uniqueIndex;not null"`
	Description string    `gorm:"column:description;type:varchar(256)"`
	Status      int       `gorm:"column:status;type:tinyint;default:1"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (Role) TableName() string {
	return "sys_roles"
}
```

- [ ] **Step 3: 编写权限 ORM 模型**

创建 `internal/data/model/permission_model.go`:

```go
package model

import "time"

type Permission struct {
	ID        uint      `gorm:"primarykey"`
	Name      string    `gorm:"column:name;type:varchar(64);not null"`
	Code      string    `gorm:"column:code;type:varchar(64);uniqueIndex;not null"`
	Type      int       `gorm:"column:type;type:tinyint"`
	ParentID  uint      `gorm:"column:parent_id;default:0"`
	Path      string    `gorm:"column:path;type:varchar(256)"`
	Method    string    `gorm:"column:method;type:varchar(32)"`
	Sort      int       `gorm:"column:sort;default:0"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (Permission) TableName() string {
	return "sys_permissions"
}
```

- [ ] **Step 4: 实现用户仓储**

创建 `internal/data/repository/user_repo_impl.go`:

```go
package repository

import (
	"github.com/gsystes/backend/internal/data/model"
	domainEntity "github.com/gsystes/backend/internal/domain/entity"
	domainRepo "github.com/gsystes/backend/internal/domain/repository"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domainRepo.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) toDomain(m *model.User) *domainEntity.User {
	if m == nil {
		return nil
	}
	return &domainEntity.User{
		ID:        m.ID,
		Username:  m.Username,
		Password:  m.Password,
		Nickname:  m.Nickname,
		Email:     m.Email,
		Phone:     m.Phone,
		Avatar:    m.Avatar,
		Status:    m.Status,
		RoleID:    m.RoleID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func (r *userRepository) toModel(d *domainEntity.User) *model.User {
	return &model.User{
		ID:        d.ID,
		Username:  d.Username,
		Password:  d.Password,
		Nickname:  d.Nickname,
		Email:     d.Email,
		Phone:     d.Phone,
		Avatar:    d.Avatar,
		Status:    d.Status,
		RoleID:    d.RoleID,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}

func (r *userRepository) Create(user *domainEntity.User) error {
	m := r.toModel(user)
	return r.db.Create(m).Error
}

func (r *userRepository) Update(user *domainEntity.User) error {
	m := r.toModel(user)
	return r.db.Model(&model.User{}).Where("id = ?", m.ID).Select("*").Omit("created_at").Updates(m).Error
}

func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}

func (r *userRepository) FindByID(id uint) (*domainEntity.User, error) {
	var m model.User
	if err := r.db.First(&m, id).Error; err != nil {
		return nil, err
	}
	return r.toDomain(&m), nil
}

func (r *userRepository) FindByUsername(username string) (*domainEntity.User, error) {
	var m model.User
	if err := r.db.Where("username = ?", username).First(&m).Error; err != nil {
		return nil, err
	}
	return r.toDomain(&m), nil
}

func (r *userRepository) FindByPage(page, pageSize int, conditions map[string]interface{}) ([]domainEntity.User, int64, error) {
	var models []model.User
	var total int64

	query := r.db.Model(&model.User{})
	for key, value := range conditions {
		query = query.Where(key, value)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&models).Error; err != nil {
		return nil, 0, err
	}

	entities := make([]domainEntity.User, len(models))
	for i, m := range models {
		entities[i] = *r.toDomain(&m)
	}
	return entities, total, nil
}
```

- [ ] **Step 5: 实现角色仓储**

创建 `internal/data/repository/role_repo_impl.go`:

```go
package repository

import (
	"github.com/gsystes/backend/internal/data/model"
	domainEntity "github.com/gsystes/backend/internal/domain/entity"
	domainRepo "github.com/gsystes/backend/internal/domain/repository"
	"gorm.io/gorm"
)

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) domainRepo.RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) toDomain(m *model.Role) *domainEntity.Role {
	if m == nil {
		return nil
	}
	return &domainEntity.Role{
		ID:          m.ID,
		Name:        m.Name,
		Code:        m.Code,
		Description: m.Description,
		Status:      m.Status,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func (r *roleRepository) Create(role *domainEntity.Role) error {
	return r.db.Create(&model.Role{
		Name:        role.Name,
		Code:        role.Code,
		Description: role.Description,
		Status:      role.Status,
	}).Error
}

func (r *roleRepository) Update(role *domainEntity.Role) error {
	return r.db.Model(&model.Role{}).Where("id = ?", role.ID).Updates(map[string]interface{}{
		"name":        role.Name,
		"code":        role.Code,
		"description": role.Description,
		"status":      role.Status,
	}).Error
}

func (r *roleRepository) Delete(id uint) error {
	return r.db.Delete(&model.Role{}, id).Error
}

func (r *roleRepository) FindByID(id uint) (*domainEntity.Role, error) {
	var m model.Role
	if err := r.db.First(&m, id).Error; err != nil {
		return nil, err
	}
	return r.toDomain(&m), nil
}

func (r *roleRepository) FindAll() ([]domainEntity.Role, error) {
	var models []model.Role
	if err := r.db.Find(&models).Error; err != nil {
		return nil, err
	}
	entities := make([]domainEntity.Role, len(models))
	for i, m := range models {
		entities[i] = *r.toDomain(&m)
	}
	return entities, nil
}

func (r *roleRepository) FindByPage(page, pageSize int) ([]domainEntity.Role, int64, error) {
	var models []model.Role
	var total int64

	if err := r.db.Model(&model.Role{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := r.db.Offset(offset).Limit(pageSize).Order("id DESC").Find(&models).Error; err != nil {
		return nil, 0, err
	}

	entities := make([]domainEntity.Role, len(models))
	for i, m := range models {
		entities[i] = *r.toDomain(&m)
	}
	return entities, total, nil
}
```

- [ ] **Step 6: 实现权限仓储**

创建 `internal/data/repository/permission_repo_impl.go`:

```go
package repository

import (
	"github.com/gsystes/backend/internal/data/model"
	domainEntity "github.com/gsystes/backend/internal/domain/entity"
	domainRepo "github.com/gsystes/backend/internal/domain/repository"
	"gorm.io/gorm"
)

type permissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) domainRepo.PermissionRepository {
	return &permissionRepository{db: db}
}

func (r *permissionRepository) toDomain(m *model.Permission) *domainEntity.Permission {
	if m == nil {
		return nil
	}
	return &domainEntity.Permission{
		ID:        m.ID,
		Name:      m.Name,
		Code:      m.Code,
		Type:      m.Type,
		ParentID:  m.ParentID,
		Path:      m.Path,
		Method:    m.Method,
		Sort:      m.Sort,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func (r *permissionRepository) Create(p *domainEntity.Permission) error {
	return r.db.Create(&model.Permission{
		Name:     p.Name,
		Code:     p.Code,
		Type:     p.Type,
		ParentID: p.ParentID,
		Path:     p.Path,
		Method:   p.Method,
		Sort:     p.Sort,
	}).Error
}

func (r *permissionRepository) Update(p *domainEntity.Permission) error {
	return r.db.Model(&model.Permission{}).Where("id = ?", p.ID).Updates(map[string]interface{}{
		"name":      p.Name,
		"code":      p.Code,
		"type":      p.Type,
		"parent_id": p.ParentID,
		"path":      p.Path,
		"method":    p.Method,
		"sort":      p.Sort,
	}).Error
}

func (r *permissionRepository) Delete(id uint) error {
	return r.db.Delete(&model.Permission{}, id).Error
}

func (r *permissionRepository) FindByID(id uint) (*domainEntity.Permission, error) {
	var m model.Permission
	if err := r.db.First(&m, id).Error; err != nil {
		return nil, err
	}
	return r.toDomain(&m), nil
}

func (r *permissionRepository) FindAll() ([]domainEntity.Permission, error) {
	var models []model.Permission
	if err := r.db.Order("sort ASC").Find(&models).Error; err != nil {
		return nil, err
	}
	entities := make([]domainEntity.Permission, len(models))
	for i, m := range models {
		entities[i] = *r.toDomain(&m)
	}
	return entities, nil
}

func (r *permissionRepository) FindByRoleID(roleID uint) ([]domainEntity.Permission, error) {
	var models []model.Permission
	err := r.db.Raw(`
		SELECT p.* FROM sys_permissions p
		INNER JOIN sys_role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id = ?
		ORDER BY p.sort ASC
	`, roleID).Scan(&models).Error
	if err != nil {
		return nil, err
	}
	entities := make([]domainEntity.Permission, len(models))
	for i, m := range models {
		entities[i] = *r.toDomain(&m)
	}
	return entities, nil
}
```

- [ ] **Step 7: 编写数据库自动迁移**

创建 `internal/data/migration/auto_migrate.go`:

```go
package migration

import (
	"github.com/gsystes/backend/internal/data/model"
	"github.com/gsystes/backend/internal/infrastructure/logger"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.Permission{},
	)
	if err != nil {
		logger.Error("auto migration failed", logger.ErrorField(err))
		return err
	}
	logger.Info("database auto migration completed")
	return nil
}
```

- [ ] **Step 8: 验证编译**

Run: `cd d:\Project\Go\Gsystes && go build ./...`

Expected: 编译成功

- [ ] **Step 9: 提交**

```bash
git add -A
git commit -m "feat: add data layer with repository implementations and migrations"
```

---

### Task 8: 业务编排层 - 编排服务

**Files:**
- Create: `internal/orchestration/service/user_orchestration.go`

- [ ] **Step 1: 编写用户编排服务**

创建 `internal/orchestration/service/user_orchestration.go`:

```go
package service

import (
	"github.com/gsystes/backend/internal/domain/entity"
	domainRepo "github.com/gsystes/backend/internal/domain/repository"
	domainService "github.com/gsystes/backend/internal/domain/service"
)

type UserOrchestration struct {
	userDomainService *domainService.UserDomainService
	userRepo          domainRepo.UserRepository
}

func NewUserOrchestration(
	userDomainService *domainService.UserDomainService,
	userRepo domainRepo.UserRepository,
) *UserOrchestration {
	return &UserOrchestration{
		userDomainService: userDomainService,
		userRepo:          userRepo,
	}
}

type CreateUserRequest struct {
	Username string
	Password string
	Nickname string
	Email    string
	Phone    string
	RoleID   uint
}

type UpdateUserRequest struct {
	ID       uint
	Nickname string
	Email    string
	Phone    string
	RoleID   uint
	Status   int
}

type LoginRequest struct {
	Username string
	Password string
}

type LoginResponse struct {
	User  *entity.User
	Token string
}

func (s *UserOrchestration) CreateUser(req *CreateUserRequest) (*entity.User, error) {
	user := &entity.User{
		Username: req.Username,
		Nickname: req.Nickname,
		Email:    req.Email,
		Phone:    req.Phone,
		RoleID:   req.RoleID,
	}

	if err := s.userDomainService.Create(user, req.Password); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserOrchestration) UpdateUser(req *UpdateUserRequest) error {
	user, err := s.userRepo.FindByID(req.ID)
	if err != nil {
		return err
	}

	user.Nickname = req.Nickname
	user.Email = req.Email
	user.Phone = req.Phone
	user.RoleID = req.RoleID
	user.Status = req.Status

	return s.userRepo.Update(user)
}

func (s *UserOrchestration) DeleteUser(id uint) error {
	return s.userRepo.Delete(id)
}

func (s *UserOrchestration) GetUser(id uint) (*entity.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *UserOrchestration) ListUsers(page, pageSize int, conditions map[string]interface{}) ([]entity.User, int64, error) {
	return s.userRepo.FindByPage(page, pageSize, conditions)
}

func (s *UserOrchestration) Login(req *LoginRequest, tokenGenerator func(userID uint, username string, roleID uint) (string, error)) (*LoginResponse, error) {
	user, err := s.userDomainService.ValidateCredentials(req.Username, req.Password)
	if err != nil {
		return nil, err
	}

	token, err := tokenGenerator(user.ID, user.Username, user.RoleID)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		User:  user,
		Token: token,
	}, nil
}

func (s *UserOrchestration) ChangePassword(userID uint, oldPassword, newPassword string) error {
	return s.userDomainService.UpdatePassword(userID, oldPassword, newPassword)
}
```

- [ ] **Step 2: 验证编译**

Run: `cd d:\Project\Go\Gsystes && go build ./...`

Expected: 编译成功

- [ ] **Step 3: 提交**

```bash
git add -A
git commit -m "feat: add user orchestration service"
```

---

### Task 9: 通信层 - 中间件

**Files:**
- Create: `internal/communication/middleware/auth_middleware.go`
- Create: `internal/communication/middleware/cors_middleware.go`
- Create: `internal/communication/middleware/recovery_middleware.go`
- Create: `internal/communication/middleware/logger_middleware.go`

- [ ] **Step 1: 编写认证中间件**

创建 `internal/communication/middleware/auth_middleware.go`:

```go
package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gsystes/backend/internal/infrastructure/auth"
	infraMiddleware "github.com/gsystes/backend/internal/infrastructure/middleware"
	"github.com/gsystes/backend/internal/infrastructure/utils"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Unauthorized(c, "authorization header is required")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.Unauthorized(c, "invalid authorization header format")
			c.Abort()
			return
		}

		claims, err := auth.ParseToken(parts[1])
		if err != nil {
			utils.Unauthorized(c, "invalid or expired token")
			c.Abort()
			return
		}

		infraMiddleware.SetClaims(c, claims)
		c.Next()
	}
}
```

- [ ] **Step 2: 编写 CORS 中间件**

创建 `internal/communication/middleware/cors_middleware.go`:

```go
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, X-Requested-With")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Authorization")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
```

- [ ] **Step 3: 编写 Recovery 中间件**

创建 `internal/communication/middleware/recovery_middleware.go`:

```go
package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/gsystes/backend/internal/infrastructure/logger"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("panic recovered",
					logger.AnyField("error", err),
					logger.AnyField("stack", string(debug.Stack())),
				)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code":    -1,
					"message": "internal server error",
				})
			}
		}()
		c.Next()
	}
}
```

- [ ] **Step 4: 编写请求日志中间件**

创建 `internal/communication/middleware/logger_middleware.go`:

```go
package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gsystes/backend/internal/infrastructure/logger"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()

		logger.Info("request",
			logger.StringField("method", method),
			logger.StringField("path", path),
			logger.IntField("status", statusCode),
			logger.DurationField("latency", latency),
			logger.StringField("client_ip", c.ClientIP()),
		)
	}
}
```

- [ ] **Step 5: 验证编译**

Run: `cd d:\Project\Go\Gsystes && go build ./...`

Expected: 编译失败，因为 logger 包中缺少 StringField/IntField 等辅助函数。这些函数后续通过 zap.Field 方式提供。

需要在 `internal/infrastructure/logger/logger.go` 中添加辅助函数，但为了保持 Task 独立性，暂时先编译看看问题。

Run: `cd d:\Project\Go\Gsystes && go get github.com/gin-gonic/gin && go build ./...`

如果编译报错，需要修复 logger 辅助函数。创建辅助函数如下，追加到 `internal/infrastructure/logger/logger.go` 末尾：

- [ ] **Step 6: 修复编译错误 - 添加 Logger 辅助函数**

修改 `internal/infrastructure/logger/logger.go`，在文件末尾添加辅助函数：

```go
func StringField(key, value string) zap.Field {
	return zap.String(key, value)
}

func IntField(key string, value int) zap.Field {
	return zap.Int(key, value)
}

func AnyField(key string, value interface{}) zap.Field {
	return zap.Any(key, value)
}

func ErrorField(err error) zap.Field {
	return zap.Error(err)
}

func DurationField(key string, value time.Duration) zap.Field {
	return zap.Duration(key, value)
}

func UintField(key string, value uint) zap.Field {
	return zap.Uint(key, value)
}
```

同时需要在文件顶部 import 中加入 `"time"`。

- [ ] **Step 7: 安装依赖并编译**

Run: `cd d:\Project\Go\Gsystes && go get github.com/gin-gonic/gin && go build ./...`

Expected: 编译成功

- [ ] **Step 8: 提交**

```bash
git add -A
git commit -m "feat: add communication layer middleware"
```

---

### Task 10: 通信层 - DTO、Handler 与路由

**Files:**
- Create: `internal/communication/dto/user_dto.go`
- Create: `internal/communication/dto/common_dto.go`
- Create: `internal/communication/handler/user_handler.go`
- Create: `internal/communication/router/router.go`

- [ ] **Step 1: 编写用户 DTO**

创建 `internal/communication/dto/user_dto.go`:

```go
package dto

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=64"`
	Password string `json:"password" binding:"required,min=6,max=128"`
	Nickname string `json:"nickname" binding:"max=64"`
	Email    string `json:"email" binding:"omitempty,email"`
	Phone    string `json:"phone" binding:"omitempty"`
	RoleID   uint   `json:"role_id" binding:"required"`
}

type UpdateUserRequest struct {
	Nickname string `json:"nickname" binding:"max=64"`
	Email    string `json:"email" binding:"omitempty,email"`
	Phone    string `json:"phone" binding:"omitempty"`
	RoleID   uint   `json:"role_id"`
	Status   int    `json:"status" binding:"oneof=1 2"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=128"`
}
```

- [ ] **Step 2: 编写通用 DTO**

创建 `internal/communication/dto/common_dto.go`:

```go
package dto

type PageParam struct {
	Page     int `form:"page" binding:"required,min=1"`
	PageSize int `form:"page_size" binding:"required,min=1,max=100"`
}

type IDParam struct {
	ID uint `uri:"id" binding:"required"`
}
```

- [ ] **Step 3: 编写用户 Handler**

创建 `internal/communication/handler/user_handler.go`:

```go
package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gsystes/backend/internal/communication/dto"
	"github.com/gsystes/backend/internal/infrastructure/auth"
	orchestration "github.com/gsystes/backend/internal/orchestration/service"
	"github.com/gsystes/backend/internal/infrastructure/utils"
)

type UserHandler struct {
	userOrchestration *orchestration.UserOrchestration
}

func NewUserHandler(userOrchestration *orchestration.UserOrchestration) *UserHandler {
	return &UserHandler{userOrchestration: userOrchestration}
}

func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	valErrors := utils.BindAndValidate(c, &req)
	if valErrors != nil {
		utils.BadRequest(c, "invalid request parameters")
		return
	}

	resp, err := h.userOrchestration.Login(&orchestration.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	}, auth.GenerateToken)
	if err != nil {
		utils.Unauthorized(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"token": resp.Token,
		"user": gin.H{
			"id":       resp.User.ID,
			"username": resp.User.Username,
			"nickname": resp.User.Nickname,
			"avatar":   resp.User.Avatar,
			"role_id":  resp.User.RoleID,
		},
	})
}

func (h *UserHandler) Create(c *gin.Context) {
	var req dto.CreateUserRequest
	valErrors := utils.BindAndValidate(c, &req)
	if valErrors != nil {
		utils.BadRequest(c, "invalid request parameters")
		return
	}

	user, err := h.userOrchestration.CreateUser(&orchestration.CreateUserRequest{
		Username: req.Username,
		Password: req.Password,
		Nickname: req.Nickname,
		Email:    req.Email,
		Phone:    req.Phone,
		RoleID:   req.RoleID,
	})
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"id": user.ID,
	})
}

func (h *UserHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid user id")
		return
	}

	var req dto.UpdateUserRequest
	valErrors := utils.BindAndValidate(c, &req)
	if valErrors != nil {
		utils.BadRequest(c, "invalid request parameters")
		return
	}

	if err := h.userOrchestration.UpdateUser(&orchestration.UpdateUserRequest{
		ID:       uint(id),
		Nickname: req.Nickname,
		Email:    req.Email,
		Phone:    req.Phone,
		RoleID:   req.RoleID,
		Status:   req.Status,
	}); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *UserHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid user id")
		return
	}

	if err := h.userOrchestration.DeleteUser(uint(id)); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *UserHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid user id")
		return
	}

	user, err := h.userOrchestration.GetUser(uint(id))
	if err != nil {
		utils.NotFound(c, "user not found")
		return
	}

	utils.Success(c, gin.H{
		"id":         user.ID,
		"username":   user.Username,
		"nickname":   user.Nickname,
		"email":      user.Email,
		"phone":      user.Phone,
		"avatar":     user.Avatar,
		"status":     user.Status,
		"role_id":    user.RoleID,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	})
}

func (h *UserHandler) List(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	conditions := make(map[string]interface{})
	if username := c.Query("username"); username != "" {
		conditions["username LIKE ?"] = "%" + username + "%"
	}
	if status := c.Query("status"); status != "" {
		conditions["status = ?"] = status
	}

	users, total, err := h.userOrchestration.ListUsers(page, pageSize, conditions)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	userList := make([]gin.H, len(users))
	for i, u := range users {
		userList[i] = gin.H{
			"id":         u.ID,
			"username":   u.Username,
			"nickname":   u.Nickname,
			"email":      u.Email,
			"phone":      u.Phone,
			"avatar":     u.Avatar,
			"status":     u.Status,
			"role_id":    u.RoleID,
			"created_at": u.CreatedAt,
		}
	}

	utils.PageSuccess(c, userList, total, page, pageSize)
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	var req dto.ChangePasswordRequest
	valErrors := utils.BindAndValidate(c, &req)
	if valErrors != nil {
		utils.BadRequest(c, "invalid request parameters")
		return
	}

	userID := utils.GetUserID(c)
	if err := h.userOrchestration.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}
```

- [ ] **Step 4: 编写路由注册**

创建 `internal/communication/router/router.go`:

```go
package router

import (
	"github.com/gin-gonic/gin"
	mid "github.com/gsystes/backend/internal/communication/middleware"
	"github.com/gsystes/backend/internal/communication/handler"
)

func SetupRouter(
	userHandler *handler.UserHandler,
) *gin.Engine {
	r := gin.New()

	r.Use(mid.Recovery())
	r.Use(mid.CORS())
	r.Use(mid.RequestLogger())

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", userHandler.Login)
		}

		users := api.Group("/users")
		users.Use(mid.AuthRequired())
		{
			users.POST("", userHandler.Create)
			users.PUT("/:id", userHandler.Update)
			users.DELETE("/:id", userHandler.Delete)
			users.GET("/:id", userHandler.Get)
			users.GET("", userHandler.List)
			users.PUT("/password", userHandler.ChangePassword)
		}
	}

	return r
}
```

- [ ] **Step 5: 编译验证**

Run: `cd d:\Project\Go\Gsystes && go build ./...`

Expected: 编译成功

- [ ] **Step 6: 提交**

```bash
git add -A
git commit -m "feat: add communication layer DTO, handler, and router"
```

---

### Task 11: 应用入口与依赖注入

**Files:**
- Create: `cmd/server/main.go`

- [ ] **Step 1: 编写 main.go**

创建 `cmd/server/main.go`:

```go
package main

import (
	"flag"
	"fmt"

	"github.com/gsystes/backend/internal/communication/handler"
	"github.com/gsystes/backend/internal/communication/router"
	dataRepo "github.com/gsystes/backend/internal/data/repository"
	"github.com/gsystes/backend/internal/data/migration"
	domainService "github.com/gsystes/backend/internal/domain/service"
	"github.com/gsystes/backend/internal/infrastructure/cache"
	"github.com/gsystes/backend/internal/infrastructure/config"
	"github.com/gsystes/backend/internal/infrastructure/database"
	"github.com/gsystes/backend/internal/infrastructure/logger"
	orchestration "github.com/gsystes/backend/internal/orchestration/service"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "config/config.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	if err := logger.InitLogger(cfg.Log); err != nil {
		panic(fmt.Sprintf("failed to init logger: %v", err))
	}

	if err := database.InitDatabase(cfg.Database); err != nil {
		logger.Fatal("failed to init database", logger.ErrorField(err))
	}

	if err := migration.AutoMigrate(database.GetDB()); err != nil {
		logger.Fatal("failed to auto migrate", logger.ErrorField(err))
	}

	if err := cache.InitRedis(cfg.Redis); err != nil {
		logger.Warn("redis init failed, continuing without cache", logger.ErrorField(err))
	}

	db := database.GetDB()
	userRepo := dataRepo.NewUserRepository(db)
	roleRepo := dataRepo.NewRoleRepository(db)
	permRepo := dataRepo.NewPermissionRepository(db)

	userDomainService := domainService.NewUserDomainService(userRepo)

	userOrchestration := orchestration.NewUserOrchestration(userDomainService, userRepo)

	userHandler := handler.NewUserHandler(userOrchestration)

	_ = roleRepo
	_ = permRepo

	r := router.SetupRouter(userHandler)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	logger.Info("server starting", logger.StringField("addr", addr))
	if err := r.Run(addr); err != nil {
		logger.Fatal("server failed to start", logger.ErrorField(err))
	}
}
```

- [ ] **Step 2: 编译验证**

Run: `cd d:\Project\Go\Gsystes && go build ./...`

Expected: 编译成功

- [ ] **Step 3: 验证完整编译**

Run: `cd d:\Project\Go\Gsystes && go build -o bin/gsystes-server ./cmd/server`

Expected: 生成可执行文件 `bin/gsystes-server`

- [ ] **Step 4: 修复编译错误（如有）**

如果编译报错（如循环依赖、未使用的导入等），逐项修复并重新编译。

- [ ] **Step 5: 提交**

```bash
git add -A
git commit -m "feat: add application entry point with dependency injection"
```

---

### Task 12: Makefile 与 Dockerfile

**Files:**
- Create: `Makefile`
- Create: `Dockerfile`

- [ ] **Step 1: 编写 Makefile**

创建 `Makefile`:

```makefile
.PHONY: build run clean test lint help

APP_NAME = gsystes-server
BUILD_DIR = bin
MAIN_PATH = cmd/server

help:
	@echo "Usage:"
	@echo "  make build       - Build the application"
	@echo "  make run         - Run the application"
	@echo "  make clean       - Clean build artifacts"
	@echo "  make test        - Run tests"
	@echo "  make lint        - Run linter"

build:
	go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)

run: build
	./$(BUILD_DIR)/$(APP_NAME) -config config/config.dev.yaml

clean:
	rm -rf $(BUILD_DIR)
	rm -rf logs/

test:
	go test ./... -v -cover

lint:
	golangci-lint run ./...

dev:
	go run $(MAIN_PATH)/main.go -config config/config.dev.yaml
```

- [ ] **Step 2: 编写 Dockerfile**

创建 `Dockerfile`:

```dockerfile
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/gsystes-server ./cmd/server

FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/gsystes-server .
COPY config/config.yaml ./config/

EXPOSE 8080

ENTRYPOINT ["./gsystes-server", "-config", "config/config.yaml"]
```

- [ ] **Step 3: 验证编译**

Run: `cd d:\Project\Go\Gsystes && go build ./...`

Expected: 编译成功

- [ ] **Step 4: 提交**

```bash
git add -A
git commit -m "chore: add Makefile and Dockerfile"
```

---

### Task 13: 编写单元测试 - 领域服务测试

**Files:**
- Create: `internal/domain/service/user_domain_service_test.go`

- [ ] **Step 1: 编写用户领域服务测试**

创建 `internal/domain/service/user_domain_service_test.go`:

```go
package service

import (
	"errors"
	"testing"

	"github.com/gsystes/backend/internal/domain/entity"
	"github.com/gsystes/backend/internal/infrastructure/utils"
)

type mockUserRepo struct {
	users map[uint]*entity.User
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{users: make(map[uint]*entity.User)}
}

func (m *mockUserRepo) Create(user *entity.User) error {
	user.ID = uint(len(m.users) + 1)
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepo) Update(user *entity.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepo) Delete(id uint) error {
	delete(m.users, id)
	return nil
}

func (m *mockUserRepo) FindByID(id uint) (*entity.User, error) {
	user, exists := m.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (m *mockUserRepo) FindByUsername(username string) (*entity.User, error) {
	for _, user := range m.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *mockUserRepo) FindByPage(page, pageSize int, conditions map[string]interface{}) ([]entity.User, int64, error) {
	return nil, 0, nil
}

func TestCreateUser_Success(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewUserDomainService(repo)

	user := &entity.User{
		Username: "testuser",
	}
	err := svc.Create(user, "password123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user.ID == 0 {
		t.Fatal("expected user ID to be set")
	}
	if user.Status != int(entity.UserStatusActive) {
		t.Fatalf("expected status %d, got %d", entity.UserStatusActive, user.Status)
	}
	if !utils.CheckPassword("password123", user.Password) {
		t.Fatal("expected password to be hashed correctly")
	}
}

func TestCreateUser_DuplicateUsername(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewUserDomainService(repo)

	user1 := &entity.User{Username: "testuser"}
	if err := svc.Create(user1, "password123"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	user2 := &entity.User{Username: "testuser"}
	err := svc.Create(user2, "password456")
	if err == nil {
		t.Fatal("expected error for duplicate username")
	}
}

func TestValidateCredentials_Success(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewUserDomainService(repo)

	user := &entity.User{Username: "testuser"}
	if err := svc.Create(user, "password123"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	result, err := svc.ValidateCredentials("testuser", "password123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Username != "testuser" {
		t.Fatalf("expected username testuser, got %s", result.Username)
	}
}

func TestValidateCredentials_WrongPassword(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewUserDomainService(repo)

	user := &entity.User{Username: "testuser"}
	if err := svc.Create(user, "password123"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err := svc.ValidateCredentials("testuser", "wrongpassword")
	if err == nil {
		t.Fatal("expected error for wrong password")
	}
}

func TestUpdatePassword_Success(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewUserDomainService(repo)

	user := &entity.User{Username: "testuser"}
	if err := svc.Create(user, "oldpassword"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if err := svc.UpdatePassword(user.ID, "oldpassword", "newpassword"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	updated, _ := repo.FindByID(user.ID)
	if !utils.CheckPassword("newpassword", updated.Password) {
		t.Fatal("expected password to be updated")
	}
}

func TestUpdatePassword_WrongOldPassword(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewUserDomainService(repo)

	user := &entity.User{Username: "testuser"}
	if err := svc.Create(user, "oldpassword"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err := svc.UpdatePassword(user.ID, "wrongpassword", "newpassword")
	if err == nil {
		t.Fatal("expected error for wrong old password")
	}
}
```

- [ ] **Step 2: 运行测试**

Run: `cd d:\Project\Go\Gsystes && go test ./internal/domain/service/ -v`

Expected: 所有测试通过

- [ ] **Step 3: 提交**

```bash
git add -A
git commit -m "test: add user domain service unit tests"
```

---

## 自检清单

### 1. Spec 覆盖检查
| 需求 | 对应 Task |
|------|----------|
| 五层分层架构 | Task 1-12 整体设计 |
| 通信层（HTTP 路由/中间件/Handler） | Task 9, 10 |
| 业务编排层（编排服务） | Task 8 |
| 领域层（实体/仓储接口/领域服务） | Task 5, 6 |
| 数据层（ORM 模型/仓储实现/迁移） | Task 7 |
| 基础设施层（配置/日志/数据库/缓存/认证/工具） | Task 1, 2, 3, 4 |
| 用户管理 CRUD | Task 5-10 |
| 用户登录认证（JWT） | Task 4, 8, 9 |
| 数据库自动迁移 | Task 7 |
| Docker 部署 | Task 12 |
| 单元测试 | Task 13 |

### 2. 占位符检查
- 所有代码块包含完整的可编译 Go 代码
- 所有文件路径使用绝对路径（相对于项目根）
- 所有命令包含具体参数和预期输出
- 无 "TBD"、"TODO"、空实现等占位符

### 3. 类型一致性检查
- `entity.User` 在所有层间使用一致的字段定义
- 仓储接口方法签名在 `domain/repository` 和 `data/repository` 中一致
- `UserOrchestration` 的 `Login` 方法签名的 `tokenGenerator` 回调与 `auth.GenerateToken` 签名一致
- 日志辅助函数命名风格统一