# 设备模块 API

## 概述

设备模块提供设备绑定、解绑、状态管理等接口。

**Base URL**: `http://localhost:8900/api/app/customer/v1/device`

**认证方式**: JWT Token (Header: `Authorization: Bearer <token>`)

---

## 一、绑定设备

### 接口说明

将设备绑定到当前用户账号。

### 请求信息

| 项目 | 说明 |
|------|------|
| URL | `/bind` |
| Method | `POST` |
| Header | `Authorization: Bearer <token>` |
| Content-Type | `application/json` |

### 请求参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| device_id | string | 是 | 设备ID（不超过64字符） |
| name | string | 否 | 设备名称（不超过100字符） |
| type | string | 否 | 设备类型 |

#### type 取值

| 值 | 说明 |
|----|------|
| gps_tracker | GPS追踪器 |
| smart_watch | 智能手表 |
| other | 其他 |

### 请求示例

```json
{
  "device_id": "ABC123",
  "name": "我的设备",
  "type": "gps_tracker"
}
```

### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "device_id": "ABC123",
    "name": "我的设备",
    "type": "gps_tracker"
  }
}
```

### 错误码

| 错误码 | 说明 |
|--------|------|
| 14002 | 设备已被其他用户绑定 |
| 14004 | 设备ID格式无效 |

---

## 二、获取设备列表

### 接口说明

获取当前用户绑定的所有设备列表。

### 请求信息

| 项目 | 说明 |
|------|------|
| URL | `/list` |
| Method | `GET` |
| Header | `Authorization: Bearer <token>` |

### 响应参数

| 参数 | 类型 | 说明 |
|------|------|------|
| devices | array | 设备列表 |

#### Device 对象

| 参数 | 类型 | 说明 |
|------|------|------|
| device_id | string | 设备ID |
| name | string | 设备名称 |
| type | string | 设备类型 |
| battery_level | int | 电量（0-100） |
| connection_status | string | 连接状态 |
| notification_enabled | bool | 通知是否开启 |
| last_location | object | 最后位置 |
| last_seen_at | string | 最后在线时间 |
| created_at | string | 绑定时间 |

#### last_location 对象

| 参数 | 类型 | 说明 |
|------|------|------|
| longitude | float | 经度 |
| latitude | float | 纬度 |
| accuracy | float | 精度（米） |
| updated_at | string | 更新时间 |

#### connection_status 取值

| 值 | 说明 |
|----|------|
| online | 在线 |
| offline | 离线 |
| unknown | 未知 |

### 响应示例

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
        "notification_enabled": true,
        "last_location": {
          "longitude": 116.397428,
          "latitude": 39.90923,
          "accuracy": 10.5,
          "updated_at": "2024-01-01T00:00:00Z"
        },
        "last_seen_at": "2024-01-01T00:00:00Z",
        "created_at": "2024-01-01T00:00:00Z"
      }
    ]
  }
}
```

---

## 三、获取设备详情

### 接口说明

获取指定设备的详细信息。

### 请求信息

| 项目 | 说明 |
|------|------|
| URL | `/:device_id` |
| Method | `GET` |
| Header | `Authorization: Bearer <token>` |

### 路径参数

| 参数 | 类型 | 说明 |
|------|------|------|
| device_id | string | 设备ID |

### 响应参数

同"获取设备列表"中的 Device 对象

### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "device_id": "ABC123",
    "name": "我的设备",
    "type": "gps_tracker",
    "battery_level": 85,
    "connection_status": "online",
    "notification_enabled": true,
    "last_location": {
      "longitude": 116.397428,
      "latitude": 39.90923,
      "accuracy": 10.5,
      "updated_at": "2024-01-01T00:00:00Z"
    },
    "last_seen_at": "2024-01-01T00:00:00Z",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

---

## 四、更新设备设置

### 接口说明

更新设备的名称和通知设置。

### 请求信息

| 项目 | 说明 |
|------|------|
| URL | `/:device_id/settings` |
| Method | `PUT` |
| Header | `Authorization: Bearer <token>` |
| Content-Type | `application/json` |

### 路径参数

| 参数 | 类型 | 说明 |
|------|------|------|
| device_id | string | 设备ID |

### 请求参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 否 | 设备名称（不超过100字符） |
| notification_enabled | bool | 否 | 通知是否开启 |

### 请求示例

```json
{
  "name": "新名称",
  "notification_enabled": true
}
```

### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "device_id": "ABC123",
    "name": "新名称",
    "notification_enabled": true
  }
}
```

---

## 五、更新设备状态

### 接口说明

更新设备的电池电量和连接状态（通常由设备调用）。

### 请求信息

| 项目 | 说明 |
|------|------|
| URL | `/:device_id/status` |
| Method | `PUT` |
| Header | `Authorization: Bearer <token>` |
| Content-Type | `application/json` |

### 路径参数

| 参数 | 类型 | 说明 |
|------|------|------|
| device_id | string | 设备ID |

### 请求参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| battery_level | int | 否 | 电量（0-100） |
| connection_status | string | 否 | 连接状态 |

### 请求示例

```json
{
  "battery_level": 75,
  "connection_status": "online"
}
```

### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "device_id": "ABC123",
    "battery_level": 75,
    "connection_status": "online"
  }
}
```

---

## 六、解绑设备

### 接口说明

解除设备与当前用户账号的绑定关系。

### 请求信息

| 项目 | 说明 |
|------|------|
| URL | `/:device_id` |
| Method | `DELETE` |
| Header | `Authorization: Bearer <token>` |

### 路径参数

| 参数 | 类型 | 说明 |
|------|------|------|
| device_id | string | 设备ID |

### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": null
}
```

### 注意事项

- 解绑后设备将无法继续上报位置
- 设备的地理位置历史记录将被保留
- 设备可以重新绑定到其他用户账号

---

## 七、设备上报位置

### 接口说明

设备主动上报当前位置（设备鉴权）。

### 请求信息

| 项目 | 说明 |
|------|------|
| URL | `/:device_id/location` |
| Method | `POST` |
| Content-Type | `application/json` |

### 请求头

| 参数 | 说明 |
|------|------|
| X-Device-Token | 设备鉴权令牌 |

### 路径参数

| 参数 | 类型 | 说明 |
|------|------|------|
| device_id | string | 设备ID |

### 请求参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| longitude | float | 是 | 经度 |
| latitude | float | 是 | 纬度 |
| accuracy | float | 否 | 精度（米） |
| battery_level | int | 否 | 电量（0-100） |

### 请求示例

```json
{
  "longitude": 116.397428,
  "latitude": 39.90923,
  "accuracy": 10.5,
  "battery_level": 85
}
```

### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "location_id": 10001
  }
}
```

---

## 八、获取设备位置历史

### 接口说明

获取指定设备的历史位置轨迹。

### 请求信息

| 项目 | 说明 |
|------|------|
| URL | `/:device_id/history` |
| Method | `GET` |
| Header | `Authorization: Bearer <token>` |

### 路径参数

| 参数 | 类型 | 说明 |
|------|------|------|
| device_id | string | 设备ID |

### 查询参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| start_time | int64 | 是 | 开始时间（毫秒时间戳） |
| end_time | int64 | 是 | 结束时间（毫秒时间戳） |
| limit | int | 否 | 返回数量限制（默认100，最大500） |

### 响应参数

| 参数 | 类型 | 说明 |
|------|------|------|
| locations | array | 位置历史列表 |
| total | int | 总数量 |

### 响应示例

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
        "battery_level": 85,
        "timestamp": 1704067200000
      }
    ],
    "total": 100
  }
}
```

---

## 九、批量更新设备状态

### 接口说明

批量更新多个设备的状态信息。

### 请求信息

| 项目 | 说明 |
|------|------|
| URL | `/batch/status` |
| Method | `PUT` |
| Header | `Authorization: Bearer <token>` |
| Content-Type | `application/json` |

### 请求参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| devices | array | 是 | 设备状态列表 |

#### DeviceStatus 对象

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| device_id | string | 是 | 设备ID |
| battery_level | int | 否 | 电量（0-100） |
| connection_status | string | 否 | 连接状态 |

### 请求示例

```json
{
  "devices": [
    {
      "device_id": "ABC123",
      "battery_level": 75,
      "connection_status": "online"
    },
    {
      "device_id": "DEF456",
      "battery_level": 60,
      "connection_status": "offline"
    }
  ]
}
```

### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "updated": 2
  }
}
```
