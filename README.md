# App 基础后端框架

基于 Go + Gin + GORM 的单体服务后端基础架构，包含用户认证、RBAC权限管理等基础功能。

## 技术栈

- Go 1.23+
- Gin (Web框架)
- GORM + Gen (ORM)
- MySQL
- Redis
- Viper (配置管理，支持etcd热更新)
- Zap (日志)

## 项目结构

```
├── adaptor/                # 适配层
│   ├── redis/              # Redis操作封装
│   └── repo/               # 数据库操作
│       ├── admin/          # 业务仓储实现
│       ├── model/          # GORM模型 (gen生成)
│       └── query/          # 查询构建器 (gen生成)
├── api/                    # API层 (Controller)
│   ├── admin/              # 管理后台接口
│   └── customer/           # 用户端接口
├── common/                 # 公共定义
│   ├── auth.go             # 用户结构体
│   └── errno.go            # 错误码
├── config/                 # 配置
├── consts/                 # 常量
├── router/                 # 路由 & 中间件
├── service/                # 业务逻辑层
│   ├── admin/              # 管理后台服务
│   ├── do/                 # 领域对象
│   └── dto/                # 数据传输对象
├── utils/                  # 工具类
│   ├── captcha/            # 滑块验证码
│   ├── logger/             # 日志
│   └── tools/              # 通用工具
├── app_local.yml           # 本地配置文件
└── main.go                 # 入口
```

## 快速开始

### 1. 配置数据库

创建数据库并导入基础表结构（admin_user, roles, permission等RBAC相关表）。

### 2. 修改配置

编辑 `app_local.yml`：

```yaml
server:
  http_port: 8900
  log_level: debug

mysql:
  host: 127.0.0.1
  port: 3306
  user: root
  password: 123456
  database: app

redis:
  addr: 127.0.0.1:6379
```

### 3. 运行

```bash
# 使用本地配置
go run main.go -c app_local.yml

# 使用etcd配置
ETCD_ADDR=127.0.0.1:2379 go run main.go
```

### 4. 生成数据库模型

```bash
cd adaptor/repo
go run gorm.io/gen/tools/gentool -c gen.yaml
```

## API接口

| 路径 | 方法 | 说明 | 鉴权 |
|------|------|------|------|
| /ping | ANY | 健康检查 | 否 |
| /api/app/admin/v1/user/verify/captcha | GET | 获取滑块验证码 | 否 |
| /api/app/admin/v1/user/verify/captcha/check | POST | 校验滑块验证码 | 否 |
| /api/app/admin/v1/user/info | GET | 获取用户信息 | 是 |
| /api/app/admin/v1/user/create | POST | 创建用户 | 是 |
| /api/app/admin/v1/user/update | POST | 更新用户 | 是 |

## 特性

- 支持本地配置文件和etcd远程配置（热更新）
- 内置访问日志中间件
- 内置JWT鉴权中间件（需实现token解析逻辑）
- 内置滑块验证码
- 支持pprof性能分析
- 优雅关闭
