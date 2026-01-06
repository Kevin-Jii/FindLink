# FindLink - 位置共享服务后端

基于 Go + Gin + GORM 的实时位置共享应用后端服务，支持位置上报、好友管理、设备追踪、地理围栏和 WebSocket 实时推送。

## 功能特性

### C端用户功能
- **用户认证**: 注册、登录、JWT 令牌管理
- **位置共享**: 实时位置上报、历史轨迹查询、附近好友发现
- **好友系统**: 添加好友、好友请求管理、双向位置共享
- **设备管理**: 绑定/解绑设备、设备状态监控
- **地理围栏**: 创建/管理地理围栏、进出区域提醒
- **实时通信**: WebSocket 长连接、位置实时推送

### B端管理后台
- **RBAC 权限管理**: 基于角色的访问控制
- **用户管理**: 管理员账号管理
- **验证码**: 滑块验证码安全防护

## 技术栈

| 组件 | 技术选型 | 用途 |
|------|----------|------|
| 语言 | Go 1.24+ | 后端开发 |
| Web框架 | Gin | HTTP API |
| ORM | GORM + Gen | 数据库操作 |
| 数据库 | MySQL | 主数据存储 |
| 缓存 | Redis | 位置缓存、Pub/Sub |
| 配置 | Viper | 配置管理（支持 etcd 热更新） |
| 日志 | Zap | 结构化日志 |
| API文档 | Swagger | 接口文档 |
| 实时通信 | gorilla/websocket | WebSocket |

## 项目结构

```
FindLink/
├── adaptor/                  # 适配层
│   ├── repo/                 # 数据访问层
│   │   ├── admin/            # 管理后台仓储
│   │   ├── device/           # 设备仓储
│   │   ├── friend/           # 好友仓储
│   │   ├── geofence/         # 地理围栏仓储
│   │   ├── location/         # 位置仓储
│   │   ├── model/            # GORM 数据模型
│   │   ├── query/            # GORM Gen 查询构建器
│   │   └── user/             # 用户仓储
│   └── redis/                # Redis 操作封装
│       ├── cache.go          # 通用缓存
│       ├── location_cache.go # 位置缓存
│       ├── pubsub.go         # 发布订阅
│       └── verify.go         # 验证码存储
├── api/                      # API 控制器层
│   ├── admin/                # 管理后台接口
│   │   ├── admin.go          # 管理后台入口
│   │   ├── login.go          # 登录接口
│   │   ├── perm.go           # 权限接口
│   │   ├── role.go           # 角色接口
│   │   └── user.go           # 用户管理接口
│   └── customer/             # C端用户接口
│       ├── customer.go       # C端入口
│       ├── device.go         # 设备接口
│       ├── friend.go         # 好友接口
│       ├── geofence.go       # 地理围栏接口
│       ├── location.go       # 位置接口
│       ├── user.go           # 用户接口
│       └── websocket.go      # WebSocket 接口
├── common/                   # 公共定义
│   ├── auth.go               # 用户结构体定义
│   └── errno.go              # 错误码定义
├── config/                   # 配置管理
│   └── config.go             # 配置初始化
├── consts/                   # 常量定义
├── migrations/               # 数据库迁移脚本
├── router/                   # 路由 & 中间件
│   ├── access.go             # 访问日志中间件
│   ├── auth.go               # 认证中间件
│   ├── pprof.go              # 性能分析路由
│   ├── router.go             # 主路由配置
│   └── white_list.go         # 白名单配置
├── service/                  # 业务逻辑层
│   ├── admin/                # 管理后台服务
│   │   ├── login.go          # 登录服务
│   │   ├── service.go        # 服务入口
│   │   ├── user.go           # 用户服务
│   │   └── verify.go         # 验证码服务
│   ├── device/               # 设备服务
│   ├── friend/               # 好友服务
│   ├── geofence/             # 地理围栏服务
│   ├── location/             # 位置服务
│   ├── user/                 # 用户服务
│   ├── websocket/            # WebSocket 服务
│   ├── do/                   # 领域对象
│   └── dto/                  # 数据传输对象
├── utils/                    # 工具类
│   ├── captcha/              # 滑块验证码
│   ├── logger/               # 日志封装
│   └── tools/                # 通用工具
├── docs/                     # Swagger 文档
├── app_local.yml             # 本地配置示例
└── main.go                   # 应用入口
```

## 快速开始

### 环境要求

- Go 1.24+
- MySQL 5.7+ / 8.0+
- Redis 6.0+

### 1. 配置数据库

```sql
CREATE DATABASE FindLink DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 2. 修改配置

编辑 `app_local.yml`：

```yaml
server:
  http_port: 8900
  enable_pprof: true
  log_level: debug
  env: dev

mysql:
  dialect: mysql
  user: root
  password: your_password
  host: 127.0.0.1
  port: 3306
  database: FindLink
  charset: utf8mb4
  show_mysql: true
  max_open: 20
  max_idle: 5

redis:
  addr: 127.0.0.1:6379
  password: your_password
  db_index: 0
  max_idle: 2
  max_open: 10
```

### 3. 运行服务

```bash
# 使用本地配置运行
go run main.go -c app_local.yml

# 使用 etcd 远程配置运行
ETCD_ADDR=127.0.0.1:2379 go run main.go
```

### 4. 生成数据库模型

```bash
make gendb
# 或手动执行
cd adaptor/repo && gentool -c gen.yaml
```

### 5. 访问 API 文档

启动服务后访问：`http://localhost:8900/swagger/index.html`

## API 接口

### C端用户接口

| 路径 | 方法 | 说明 | 鉴权 |
|------|------|------|------|
| `/api/app/customer/v1/user/login` | POST | 用户登录 | 否 |
| `/api/app/customer/v1/user/register` | POST | 用户注册 | 否 |
| `/api/app/customer/v1/user/info` | GET | 获取用户信息 | 是 |
| `/api/app/customer/v1/user/logout` | POST | 用户登出 | 是 |
| `/api/app/customer/v1/location/report` | POST | 上报位置 | 是 |
| `/api/app/customer/v1/location/batch` | POST | 批量上报位置 | 是 |
| `/api/app/customer/v1/location/user/:user_id` | GET | 获取用户位置 | 是 |
| `/api/app/customer/v1/location/device/:device_id` | GET | 获取设备位置 | 是 |
| `/api/app/customer/v1/location/history` | GET | 获取位置历史 | 是 |
| `/api/app/customer/v1/location/nearby` | GET | 获取附近好友 | 是 |
| `/api/app/customer/v1/friend/request` | POST | 发送好友请求 | 是 |
| `/api/app/customer/v1/friend/requests` | GET | 获取好友请求列表 | 是 |
| `/api/app/customer/v1/friend/accept` | POST | 接受好友请求 | 是 |
| `/api/app/customer/v1/friend/reject` | POST | 拒绝好友请求 | 是 |
| `/api/app/customer/v1/friend/list` | GET | 获取好友列表 | 是 |
| `/api/app/customer/v1/friend/:friend_id` | DELETE | 删除好友 | 是 |
| `/api/app/customer/v1/friend/search` | GET | 搜索用户 | 是 |
| `/api/app/customer/v1/device/bind` | POST | 绑定设备 | 是 |
| `/api/app/customer/v1/device/list` | GET | 获取设备列表 | 是 |
| `/api/app/customer/v1/device/:device_id/settings` | PUT | 更新设备设置 | 是 |
| `/api/app/customer/v1/device/:device_id/status` | PUT | 更新设备状态 | 是 |
| `/api/app/customer/v1/device/:device_id` | DELETE | 解绑设备 | 是 |
| `/api/app/customer/v1/geofence` | POST | 创建地理围栏 | 是 |
| `/api/app/customer/v1/geofence/list` | GET | 获取地理围栏列表 | 是 |
| `/api/app/customer/v1/geofence/:geofence_id` | PUT | 更新地理围栏 | 是 |
| `/api/app/customer/v1/geofence/:geofence_id` | DELETE | 删除地理围栏 | 是 |
| `/api/app/customer/v1/ws` | GET | WebSocket 连接 | 是 |

### B端管理后台接口

| 路径 | 方法 | 说明 | 鉴权 |
|------|------|------|------|
| `/api/app/admin/v1/user/login` | POST | 管理员登录 | 否 |
| `/api/app/admin/v1/user/verify/captcha` | GET | 获取滑块验证码 | 否 |
| `/api/app/admin/v1/user/verify/captcha/check` | POST | 校验滑块验证码 | 否 |
| `/api/app/admin/v1/user/info` | GET | 获取管理员信息 | 是 |
| `/api/app/admin/v1/user/create` | POST | 创建管理员用户 | 是 |
| `/api/app/admin/v1/user/update` | POST | 更新管理员用户 | 是 |

### 系统接口

| 路径 | 方法 | 说明 | 鉴权 |
|------|------|------|------|
| `/ping` | ANY | 健康检查 | 否 |
| `/swagger/*any` | GET | Swagger 文档 | 否 |
| `/debug/pprof/*` | GET | 性能分析 | 否 |

## 默认账号

- **超级管理员**: `admin` / `admin123`

## 特性

- 支持本地配置文件和 etcd 远程配置（热更新）
- 内置访问日志中间件
- 内置 JWT 认证中间件
- 内置滑块验证码
- 支持 pprof 性能分析
- 自动数据库迁移
- WebSocket 实时位置推送
- Redis 位置缓存加速
- Swagger API 文档自动生成

## 文档

- [API 接口文档](docs/api.md)
- [数据库设计文档](docs/database.md)
- [部署文档](docs/deployment.md)
