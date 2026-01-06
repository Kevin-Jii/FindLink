# 认证模块 API

## 概述

认证模块提供用户登录、注册、登出等认证相关接口。

**Base URL**: `http://localhost:8900/api/app/customer/v1`

**认证方式**: JWT Token (Header: `Authorization: Bearer <token>`)

---

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

---

## 一、用户登录

### 接口说明

用户使用手机号和密码登录，获取 JWT 令牌。

### 请求信息

| 项目 | 说明 |
|------|------|
| URL | `/user/login` |
| Method | `POST` |
| Content-Type | `application/json` |

### 请求参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| phone | string | 是 | 手机号，格式：`1[3-9]\d{9}` |
| password | string | 是 | 密码，最小6位，最大20位 |

### 请求示例

```json
{
  "phone": "13800138000",
  "password": "password123"
}
```

### 响应参数

| 参数 | 类型 | 说明 |
|------|------|------|
| token | string | JWT 访问令牌 |
| user_id | int64 | 用户ID |
| expires_in | int | 令牌有效期（秒） |

### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user_id": 10001,
    "expires_in": 86400
  }
}
```

### 错误码

| 错误码 | 说明 |
|--------|------|
| 10001 | 参数错误 |
| 10002 | 用户名或密码错误 |
| 10004 | 用户不存在 |
| 10005 | 用户已被禁用 |

---

## 二、用户注册

### 接口说明

新用户使用手机号注册账号，需要短信验证码验证。

### 请求信息

| 项目 | 说明 |
|------|------|
| URL | `/user/register` |
| Method | `POST` |
| Content-Type | `application/json` |

### 请求参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| phone | string | 是 | 手机号 |
| password | string | 是 | 密码（6-20位） |
| nickname | string | 是 | 昵称（2-20字符） |
| verify_code | string | 是 | 短信验证码（6位数字） |

### 请求示例

```json
{
  "phone": "13800138000",
  "password": "password123",
  "nickname": "新用户",
  "verify_code": "123456"
}
```

### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user_id": 10001,
    "expires_in": 86400
  }
}
```

### 错误码

| 错误码 | 说明 |
|--------|------|
| 10001 | 参数错误 |
| 10006 | 验证码错误或已过期 |
| 10007 | 手机号已被注册 |

---

## 三、用户登出

### 接口说明

用户登出，使当前 token 失效。

### 请求信息

| 项目 | 说明 |
|------|------|
| URL | `/user/logout` |
| Method | `POST` |
| Header | `Authorization: Bearer <token>` |

### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": null
}
```

---

## 四、发送短信验证码

### 接口说明

发送手机短信验证码，用于注册和找回密码。

### 请求信息

| 项目 | 说明 |
|------|------|
| URL | `/user/verify/sms` |
| Method | `POST` |
| Content-Type | `application/json` |

### 请求参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| phone | string | 是 | 手机号 |
| type | string | 是 | 验证码类型：`register` / `reset_password` |

### 请求示例

```json
{
  "phone": "13800138000",
  "type": "register"
}
```

### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "expires_in": 300
  }
}
```

### 验证码有效期

- 有效期：5分钟
- 发送间隔：60秒
- 每个手机号每天最多发送：10次

---

## 五、刷新令牌

### 接口说明

使用刷新令牌获取新的访问令牌。

### 请求信息

| 项目 | 说明 |
|------|------|
| URL | `/user/refresh` |
| Method | `POST` |
| Content-Type | `application/json` |

### 请求参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| refresh_token | string | 是 | 刷新令牌 |

### 响应参数

| 参数 | 类型 | 说明 |
|------|------|------|
| token | string | 新的访问令牌 |
| refresh_token | string | 新的刷新令牌 |
| expires_in | int | 有效期（秒） |

---

## 六、忘记密码

### 接口说明

用户忘记密码，通过短信验证码重置密码。

### 请求信息

| 项目 | 说明 |
|------|------|
| URL | `/user/password/reset` |
| Method | `POST` |
| Content-Type | `application/json` |

### 请求参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| phone | string | 是 | 手机号 |
| verify_code | string | 是 | 短信验证码 |
| new_password | string | 是 | 新密码（6-20位） |

### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": null
}
```
