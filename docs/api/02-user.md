# 用户模块 API

## 概述

用户模块提供用户信息获取、更新等接口。

**Base URL**: `http://localhost:8900/api/app/customer/v1`

**认证方式**: JWT Token (Header: `Authorization: Bearer <token>`)

---

## 一、获取用户信息

### 接口说明

获取当前登录用户的详细信息。

### 请求信息

| 项目 | 说明 |
|------|------|
| URL | `/user/info` |
| Method | `GET` |
| Header | `Authorization: Bearer <token>` |

### 响应参数

| 参数 | 类型 | 说明 |
|------|------|------|
| user_id | int64 | 用户ID |
| phone | string | 手机号 |
| nickname | string | 昵称 |
| avatar | string | 头像URL |
| sex | int | 性别：0-未知 1-男 2-女 |
| created_at | string | 创建时间 |

### 响应示例

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

## 二、更新用户信息

### 接口说明

更新当前登录用户的个人信息。

### 请求信息

| 项目 | 说明 |
|------|------|
| URL | `/user/update` |
| Method | `PUT` |
| Header | `Authorization: Bearer <token>` |
| Content-Type | `application/json` |

### 请求参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| nickname | string | 否 | 昵称（2-20字符） |
| avatar | string | 否 | 头像URL |
| sex | int | 否 | 性别：0-未知 1-男 2-女 |

### 请求示例

```json
{
  "nickname": "新昵称",
  "avatar": "https://example.com/new-avatar.png",
  "sex": 1
}
```

### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "user_id": 10001,
    "nickname": "新昵称",
    "avatar": "https://example.com/new-avatar.png",
    "sex": 1
  }
}
```

---

## 三、获取用户设置

### 接口说明

获取当前用户的隐私设置和偏好设置。

### 请求信息

| 项目 | 说明 |
|------|------|
| URL | `/user/settings` |
| Method | `GET` |
| Header | `Authorization: Bearer <token>` |

### 响应参数

| 参数 | 类型 | 说明 |
|------|------|------|
| share_location | bool | 是否分享位置：true-是 false-否 |
| ghost_mode | bool | 幽灵模式：true-开启 false-关闭 |
| smart_alerts | bool | 智能提醒：true-开启 false-关闭 |
| sos_alerts | bool | SOS提醒：true-开启 false-关闭 |
| map_style | string | 地图风格：dark-暗色 light-亮色 |
| distance_unit | string | 距离单位：km-公里 mi-英里 |

### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "share_location": true,
    "ghost_mode": false,
    "smart_alerts": true,
    "sos_alerts": true,
    "map_style": "dark",
    "distance_unit": "km"
  }
}
```

---

## 四、更新用户设置

### 接口说明

更新当前用户的隐私设置和偏好设置。

### 请求信息

| 项目 | 说明 |
|------|------|
| URL | `/user/settings` |
| Method | `PUT` |
| Header | `Authorization: Bearer <token>` |
| Content-Type | `application/json` |

### 请求参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| share_location | bool | 否 | 是否分享位置 |
| ghost_mode | bool | 否 | 幽灵模式 |
| smart_alerts | bool | 否 | 智能提醒 |
| sos_alerts | bool | 否 | SOS提醒 |
| map_style | string | 否 | 地图风格 |
| distance_unit | string | 否 | 距离单位 |

### 请求示例

```json
{
  "share_location": true,
  "ghost_mode": true,
  "smart_alerts": false,
  "sos_alerts": true,
  "map_style": "light",
  "distance_unit": "mi"
}
```

### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "share_location": true,
    "ghost_mode": true,
    "smart_alerts": false,
    "sos_alerts": true,
    "map_style": "light",
    "distance_unit": "mi"
  }
}
```

---

## 五、获取指定用户信息

### 接口说明

获取其他用户的公开信息（无需登录）。

### 请求信息

| 项目 | 说明 |
|------|------|
| URL | `/user/profile/:user_id` |
| Method | `GET` |

### 路径参数

| 参数 | 类型 | 说明 |
|------|------|------|
| user_id | int64 | 用户ID |

### 响应参数

| 参数 | 类型 | 说明 |
|------|------|------|
| user_id | int64 | 用户ID |
| nickname | string | 昵称 |
| avatar | string | 头像URL |
| is_online | bool | 是否在线 |

### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "user_id": 10002,
    "nickname": "好友昵称",
    "avatar": "https://example.com/avatar.png",
    "is_online": true
  }
}
```

---

## 六、修改密码

### 接口说明

修改当前用户的登录密码。

### 请求信息

| 项目 | 说明 |
|------|------|
| URL | `/user/password` |
| Method | `PUT` |
| Header | `Authorization: Bearer <token>` |
| Content-Type | `application/json` |

### 请求参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| old_password | string | 是 | 原密码 |
| new_password | string | 是 | 新密码（6-20位） |

### 请求示例

```json
{
  "old_password": "old_password123",
  "new_password": "new_password123"
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

## 七、注销账号

### 接口说明

注销当前用户账号（软删除）。

### 请求信息

| 项目 | 说明 |
|------|------|
| URL | `/user/account` |
| Method | `DELETE` |
| Header | `Authorization: Bearer <token>` |

### 请求参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| password | string | 是 | 当前密码 |
| verify_code | string | 是 | 短信验证码 |

### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": null
}
```

### 注意事项

- 账号注销后不可恢复
- 注销后 30 天内可申请恢复
- 30 天后数据将被永久删除
