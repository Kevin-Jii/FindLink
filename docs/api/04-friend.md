# 好友模块 API

## 概述

好友模块提供好友关系管理、好友请求处理等功能。

**Base URL**: `http://localhost:8900/api/app/customer/v1/friend`

**认证方式**: JWT Token (Header: `Authorization: Bearer <token>`)

---

## 一、发送好友请求

### 接口说明

向指定用户发送好友请求。

### 请求信息

| 项目 | 说明 |
|------|------|
| URL | `/request` |
| Method | `POST` |
| Header | `Authorization: Bearer <token>` |
| Content-Type | `application/json` |

### 请求参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| target_user_id | int64 | 是 | 目标用户ID |
| message | string | 否 | 验证消息（最多256字符） |

### 请求示例

```json
{
  "target_user_id": 10002,
  "message": "我是xxx，想加你好友"
}
```

### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "request_id": 10001
  }
}
```

### 错误码

| 错误码 | 说明 |
|--------|------|
| 13004 | 不能添加自己为好友 |
| 13003 | 已经是好友关系 |
| 13002 | 已存在待处理的好友请求 |

---

## 二、获取好友请求列表

### 接口说明

获取收到或发出的好友请求列表。

### 请求信息

| 项目 | 说明 |
|------|------|
| URL | `/requests` |
| Method | `GET` |
| Header | `Authorization: Bearer <token>` |

### 查询参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| type | string | 否 | 请求类型：received-收到的 / sent-发出的（默认received） |
| status | string | 否 | 状态：pending-待处理 / accepted-已接受 / rejected-已拒绝 |
| page | int | 否 | 页码（默认1） |
| page_size | int | 否 | 每页数量（默认20，最大100） |

### 响应参数

| 参数 | 类型 | 说明 |
|------|------|------|
| requests | array | 请求列表 |
| total | int | 总数量 |
| page | int | 当前页码 |
| page_size | int | 每页数量 |

#### Request 对象（收到的请求）

| 参数 | 类型 | 说明 |
|------|------|------|
| request_id | int64 | 请求ID |
| from_user_id | int64 | 发起用户ID |
| from_user_nickname | string | 发起用户昵称 |
| from_user_avatar | string | 发起用户头像 |
| message | string | 验证消息 |
| status | string | 状态 |
| created_at | string | 发起时间 |

#### Request 对象（发出的请求）

| 参数 | 类型 | 说明 |
|------|------|------|
| request_id | int64 | 请求ID |
| to_user_id | int64 | 目标用户ID |
| to_user_nickname | string | 目标用户昵称 |
| to_user_avatar | string | 目标用户头像 |
| message | string | 验证消息 |
| status | string | 状态 |
| created_at | string | 发起时间 |

### 响应示例

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
    ],
    "total": 10,
    "page": 1,
    "page_size": 20
  }
}
```

---

## 三、接受好友请求

### 接口说明

接受指定的好友请求，建立双向好友关系。

### 请求信息

| 项目 | 说明 |
|------|------|
| URL | `/accept` |
| Method | `POST` |
| Header | `Authorization: Bearer <token>` |
| Content-Type | `application/json` |

### 请求参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| request_id | int64 | 是 | 好友请求ID |

### 请求示例

```json
{
  "request_id": 10001
}
```

### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "friend_id": 10002,
    "nickname": "用户A",
    "avatar": "https://example.com/avatar.png"
  }
}
```

---

## 四、拒绝好友请求

### 接口说明

拒绝指定的好友请求。

### 请求信息

| 项目 | 说明 |
|------|------|
| URL | `/reject` |
| Method | `POST` |
| Header | `Authorization: Bearer <token>` |
| Content-Type | `application/json` |

### 请求参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| request_id | int64 | 是 | 好友请求ID |

### 请求示例

```json
{
  "request_id": 10001
}
```

### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": null
}
```

---

## 五、获取好友列表

### 接口说明

获取当前用户的好友列表。

### 请求信息

| 项目 | 说明 |
|------|------|
| URL | `/list` |
| Method | `GET` |
| Header | `Authorization: Bearer <token>` |

### 查询参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| status | string | 否 | 共享状态：sharing-共享中 / paused-已暂停 / hidden-隐藏 |
| page | int | 否 | 页码（默认1） |
| page_size | int | 否 | 每页数量（默认20，最大100） |

### 响应参数

| 参数 | 类型 | 说明 |
|------|------|------|
| friends | array | 好友列表 |
| total | int | 总数量 |

#### Friend 对象

| 参数 | 类型 | 说明 |
|------|------|------|
| friend_id | int64 | 好友用户ID |
| nickname | string | 好友昵称 |
| avatar | string | 好友头像 |
| sharing_status | string | 共享状态 |
| last_location | object | 最后位置 |
| is_online | bool | 是否在线 |
| created_at | string | 成为好友时间 |

#### last_location 对象

| 参数 | 类型 | 说明 |
|------|------|------|
| longitude | float | 经度 |
| latitude | float | 纬度 |
| updated_at | string | 更新时间 |

### 响应示例

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
        },
        "is_online": true,
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 10
  }
}
```

---

## 六、删除好友

### 接口说明

删除指定好友，同时解除双向好友关系。

### 请求信息

| 项目 | 说明 |
|------|------|
| URL | `/:friend_id` |
| Method | `DELETE` |
| Header | `Authorization: Bearer <token>` |

### 路径参数

| 参数 | 类型 | 说明 |
|------|------|------|
| friend_id | int64 | 好友用户ID |

### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": null
}
```

---

## 七、搜索用户

### 接口说明

通过手机号或昵称搜索用户。

### 请求信息

| 项目 | 说明 |
|------|------|
| URL | `/search` |
| Method | `GET` |
| Header | `Authorization: Bearer <token>` |

### 查询参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| keyword | string | 是 | 搜索关键词（手机号或昵称） |

### 响应参数

| 参数 | 类型 | 说明 |
|------|------|------|
| users | array | 用户列表 |

#### User 对象

| 参数 | 类型 | 说明 |
|------|------|------|
| user_id | int64 | 用户ID |
| phone | string | 手机号（脱敏） |
| nickname | string | 昵称 |
| avatar | string | 头像URL |
| is_friend | bool | 是否已是好友 |
| is_pending | bool | 是否有待处理请求 |

### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "users": [
      {
        "user_id": 10002,
        "phone": "138****0000",
        "nickname": "用户A",
        "avatar": "https://example.com/avatar.png",
        "is_friend": false,
        "is_pending": false
      }
    ]
  }
}
```

---

## 八、取消好友请求

### 接口说明

取消自己发出的好友请求。

### 请求信息

| 项目 | 说明 |
|------|------|
| URL | `/request/:request_id` |
| Method | `DELETE` |
| Header | `Authorization: Bearer <token>` |

### 路径参数

| 参数 | 类型 | 说明 |
|------|------|------|
| request_id | int64 | 请求ID |

### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": null
}
```

---

## 九、更新好友共享状态

### 接口说明

更新与指定好友的位置共享状态。

### 请求信息

| 项目 | 说明 |
|------|------|
| URL | `/:friend_id/status` |
| Method | `PUT` |
| Header | `Authorization: Bearer <token>` |
| Content-Type | `application/json` |

### 路径参数

| 参数 | 类型 | 说明 |
|------|------|------|
| friend_id | int64 | 好友用户ID |

### 请求参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| sharing_status | string | 是 | 共享状态：sharing-共享 / paused-暂停 / hidden-隐藏 |

### 请求示例

```json
{
  "sharing_status": "paused"
}
```

### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "friend_id": 10002,
    "sharing_status": "paused"
  }
}
```

---

## 状态说明

### 好友请求状态

| 状态 | 说明 |
|------|------|
| pending | 待处理 |
| accepted | 已接受 |
| rejected | 已拒绝 |

### 共享状态

| 状态 | 说明 |
|------|------|
| sharing | 共享中 |
| paused | 已暂停 |
| hidden | 隐藏 |
