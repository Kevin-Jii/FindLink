# 位置模块 API

## 概述

位置模块提供位置上报、查询、历史轨迹、附近好友发现等功能。

**Base URL**: `http://localhost:8900/api/app/customer/v1/location`

**认证方式**: JWT Token (Header: `Authorization: Bearer <token>`)

---

## 一、上报位置

### 接口说明

用户或设备上报当前位置信息。

### 请求信息

| 项目 | 说明 |
|------|------|
| URL | `/report` |
| Method | `POST` |
| Header | `Authorization: Bearer <token>` |
| Content-Type | `application/json` |

### 请求参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| longitude | float | 是 | 经度，范围：-180 ~ 180 |
| latitude | float | 是 | 纬度，范围：-90 ~ 90 |
| accuracy | float | 否 | 定位精度（米），范围：0 ~ 99999 |
| altitude | float | 否 | 海拔高度（米） |
| speed | float | 否 | 速度（米/秒），范围：0 ~ 9999.99 |
| bearing | float | 否 | 方向（度），范围：0 ~ 360 |
| battery_level | int | 否 | 电量（0-100） |
| location_mode | string | 否 | 定位模式 |

#### location_mode 取值

| 值 | 说明 |
|----|------|
| foreground | 前台定位 |
| background | 后台定位 |
| significant_change | 显著位置变化 |

### 请求示例

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

## 二、批量上报位置

### 接口说明

批量上报多个位置信息，用于后台恢复或离线数据同步。

### 请求信息

| 项目 | 说明 |
|------|------|
| URL | `/batch` |
| Method | `POST` |
| Header | `Authorization: Bearer <token>` |
| Content-Type | `application/json` |

### 请求参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| locations | array | 是 | 位置数组 |

#### Location 对象

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| longitude | float | 是 | 经度 |
| latitude | float | 是 | 纬度 |
| accuracy | float | 否 | 精度（米） |
| altitude | float | 否 | 海拔（米） |
| speed | float | 否 | 速度（米/秒） |
| bearing | float | 否 | 方向（度） |
| timestamp | int64 | 是 | 时间戳（毫秒） |
| battery_level | int | 否 | 电量（0-100） |
| location_mode | string | 否 | 定位模式 |

### 请求示例

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

### 响应参数

| 参数 | 类型 | 说明 |
|------|------|------|
| total | int | 接收的位置数量 |
| saved | int | 成功保存的位置数量 |

### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 2,
    "saved": 2
  }
}
```

---

## 三、获取用户位置

### 接口说明

获取指定用户的最新位置信息。

### 请求信息

| 项目 | 说明 |
|------|------|
| URL | `/user/:user_id` |
| Method | `GET` |
| Header | `Authorization: Bearer <token>` |

### 路径参数

| 参数 | 类型 | 说明 |
|------|------|------|
| user_id | int64 | 用户ID |

### 响应参数

| 参数 | 类型 | 说明 |
|------|------|------|
| user_id | int64 | 用户ID |
| longitude | float | 经度 |
| latitude | float | 纬度 |
| accuracy | float | 定位精度（米） |
| battery_level | int | 电量（0-100） |
| location_mode | string | 定位模式 |
| updated_at | string | 更新时间 |
| is_online | bool | 是否在线 |

### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "user_id": 10002,
    "longitude": 116.397428,
    "latitude": 39.90923,
    "accuracy": 10.5,
    "battery_level": 85,
    "location_mode": "foreground",
    "updated_at": "2024-01-01T00:00:00Z",
    "is_online": true
  }
}
```

### 错误码

| 错误码 | 说明 |
|--------|------|
| 12001 | 位置不存在 |
| 12005 | 无权限查看该用户位置 |

---

## 四、获取设备位置

### 接口说明

获取指定设备的最新位置信息。

### 请求信息

| 项目 | 说明 |
|------|------|
| URL | `/device/:device_id` |
| Method | `GET` |
| Header | `Authorization: Bearer <token>` |

### 路径参数

| 参数 | 类型 | 说明 |
|------|------|------|
| device_id | string | 设备ID |

### 响应参数

| 参数 | 类型 | 说明 |
|------|------|------|
| device_id | string | 设备ID |
| longitude | float | 经度 |
| latitude | float | 纬度 |
| accuracy | float | 定位精度（米） |
| battery_level | int | 电量（0-100） |
| connection_status | string | 连接状态 |
| updated_at | string | 更新时间 |

### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "device_id": "ABC123",
    "longitude": 116.397428,
    "latitude": 39.90923,
    "accuracy": 10.5,
    "battery_level": 75,
    "connection_status": "online",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

---

## 五、获取位置历史

### 接口说明

获取指定时间段内的位置历史轨迹。

### 请求信息

| 项目 | 说明 |
|------|------|
| URL | `/history` |
| Method | `GET` |
| Header | `Authorization: Bearer <token>` |

### 查询参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| start_time | int64 | 是 | 开始时间（毫秒时间戳） |
| end_time | int64 | 是 | 结束时间（毫秒时间戳） |
| limit | int | 否 | 返回数量限制（默认100，最大500） |
| entity_type | string | 否 | 实体类型：user/device |
| entity_id | string | 否 | 实体ID（用户ID或设备ID） |

### 响应参数

| 参数 | 类型 | 说明 |
|------|------|------|
| locations | array | 位置数组 |
| total | int | 总数量 |

#### Location 对象

| 参数 | 类型 | 说明 |
|------|------|------|
| longitude | float | 经度 |
| latitude | float | 纬度 |
| accuracy | float | 定位精度（米） |
| altitude | float | 海拔（米） |
| speed | float | 速度（米/秒） |
| bearing | float | 方向（度） |
| timestamp | int64 | 时间戳（毫秒） |
| location_mode | string | 定位模式 |

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
        "altitude": 100.0,
        "speed": 0.0,
        "bearing": 0.0,
        "timestamp": 1704067200000,
        "location_mode": "foreground"
      }
    ],
    "total": 100
  }
}
```

---

## 六、获取附近好友

### 接口说明

获取附近的好友列表（按距离排序）。

### 请求信息

| 项目 | 说明 |
|------|------|
| URL | `/nearby` |
| Method | `GET` |
| Header | `Authorization: Bearer <token>` |

### 查询参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| radius | float | 否 | 搜索半径（米，默认1000，最大10000） |
| limit | int | 否 | 返回数量限制（默认20，最大100） |

### 响应参数

| 参数 | 类型 | 说明 |
|------|------|------|
| friends | array | 附近好友列表 |

#### Friend 对象

| 参数 | 类型 | 说明 |
|------|------|------|
| user_id | int64 | 用户ID |
| nickname | string | 昵称 |
| avatar | string | 头像URL |
| longitude | float | 经度 |
| latitude | float | 纬度 |
| distance | float | 距离（米） |
| updated_at | string | 位置更新时间 |
| is_online | bool | 是否在线 |

### 响应示例

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
        "updated_at": "2024-01-01T00:00:00Z",
        "is_online": true
      }
    ]
  }
}
```

---

## 七、删除位置历史

### 接口说明

删除指定时间段内的位置历史记录。

### 请求信息

| 项目 | 说明 |
|------|------|
| URL | `/history` |
| Method | `DELETE` |
| Header | `Authorization: Bearer <token>` |
| Content-Type | `application/json` |

### 请求参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| start_time | int64 | 是 | 开始时间（毫秒时间戳） |
| end_time | int64 | 是 | 结束时间（毫秒时间戳） |

### 请求示例

```json
{
  "start_time": 1704067200000,
  "end_time": 1704153600000
}
```

### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "deleted": 100
  }
}
```

---

## 八、获取我的最新位置

### 接口说明

获取当前登录用户的最新位置信息。

### 请求信息

| 项目 | 说明 |
|------|------|
| URL | `/me` |
| Method | `GET` |
| Header | `Authorization: Bearer <token>` |

### 响应参数

同"获取用户位置"接口

### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "user_id": 10001,
    "longitude": 116.397428,
    "latitude": 39.90923,
    "accuracy": 10.5,
    "battery_level": 85,
    "location_mode": "foreground",
    "updated_at": "2024-01-01T00:00:00Z",
    "is_online": true
  }
}
```
