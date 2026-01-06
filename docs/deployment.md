# 部署文档

## 环境要求

| 组件 | 最低版本 | 推荐版本 |
|------|----------|----------|
| Go | 1.24 | 1.24.x |
| MySQL | 5.7 | 8.0+ |
| Redis | 6.0 | 7.0+ |

## 一、准备工作

### 1.1 安装 Go

```bash
# Windows (使用 Chocolatey)
choco install golang

# macOS (使用 Homebrew)
brew install go@1.24

# Linux (Ubuntu/Debian)
wget https://go.dev/dl/go1.24.0.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.24.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

验证安装：
```bash
go version
```

### 1.2 安装 MySQL

```bash
# macOS
brew install mysql@8.0
brew services start mysql@8.0

# Ubuntu/Debian
sudo apt update
sudo apt install mysql-server-8.0
sudo systemctl start mysql
sudo systemctl enable mysql

# Windows
# 下载安装包: https://dev.mysql.com/downloads/mysql/
```

初始化数据库：
```bash
sudo mysql
CREATE DATABASE FindLink DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'findlink'@'localhost' IDENTIFIED BY 'your_password';
GRANT ALL PRIVILEGES ON FindLink.* TO 'findlink'@'localhost';
FLUSH PRIVILEGES;
EXIT;
```

### 1.3 安装 Redis

```bash
# macOS
brew install redis
brew services start redis

# Ubuntu/Debian
sudo apt install redis-server
sudo systemctl start redis
sudo systemctl enable redis

# Windows
# 下载: https://github.com/microsoftarchive/redis/releases
```

## 二、编译部署

### 2.1 获取代码

```bash
git clone https://your-repo/FindLink.git
cd FindLink
```

### 2.2 编译二进制文件

```bash
# Linux/macOS
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o FindLink main.go

# Windows
go build -o FindLink.exe main.go
```

### 2.3 交叉编译示例

```bash
# 编译不同平台版本
# Linux AMD64
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o FindLink-linux-amd64 main.go

# Linux ARM64
CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build -o FindLink-linux-arm64 main.go

# macOS AMD64
CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o FindLink-darwin-amd64 main.go

# macOS ARM64 (Apple Silicon)
CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -o FindLink-darwin-arm64 main.go

# Windows AMD64
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -o FindLink.exe main.go
```

## 三、配置说明

### 3.1 配置文件结构

创建 `config.yml`：

```yaml
server:
  http_port: 8900
  enable_pprof: true
  log_level: info
  env: production

mysql:
  dialect: mysql
  user: findlink
  password: your_password
  host: 127.0.0.1
  port: 3306
  database: FindLink
  charset: utf8mb4
  show_mysql: false
  max_open: 50
  max_idle: 10

redis:
  addr: 127.0.0.1:6379
  password: your_password
  db_index: 0
  max_idle: 10
  max_open: 50

# 可选: etcd 配置（远程配置中心）
etcd:
  enabled: false
  endpoints:
    - 127.0.0.1:2379
  timeout: 5
```

### 3.2 环境变量覆盖

支持通过环境变量覆盖配置：

| 配置项 | 环境变量 |
|--------|----------|
| HTTP 端口 | `SERVER_HTTP_PORT` |
| 日志级别 | `SERVER_LOG_LEVEL` |
| MySQL 主机 | `MYSQL_HOST` |
| MySQL 密码 | `MYSQL_PASSWORD` |
| Redis 地址 | `REDIS_ADDR` |
| Redis 密码 | `REDIS_PASSWORD` |

### 3.3 日志配置

日志输出到文件和控制台：

```yaml
server:
  log_level: info
  # 日志文件配置在代码中实现
```

## 四、运行服务

### 4.1 前台运行

```bash
# 使用配置文件
./FindLink -c config.yml

# 使用环境变量
./FindLink
```

### 4.2 后台运行 (Systemd)

创建服务文件 `/etc/systemd/system/findlink.service`：

```ini
[Unit]
Description=FindLink Location Service
After=network.target mysql.service redis.service
Wants=mysql.service redis.service

[Service]
Type=simple
User=www-data
Group=www-data
WorkingDirectory=/opt/findlink
ExecStart=/opt/findlink/FindLink -c /opt/findlink/config.yml
Restart=always
RestartSec=5

# 环境变量
Environment="SERVER_ENV=production"
EnvironmentFile=/opt/findlink/.env

# 日志配置
StandardOutput=journal
StandardError=journal
SyslogIdentifier=findlink

# 资源限制
LimitNOFILE=65536
LimitNPROC=4096

[Install]
WantedBy=multi-user.target
```

启动服务：

```bash
sudo systemctl daemon-reload
sudo systemctl enable findlink
sudo systemctl start findlink
sudo systemctl status findlink
```

### 4.3 后台运行 (Docker)

创建 `Dockerfile`:

```dockerfile
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o FindLink main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /opt/findlink

COPY --from=builder /app/FindLink .
COPY --from=builder /app/app_local.yml .

EXPOSE 8900

CMD ["sh", "-c", "./FindLink -c app_local.yml"]
```

构建和运行：

```bash
# 构建镜像
docker build -t findlink:latest .

# 运行容器
docker run -d \
  --name findlink \
  -p 8900:8900 \
  -v $(pwd)/config.yml:/opt/findlink/config.yml \
  -e SERVER_ENV=production \
  findlink:latest
```

### 4.4 Docker Compose

创建 `docker-compose.yml`:

```yaml
version: '3.8'

services:
  mysql:
    image: mysql:8.0
    container_name: findlink-mysql
    restart: always
    environment:
      MYSQL_ROOT_PASSWORD: root_password
      MYSQL_DATABASE: FindLink
      MYSQL_USER: findlink
      MYSQL_PASSWORD: findlink_password
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql
      - ./migrations:/docker-entrypoint-initdb.d
    command: --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost"]
      interval: 10s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    container_name: findlink-redis
    restart: always
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    command: redis-server --appendonly yes
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5

  findlink:
    build: .
    container_name: findlink-app
    restart: always
    ports:
      - "8900:8900"
    depends_on:
      mysql:
        condition: service_healthy
      redis:
        condition: service_healthy
    environment:
      - SERVER_ENV=production
      - MYSQL_HOST=mysql
      - MYSQL_PASSWORD=findlink_password
      - REDIS_ADDR=redis:6379
    volumes:
      - ./config.yml:/opt/findlink/config.yml
      - logs:/opt/findlink/logs

volumes:
  mysql_data:
  redis_data:
  logs:
```

启动服务：

```bash
docker-compose up -d
docker-compose logs -f findlink
```

## 五、监控运维

### 5.1 健康检查

```bash
# HTTP 健康检查
curl http://localhost:8900/ping

# 响应示例
{"message":"pong"}
```

### 5.2 pprof 性能分析

```bash
# 查看 CPU 使用情况
go tool pprof http://localhost:8900/debug/pprof/profile

# 查看内存使用情况
go tool pprof http://localhost:8900/debug/pprof/heap

# 查看 goroutine
go tool pprof http://localhost:8900/debug/pprof/goroutine
```

### 5.3 日志查看

```bash
# 查看系统日志 (journald)
journalctl -u findlink -f

# 查看 Docker 日志
docker logs -f findlink
```

### 5.4 优雅重启

```bash
# Systemd 重启
sudo systemctl restart findlink

# Docker Compose 重启
docker-compose restart findlink

# 发送 SIGHUP 信号（配置热重载）
kill -HUP $(cat /var/run/findlink.pid)
```

## 六、故障排查

### 6.1 常见问题

**Q: MySQL 连接失败**
```bash
# 检查 MySQL 服务状态
sudo systemctl status mysql

# 测试连接
mysql -u findlink -p -h 127.0.0.1
```

**Q: Redis 连接失败**
```bash
# 检查 Redis 服务状态
sudo systemctl status redis

# 测试连接
redis-cli ping
```

**Q: 端口被占用**
```bash
# 查看端口占用
netstat -tlnp | grep 8900

# 杀掉占用进程
kill -9 <pid>
```

### 6.2 日志级别调试

在排查问题时，可以临时降低日志级别：

```yaml
server:
  log_level: debug
```

## 七、性能优化

### 7.1 MySQL 优化

```sql
-- 优化表结构
OPTIMIZE TABLE user_locations;
OPTIMIZE TABLE device_locations;

-- 查看慢查询
SHOW VARIABLES LIKE 'slow_query_log';
SET GLOBAL slow_query_log = 'ON';
```

### 7.2 Redis 优化

```bash
# 查看内存使用
redis-cli info memory

# 内存碎片整理
redis-cli memory purge
```

### 7.3 连接池配置

```yaml
mysql:
  max_open: 100
  max_idle: 20

redis:
  max_open: 100
  max_idle: 20
```
