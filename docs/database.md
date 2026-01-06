# 数据库设计文档

## 概述

本文档描述 FindLink 位置共享服务的数据库设计，包括表结构、索引设计和 ER 图。

## 数据库信息

- **数据库**: MySQL 8.0+
- **字符集**: utf8mb4
- **排序规则**: utf8mb4_unicode_ci

## 数据库表列表

| 表名 | 说明 |
|------|------|
| admin_user | 管理员用户表 |
| admin_user_role | 管理员用户角色关联表 |
| role | 角色表 |
| permission | 权限表 |
| role_permission | 角色权限关联表 |
| user | C端用户表 |
| user_settings | 用户设置表 |
| user_location | 用户位置历史表 |
| device | 设备表 |
| device_location | 设备位置历史表 |
| friend | 好友关系表 |
| friend_request | 好友请求表 |
| geofence | 地理围栏表 |
| geofence_event | 地理围栏事件表 |

---

## 一、管理后台相关表

### 1.1 admin_user - 管理员用户表

```sql
CREATE TABLE `admin_user` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `username` varchar(64) NOT NULL COMMENT '用户名',
  `password` varchar(256) NOT NULL COMMENT '密码(SHA256)',
  `nickname` varchar(64) NOT NULL DEFAULT '' COMMENT '昵称',
  `avatar` varchar(256) DEFAULT '' COMMENT '头像URL',
  `phone` varchar(20) DEFAULT '' COMMENT '手机号',
  `email` varchar(128) DEFAULT '' COMMENT '邮箱',
  `sex` tinyint DEFAULT '0' COMMENT '性别: 0-未知 1-男 2-女',
  `status` tinyint DEFAULT '1' COMMENT '状态: 0-禁用 1-启用',
  `login_at` datetime DEFAULT NULL COMMENT '最后登录时间',
  `login_ip` varchar(128) DEFAULT '' COMMENT '最后登录IP',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`),
  KEY `idx_phone` (`phone`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='管理员用户表';
```

### 1.2 role - 角色表

```sql
CREATE TABLE `role` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `name` varchar(64) NOT NULL COMMENT '角色名称',
  `code` varchar(64) NOT NULL COMMENT '角色编码',
  `description` varchar(256) DEFAULT '' COMMENT '描述',
  `status` tinyint DEFAULT '1' COMMENT '状态: 0-禁用 1-启用',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色表';
```

### 1.3 permission - 权限表

```sql
CREATE TABLE `permission` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `name` varchar(64) NOT NULL COMMENT '权限名称',
  `code` varchar(128) NOT NULL COMMENT '权限编码',
  `type` varchar(32) DEFAULT '' COMMENT '类型: menu-菜单 button-按钮',
  `parent_id` bigint DEFAULT '0' COMMENT '父级ID',
  `path` varchar(256) DEFAULT '' COMMENT '路由路径',
  `icon` varchar(128) DEFAULT '' COMMENT '图标',
  `sort` int DEFAULT '0' COMMENT '排序',
  `status` tinyint DEFAULT '1' COMMENT '状态: 0-禁用 1-启用',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_parent_id` (`parent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='权限表';
```

### 1.4 admin_user_role - 管理员用户角色关联表

```sql
CREATE TABLE `admin_user_role` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `role_id` bigint NOT NULL COMMENT '角色ID',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_role` (`user_id`, `role_id`),
  KEY `idx_role_id` (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='管理员用户角色关联表';
```

### 1.5 role_permission - 角色权限关联表

```sql
CREATE TABLE `role_permission` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `role_id` bigint NOT NULL COMMENT '角色ID',
  `permission_id` bigint NOT NULL COMMENT '权限ID',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_role_permission` (`role_id`, `permission_id`),
  KEY `idx_permission_id` (`permission_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色权限关联表';
```

---

## 二、C端用户相关表

### 2.1 user - 用户表

```sql
CREATE TABLE `user` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `phone` varchar(20) NOT NULL COMMENT '手机号',
  `password` varchar(256) NOT NULL COMMENT '密码(SHA256)',
  `nickname` varchar(64) NOT NULL DEFAULT '' COMMENT '昵称',
  `avatar` varchar(256) DEFAULT '' COMMENT '头像URL',
  `sex` tinyint DEFAULT '0' COMMENT '性别: 0-未知 1-男 2-女',
  `status` tinyint DEFAULT '1' COMMENT '状态: 0-禁用 1-启用',
  `last_login_at` datetime DEFAULT NULL COMMENT '最后登录时间',
  `last_login_ip` varchar(128) DEFAULT '' COMMENT '最后登录IP',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_phone` (`phone`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';
```

### 2.2 user_settings - 用户设置表

```sql
CREATE TABLE `user_settings` (
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `share_location` tinyint DEFAULT '1' COMMENT '是否分享位置: 0-否 1-是',
  `ghost_mode` tinyint DEFAULT '0' COMMENT '幽灵模式: 0-关闭 1-开启',
  `smart_alerts` tinyint DEFAULT '1' COMMENT '智能提醒: 0-关闭 1-开启',
  `sos_alerts` tinyint DEFAULT '1' COMMENT 'SOS提醒: 0-关闭 1-开启',
  `map_style` varchar(20) DEFAULT 'dark' COMMENT '地图风格: dark-暗色 light-亮色',
  `distance_unit` varchar(10) DEFAULT 'km' COMMENT '距离单位: km-公里 mi-英里',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户设置表';
```

---

## 三、位置相关表

### 3.1 user_location - 用户位置历史表

```sql
CREATE TABLE `user_location` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `longitude` decimal(10,7) NOT NULL COMMENT '经度',
  `latitude` decimal(10,7) NOT NULL COMMENT '纬度',
  `accuracy` float DEFAULT NULL COMMENT '精度(米)',
  `altitude` float DEFAULT NULL COMMENT '海拔(米)',
  `speed` float DEFAULT NULL COMMENT '速度(米/秒)',
  `bearing` float DEFAULT NULL COMMENT '方向(0-360)',
  `battery_level` tinyint DEFAULT NULL COMMENT '电量(0-100)',
  `location_mode` varchar(20) DEFAULT 'foreground' COMMENT '定位模式: foreground-前台 background-后台 significant_change-显著变化',
  `is_low_accuracy` tinyint DEFAULT '0' COMMENT '是否低精度: 0-否 1-是',
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_user_created` (`user_id`, `created_at`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户位置历史表';
```

### 3.2 device_location - 设备位置历史表

```sql
CREATE TABLE `device_location` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `device_id` varchar(64) NOT NULL COMMENT '设备ID',
  `longitude` decimal(10,7) NOT NULL COMMENT '经度',
  `latitude` decimal(10,7) NOT NULL COMMENT '纬度',
  `accuracy` float DEFAULT NULL COMMENT '精度(米)',
  `battery_level` tinyint DEFAULT NULL COMMENT '电量(0-100)',
  `connection_status` varchar(20) DEFAULT 'unknown' COMMENT '连接状态: online-在线 offline-离线 unknown-未知',
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_device_id` (`device_id`),
  KEY `idx_device_created` (`device_id`, `created_at`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='设备位置历史表';
```

---

## 四、设备相关表

### 4.1 device - 设备表

```sql
CREATE TABLE `device` (
  `id` varchar(64) NOT NULL COMMENT '设备ID',
  `user_id` bigint NOT NULL COMMENT '绑定用户ID',
  `name` varchar(100) DEFAULT '' COMMENT '设备名称',
  `type` varchar(50) DEFAULT 'other' COMMENT '设备类型: gps_tracker-追踪器 smart_watch-智能手表 other-其他',
  `battery_level` tinyint DEFAULT NULL COMMENT '电量(0-100)',
  `connection_status` varchar(20) DEFAULT 'unknown' COMMENT '连接状态: online-在线 offline-离线 unknown-未知',
  `notification_enabled` tinyint DEFAULT '1' COMMENT '通知开关: 0-关 1-开',
  `last_seen_at` datetime DEFAULT NULL COMMENT '最后在线时间',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='设备表';
```

---

## 五、好友相关表

### 5.1 friend - 好友关系表

```sql
CREATE TABLE `friend` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `friend_id` bigint NOT NULL COMMENT '好友ID',
  `status` varchar(20) NOT NULL DEFAULT 'pending' COMMENT '状态: pending-待接受 accepted-已接受 rejected-已拒绝',
  `sharing_status` varchar(20) DEFAULT 'sharing' COMMENT '共享状态: sharing-共享 paused-暂停 hidden-隐藏',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_friend` (`user_id`, `friend_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_friend_id` (`friend_id`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='好友关系表';
```

### 5.2 friend_request - 好友请求表

```sql
CREATE TABLE `friend_request` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `from_user_id` bigint NOT NULL COMMENT '发起用户ID',
  `to_user_id` bigint NOT NULL COMMENT '目标用户ID',
  `message` varchar(256) DEFAULT '' COMMENT '验证消息',
  `status` varchar(20) NOT NULL DEFAULT 'pending' COMMENT '状态: pending-待处理 accepted-已接受 rejected-已拒绝',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_from_to` (`from_user_id`, `to_user_id`),
  KEY `idx_from_user` (`from_user_id`),
  KEY `idx_to_user` (`to_user_id`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='好友请求表';
```

---

## 六、地理围栏相关表

### 6.1 geofence - 地理围栏表

```sql
CREATE TABLE `geofence` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `name` varchar(100) NOT NULL COMMENT '围栏名称',
  `longitude` decimal(10,7) NOT NULL COMMENT '中心点经度',
  `latitude` decimal(10,7) NOT NULL COMMENT '中心点纬度',
  `radius` float NOT NULL COMMENT '半径(米)',
  `notify_on_enter` tinyint DEFAULT '1' COMMENT '进入通知: 0-否 1-是',
  `notify_on_exit` tinyint DEFAULT '1' COMMENT '离开通知: 0-否 1-是',
  `is_active` tinyint DEFAULT '1' COMMENT '是否激活: 0-否 1-是',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='地理围栏表';
```

### 6.2 geofence_event - 地理围栏事件表

```sql
CREATE TABLE `geofence_event` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `geofence_id` bigint NOT NULL COMMENT '围栏ID',
  `entity_type` varchar(20) NOT NULL COMMENT '实体类型: user-用户 device-设备',
  `entity_id` varchar(64) NOT NULL COMMENT '实体ID',
  `event_type` varchar(20) NOT NULL COMMENT '事件类型: enter-进入 exit-离开',
  `longitude` decimal(10,7) NOT NULL COMMENT '触发时经度',
  `latitude` decimal(10,7) NOT NULL COMMENT '触发时纬度',
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_geofence_id` (`geofence_id`),
  KEY `idx_entity` (`entity_type`, `entity_id`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='地理围栏事件表';
```

---

## 七、ER 图

```
┌─────────────────┐       ┌─────────────────┐
│   admin_user    │       │      role       │
├─────────────────┤       ├─────────────────┤
│ id (PK)         │──────<│ id (PK)         │
│ username        │       │ name            │
│ password        │       │ code            │
│ nickname        │       └────────┬────────┘
│ ...             │                │
└────────┬────────┘                │
         │                         │
         │    ┌────────────────────┘
         │    │
         v    v
┌─────────────────────┐
│  admin_user_role    │
├─────────────────────┤
│ id (PK)             │
│ user_id (FK)        │
│ role_id (FK)        │
└─────────────────────┘
         │
         │         ┌─────────────────┐
         └───────< │    permission   │
                   ├─────────────────┤
                   │ id (PK)         │
                   │ name            │
                   │ code            │
                   │ parent_id       │
                   └─────────────────┘
                           │
                           v
                   ┌─────────────────────┐
                   │  role_permission    │
                   ├─────────────────────┤
                   │ id (PK)             │
                   │ role_id (FK)        │
                   │ permission_id (FK)  │
                   └─────────────────────┘


┌─────────────────┐       ┌─────────────────────┐
│      user       │       │    user_settings    │
├─────────────────┤       ├─────────────────────┤
│ id (PK)         │──┐    │ user_id (PK,FK)     │
│ phone           │  │    │ share_location      │
│ password        │  │    │ ghost_mode          │
│ nickname        │  │    │ ...                 │
│ ...             │  │    └─────────────────────┘
└─────────────────┘  │
                    v
          ┌─────────────────────┐
          │   user_location     │
          ├─────────────────────┤
          │ id (PK)             │
          │ user_id (FK)        │
          │ longitude           │
          │ latitude            │
          │ ...                 │
          └─────────────────────┘


┌─────────────────┐       ┌─────────────────────┐
│     device      │       │   device_location   │
├─────────────────┤       ├─────────────────────┤
│ id (PK)         │──┐    │ id (PK)             │
│ user_id (FK)    │  │    │ device_id (FK)      │
│ name            │  │    │ longitude           │
│ type            │  │    │ latitude            │
│ ...             │  │    │ ...                 │
└─────────────────┘  │    └─────────────────────┘
                     │
                     v
          ┌─────────────────────┐
          │   geofence          │
          ├─────────────────────┤
          │ id (PK)             │
          │ user_id (FK)        │
          │ name                │
          │ longitude           │
          │ latitude            │
          │ radius              │
          └──────────┬──────────┘
                     │
                     v
          ┌─────────────────────┐
          │   geofence_event    │
          ├─────────────────────┤
          │ id (PK)             │
          │ geofence_id (FK)    │
          │ entity_type         │
          │ entity_id           │
          │ event_type          │
          └─────────────────────┘


┌─────────────────┐       ┌─────────────────────┐
│      friend     │       │   friend_request    │
├─────────────────┤       ├─────────────────────┤
│ id (PK)         │       │ id (PK)             │
│ user_id (FK)    │       │ from_user_id (FK)   │
│ friend_id (FK)  │       │ to_user_id (FK)     │
│ status          │       │ message             │
│ sharing_status  │       │ status              │
└─────────────────┘       └─────────────────────┘
```

---

## 八、索引优化建议

### 8.1 常用查询优化

```sql
-- 用户位置查询优化
-- 按用户ID和时间范围查询
CREATE INDEX idx_user_location_query ON user_location (user_id, created_at DESC);

-- 设备位置查询优化
-- 按设备ID和时间范围查询
CREATE INDEX idx_device_location_query ON device_location (device_id, created_at DESC);

-- 好友查询优化
-- 按用户ID查询所有好友
CREATE INDEX idx_friend_user_list ON friend (user_id, status);

-- 地理围栏事件查询优化
-- 按实体查询最近事件
CREATE INDEX idx_geofence_event_entity ON geofence_event (entity_type, entity_id, created_at DESC);
```

### 8.2 分区建议

对于大表 `user_location` 和 `device_location`，建议按时间分区：

```sql
-- 按月分区示例
ALTER TABLE user_location PARTITION BY RANGE (TO_DAYS(created_at)) (
    PARTITION p202401 VALUES LESS THAN (TO_DAYS('2024-02-01')),
    PARTITION p202402 VALUES LESS THAN (TO_DAYS('2024-03-01')),
    PARTITION p202403 VALUES LESS THAN (TO_DAYS('2024-04-01')),
    ...
);
```
