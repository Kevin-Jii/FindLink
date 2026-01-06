# API 接口文档

## 概述

- Base URL: `http://localhost:8900/api/app`
- 认证方式: JWT Token (Header: `Authorization: Bearer <token>`)
- 响应格式: JSON

## 通用响应结构

### 成功响应
```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

### 错误响应
```json
{
  "code": 10001,
  "message": "invalid parameters",
  "data": null
}
```

### 错误码定义

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 10001 | 参数错误 |
| 10002 | 未授权 |
| 10003 | 禁止访问 |
| 10004 | 资源不存在 |
| 10005 | 内部错误 |

---

## 一、用户认证

### 1.1 用户登录

**POST** `/customer/v1/user/login`

**请求参数**
```json
{
  "phone": "13800138000",
  "password": "password123"
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| phone | string | 是 | 手机号 |
| password | string | 是 | 密码 |

**响应**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user_id": 10001,
    "expires_in": 86400
  }
}
```

---

### 1.2 用户注册

**POST** `/customer/v1/user/register`

**请求参数**
```json
{
  "phone": "13800138000",
  "password": "password123",
  "nickname": "用户名",
  "verify_code": "123456"
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| phone | string | 是 | 手机号 |
| password | string | 是 | 密码 (6-20位) |
| nickname | string | 是 | 昵称 |
| verify_code | string | 是 | 短信验证码 |

---

### 1.3 获取用户信息

**GET** `/customer/v1/user/info`

**请求头**
```
Authorization: Bearer <token>
```

**响应**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "user_id": 10001,
    "phone": "13800138000",
    "nickname": "用户名",
    "avatar": "https://example.com/avatar.png",
    "sex": 1,
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

---

### 1.4 用户登出

**POST** `/customer/v1/user/logout`

**请求头**
```
Authorization: Bearer <token>
```

---

## 二、位置服务

### 2.1 上报位置

**POST** `/customer/v1/location/report`

**请求头**
```
Authorization: Bearer <token>
```

**请求参数**
```json
{
  "longitude": 116.397428,
  "latitude": 39.90923,
  "accuracy": 10.5,
  "altitude": 100.0,
  "speed": 0.0,
  "bearing": 0.0,
  "battery_level": 85,
  "location_mode": "foreground"
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| longitude | float | 是 | 经度 |
| latitude | float | 是 | 纬度 |
| accuracy | float | 否 | 精度 (米) |
| altitude | float | 否 | 海拔 (米) |
| speed | float | 否 | 速度 (米/秒) |
| bearing | float | 否 | 方向 (0-360) |
| battery_level | int | 否 | 电量 (0-100) |
| location_mode | string | 否 | 定位模式 (foreground/background/significant_change) |

---

### 2.2 批量上报位置

**POST** `/customer/v1/location/batch`

**请求头**
```
Authorization: Bearer <token>
```

**请求参数**
```json
{
  "locations": [
    {
      "longitude": 116.397428,
      "latitude": 39.90923,
      "accuracy": 10.5,
      "timestamp": 1704067200000
    },
    {
      "longitude": 116.397429,
      "latitude": 39.90924,
      "accuracy": 12.0,
      "timestamp": 1704067201000
    }
  ]
}
```

---

### 2.3 获取用户位置

**GET** `/customer/v1/location/user/:user_id`

**请求头**
```
Authorization: Bearer <token>
```

**响应**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "user_id": 10001,
    "longitude": 116.397428,
    "latitude": 39.90923,
    "accuracy": 10.5,
    "updated_at": "2024-01-01T00:00:00Z",
    "is_online": true
  }
}
```

---

### 2.4 获取设备位置

**GET** `/customer/v1/location/device/:device_id`

**请求头**
```
Authorization: Bearer <token>
```

---

### 2.5 获取位置历史

**GET** `/customer/v1/location/history`

**请求参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| start_time | int64 | 是 | 开始时间 (毫秒时间戳) |
| end_time | int64 | 是 | 结束时间 (毫秒时间戳) |
| limit | int | 否 | 返回数量限制 (默认100) |

**响应**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "locations": [
      {
        "longitude": 116.397428,
        "latitude": 39.90923,
        "accuracy": 10.5,
        "timestamp": 1704067200000
      }
    ],
    "total": 100
  }
}
```

---

### 2.6 获取附近好友

**GET** `/customer/v1/location/nearby`

**请求参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| radius | float | 否 | 搜索半径 (米, 默认1000) |
| limit | int | 否 | 返回数量限制 (默认20) |

**响应**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "friends": [
      {
        "user_id": 10002,
        "nickname": "好友昵称",
        "avatar": "https://example.com/avatar.png",
        "longitude": 116.397428,
        "latitude": 39.90923,
        "distance": 500.5,
        "updated_at": "2024-01-01T00:00:00Z"
      }
    ]
  }
}
```

---

## 三、好友服务

### 3.1 发送好友请求

**POST** `/customer/v1/friend/request`

**请求头**
```
Authorization: Bearer <token>
```

**请求参数**
```json
{
  "target_user_id": 10002,
  "message": "我是xxx"
}
```

---

### 3.2 获取好友请求列表

**GET** `/customer/v1/friend/requests`

**请求参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| type | string | 否 | 请求类型 (received/sent, 默认received) |
| status | string | 否 | 状态 (pending/accepted/rejected) |

**响应**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "requests": [
      {
        "request_id": 10001,
        "from_user_id": 10002,
        "from_user_nickname": "用户A",
        "from_user_avatar": "https://example.com/avatar.png",
        "message": "我是xxx",
        "status": "pending",
        "created_at": "2024-01-01T00:00:00Z"
      }
    ]
  }
}
```

---

### 3.3 接受好友请求

**POST** `/customer/v1/friend/accept`

**请求参数**
```json
{
  "request_id": 10001
}
```

---

### 3.4 拒绝好友请求

**POST** `/customer/v1/friend/reject`

**请求参数**
```json
{
  "request_id": 10001
}
```

---

### 3.5 获取好友列表

**GET** `/customer/v1/friend/list`

**响应**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "friends": [
      {
        "friend_id": 10002,
        "nickname": "好友昵称",
        "avatar": "https://example.com/avatar.png",
        "sharing_status": "sharing",
        "last_location": {
          "longitude": 116.397428,
          "latitude": 39.90923,
          "updated_at": "2024-01-01T00:00:00Z"
        }
      }
    ]
  }
}
```

---

### 3.6 删除好友

**DELETE** `/customer/v1/friend/:friend_id`

---

### 3.7 搜索用户

**GET** `/customer/v1/friend/search`

**请求参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| keyword | string | 是 | 搜索关键词 |

---

## 四、设备服务

### 4.1 绑定设备

**POST** `/customer/v1/device/bind`

**请求参数**
```json
{
  "device_id": "ABC123",
  "name": "我的设备",
  "type": "gps_tracker"
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| device_id | string | 是 | 设备ID |
| name | string | 否 | 设备名称 |
| type | string | 否 | 设备类型 (gps_tracker/smart_watch/other) |

---

### 4.2 获取设备列表

**GET** `/customer/v1/device/list`

**响应**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "devices": [
      {
        "device_id": "ABC123",
        "name": "我的设备",
        "type": "gps_tracker",
        "battery_level": 85,
        "connection_status": "online",
        "last_location": {
          "longitude": 116.397428,
          "latitude": 39.90923,
          "updated_at": "2024-01-01T00:00:00Z"
        },
        "created_at": "2024-01-01T00:00:00Z"
      }
    ]
  }
}
```

---

### 4.3 更新设备设置

**PUT** `/customer/v1/device/:device_id/settings`

**请求参数**
```json
{
  "name": "新名称",
  "notification_enabled": true
}
```

---

### 4.4 更新设备状态

**PUT** `/customer/v1/device/:device_id/status`

**请求参数**
```json
{
  "battery_level": 75,
  "connection_status": "online"
}
```

---

### 4.5 解绑设备

**DELETE** `/customer/v1/device/:device_id`

---

## 五、地理围栏服务

### 5.1 创建地理围栏

**POST** `/customer/v1/geofence`

**请求参数**
```json
{
  "name": "家",
  "longitude": 116.397428,
  "latitude": 39.90923,
  "radius": 100.0,
  "notify_on_enter": true,
  "notify_on_exit": true
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 围栏名称 |
| longitude | float | 是 | 中心点经度 |
| latitude | float | 是 | 中心点纬度 |
| radius | float | 是 | 半径 (米) |
| notify_on_enter | bool | 否 | 进入通知 (默认true) |
| notify_on_exit | bool | 否 | 离开通知 (默认true) |

---

### 5.2 获取地理围栏列表

**GET** `/customer/v1/geofence/list`

**响应**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "geofences": [
      {
        "geofence_id": 10001,
        "name": "家",
        "longitude": 116.397428,
        "latitude": 39.90923,
        "radius": 100.0,
        "notify_on_enter": true,
        "notify_on_exit": true,
        "is_active": true,
        "created_at": "2024-01-01T00:00:00Z"
      }
    ]
  }
}
```

---

### 5.3 更新地理围栏

**PUT** `/customer/v1/geofence/:geofence_id`

**请求参数**
```json
{
  "name": "新名称",
  "radius": 200.0,
  "notify_on_enter": true,
  "notify_on_exit": false
}
```

---

### 5.4 删除地理围栏

**DELETE** `/customer/v1/geofence/:geofence_id`

---

## 六、WebSocket

### 6.1 建立连接

**GET** `/customer/v1/ws`

**请求参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| token | string | 是 | JWT Token |

**连接示例**
```
ws://localhost:8900/api/app/customer/v1/ws?token=eyJhbGciOiJIUzI1NiIs...
```

### 6.2 消息格式

**客户端发送 - 订阅位置更新**
```json
{
  "action": "subscribe",
  "entities": [
    {"type": "user", "id": "10002"},
    {"type": "device", "id": "ABC123"}
  ]
}
```

**服务端推送 - 位置更新**
```json
{
  "action": "location_update",
  "data": {
    "entity_type": "user",
    "entity_id": "10002",
    "longitude": 116.397428,
    "latitude": 39.90923,
    "timestamp": 1704067200000
  }
}
```

### 6.3 支持的消息类型

| 动作 | 方向 | 说明 |
|------|------|------|
| subscribe | 客户端->服务端 | 订阅位置更新 |
| unsubscribe | 客户端->服务端 | 取消订阅 |
| location_update | 服务端->客户端 | 位置更新推送 |
| geofence_event | 服务端->客户端 | 地理围栏事件推送 |

---

## 七、管理后台

### 7.1 管理员登录

**POST** `/admin/v1/user/login`

**请求参数**
```json
{
  "username": "admin",
  "password": "admin123",
  "captcha_key": "xxx",
  "captcha_code": "abc123"
}
```

**响应**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 86400
  }
}
```

---

### 7.2 获取滑块验证码

**GET** `/admin/v1/user/verify/captcha`

**响应**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "captcha_key": "xxx",
    "captcha_image": "data:image/png;base64,..."
  }
}
```

---

### 7.3 校验滑块验证码

**POST** `/admin/v1/user/verify/captcha/check`

**请求参数**
```json
{
  "captcha_key": "xxx",
  "captcha_code": "abc123"
}
```

---

### 7.4 获取管理员信息

**GET** `/admin/v1/user/info`

**请求头**
```
Authorization: Bearer <token>
```

---

### 7.5 创建管理员用户

**POST** `/admin/v1/user/create`

**请求参数**
```json
{
  "username": "admin2",
  "password": "password123",
  "nickname": "管理员2",
  "role_ids": [1, 2]
}
```

---

### 7.6 更新管理员用户

**POST** `/admin/v1/user/update`

**请求参数**
```json
{
  "user_id": 2,
  "nickname": "新昵称",
  "role_ids": [1]
}
```
