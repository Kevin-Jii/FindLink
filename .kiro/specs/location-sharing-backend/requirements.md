# Requirements Document

## Introduction

本文档定义了实时位置共享应用后端 API 的需求规范。该系统支持用户之间的位置共享、设备管理、好友关系管理等核心功能。后端基于 Go + Gin + GORM + PostGIS 技术栈，需要支持高并发位置上报（写入）、WebSocket 实时推送、地理空间计算和 Redis 缓存加速。

系统需要处理移动端（iOS/Android）后台持续定位的特殊场景，包括低功耗模式、显著位置变化模式等。

## Glossary

- **User**: C端应用用户，可以共享位置、添加好友和管理设备
- **Friend**: 用户之间的好友关系，支持双向位置共享
- **Device**: 用户绑定的追踪设备（如GPS追踪器、智能手表等）
- **Location**: 地理位置信息，包含经纬度、时间戳、精度等
- **SharingStatus**: 位置共享状态（sharing/paused/hidden）
- **ConnectionStatus**: 设备连接状态（online/offline/unknown）
- **GeoFence**: 地理围栏，用于进出区域提醒
- **PostGIS**: PostgreSQL 的地理空间扩展，支持空间索引和地理计算
- **WebSocket**: 全双工通信协议，用于实时位置推送
- **BackgroundLocation**: 后台定位模式，支持 iOS/Android 后台持续上报
- **SignificantLocationChange**: iOS 显著位置变化模式，低功耗后台定位
- **Geohash**: 地理位置编码，用于空间索引和邻近查询

## Requirements

### Requirement 1: 用户认证与授权

**User Story:** As a user, I want to securely login and manage my account, so that I can access the location sharing features safely.

#### Acceptance Criteria

1. WHEN a user submits valid credentials THEN the System SHALL return a JWT token with user information
2. WHEN a user submits invalid credentials THEN the System SHALL return an authentication error without revealing which field is incorrect
3. WHEN a user registers with a new mobile number THEN the System SHALL create a new account and return a JWT token
4. WHEN a user requests their profile THEN the System SHALL return user information including avatar, nickname, and settings
5. WHEN a user updates their profile THEN the System SHALL persist the changes and return the updated profile
6. WHEN serializing user data to JSON THEN the System SHALL support round-trip serialization without data loss

### Requirement 2: 好友关系管理

**User Story:** As a user, I want to manage my friends list, so that I can share my location with trusted people.

#### Acceptance Criteria

1. WHEN a user sends a friend request THEN the System SHALL create a pending friend relationship and notify the target user
2. WHEN a user accepts a friend request THEN the System SHALL establish a bidirectional friend relationship
3. WHEN a user rejects a friend request THEN the System SHALL remove the pending request
4. WHEN a user queries their friends list THEN the System SHALL return all friends with their current sharing status and last known location
5. WHEN a user removes a friend THEN the System SHALL delete the bidirectional relationship and stop location sharing
6. WHEN a user searches for friends by mobile or username THEN the System SHALL return matching users who are not already friends
7. WHEN serializing friend data to JSON THEN the System SHALL support round-trip serialization without data loss

### Requirement 3: 设备管理

**User Story:** As a user, I want to manage my tracking devices, so that I can monitor their locations.

#### Acceptance Criteria

1. WHEN a user binds a device by ID THEN the System SHALL associate the device with the user account
2. WHEN a user queries their devices list THEN the System SHALL return all devices with current status, battery level, and last known location
3. WHEN a user unbinds a device THEN the System SHALL remove the device association and clear location history
4. WHEN a user updates device settings THEN the System SHALL persist the device name and notification preferences
5. WHEN a device reports its status THEN the System SHALL update battery level and connection status
6. WHEN serializing device data to JSON THEN the System SHALL support round-trip serialization without data loss

### Requirement 4: 位置上报与查询（PostGIS + Redis）

**User Story:** As a user, I want to share and view real-time locations, so that I can stay connected with friends and track my devices.

#### Acceptance Criteria

1. WHEN a user or device reports location THEN the System SHALL store the location as PostGIS POINT geometry with SRID 4326
2. WHEN a user queries a friend's location THEN the System SHALL first check Redis cache, then fall back to PostGIS if cache misses
3. WHEN a user queries a device's location THEN the System SHALL return the latest location with battery and status info from cache
4. WHEN a user queries location history THEN the System SHALL use PostGIS spatial index for efficient time-range queries
5. WHEN a user's sharing status is paused or hidden THEN the System SHALL NOT return their location to friends
6. WHEN multiple locations are reported in batch THEN the System SHALL use bulk insert with PostGIS geometry conversion
7. WHEN querying nearby friends THEN the System SHALL use PostGIS ST_DWithin for radius-based spatial queries
8. WHEN serializing location data to JSON THEN the System SHALL support round-trip serialization without data loss

### Requirement 4.1: 后台持续定位支持

**User Story:** As a mobile user, I want the app to track my location in the background, so that my friends can see my real-time position even when the app is not active.

#### Acceptance Criteria

1. WHEN a client reports location with background mode flag THEN the System SHALL accept and process the location normally
2. WHEN a client uses iOS significant location change mode THEN the System SHALL accept locations with lower frequency and larger accuracy radius
3. WHEN a client uses Android foreground service mode THEN the System SHALL accept high-frequency location updates
4. WHEN a client reports battery-optimized location THEN the System SHALL store the optimization mode metadata
5. WHEN a client reconnects after background suspension THEN the System SHALL accept batch upload of cached locations
6. WHEN location accuracy is below threshold THEN the System SHALL still store the location but mark it as low-accuracy

### Requirement 5: 位置共享设置

**User Story:** As a user, I want to control my location sharing preferences, so that I can protect my privacy.

#### Acceptance Criteria

1. WHEN a user enables location sharing THEN the System SHALL allow friends to see their location
2. WHEN a user pauses location sharing THEN the System SHALL temporarily hide their location from all friends
3. WHEN a user enables ghost mode THEN the System SHALL hide their location while still receiving others' locations
4. WHEN a user sets sharing preferences for specific friends THEN the System SHALL apply individual sharing rules
5. WHEN a user queries their sharing settings THEN the System SHALL return current sharing status and friend-specific rules
6. WHEN serializing settings data to JSON THEN the System SHALL support round-trip serialization without data loss

### Requirement 6: 通知与提醒

**User Story:** As a user, I want to receive notifications about location events, so that I can stay informed.

#### Acceptance Criteria

1. WHEN a friend enters or leaves a geofence THEN the System SHALL send a push notification
2. WHEN a device battery is low THEN the System SHALL send a warning notification
3. WHEN a device goes offline THEN the System SHALL send an alert notification
4. WHEN a user receives a friend request THEN the System SHALL send a notification
5. WHEN a user configures notification preferences THEN the System SHALL respect the settings for future notifications

### Requirement 7: 地理围栏管理

**User Story:** As a user, I want to create geofences, so that I can receive alerts when friends or devices enter or leave specific areas.

#### Acceptance Criteria

1. WHEN a user creates a geofence THEN the System SHALL store the center point, radius, and notification settings
2. WHEN a user queries their geofences THEN the System SHALL return all geofences with current status
3. WHEN a user updates a geofence THEN the System SHALL persist the changes
4. WHEN a user deletes a geofence THEN the System SHALL remove the geofence and stop monitoring
5. WHEN a tracked entity crosses a geofence boundary THEN the System SHALL trigger the appropriate notification
6. WHEN serializing geofence data to JSON THEN the System SHALL support round-trip serialization without data loss

### Requirement 8: 实时通信（WebSocket）

**User Story:** As a user, I want to receive real-time location updates, so that I can see live movements on the map.

#### Acceptance Criteria

1. WHEN a user connects via WebSocket THEN the System SHALL authenticate using JWT token and establish a persistent connection
2. WHEN a friend's location updates THEN the System SHALL push the update to all connected subscribers within 100ms
3. WHEN a device's location updates THEN the System SHALL push the update to the device owner's WebSocket connection
4. WHEN a WebSocket connection is lost THEN the System SHALL clean up subscriptions and mark user as offline
5. WHEN a user subscribes to specific friends or devices THEN the System SHALL only push relevant updates
6. WHEN a user reconnects after disconnect THEN the System SHALL restore previous subscriptions automatically
7. WHEN the server receives high-frequency updates THEN the System SHALL throttle pushes to max 1 update per second per entity
8. WHEN a WebSocket message fails to deliver THEN the System SHALL retry up to 3 times before dropping

### Requirement 8.1: WebSocket 连接管理

**User Story:** As a system, I want to efficiently manage WebSocket connections, so that the system can scale to handle many concurrent users.

#### Acceptance Criteria

1. WHEN a user has multiple devices connected THEN the System SHALL maintain separate WebSocket connections for each device
2. WHEN broadcasting location updates THEN the System SHALL use Redis Pub/Sub for cross-instance message distribution
3. WHEN a connection is idle for 5 minutes THEN the System SHALL send a ping frame to verify connection health
4. WHEN a client fails to respond to ping within 30 seconds THEN the System SHALL close the connection
5. WHEN the server restarts THEN the System SHALL gracefully close all connections with reconnect hint

### Requirement 9: 数据持久化与缓存（Redis + PostGIS）

**User Story:** As a system administrator, I want reliable data storage, so that user data is safe and queries are fast.

#### Acceptance Criteria

1. WHEN storing location data THEN the System SHALL use PostGIS spatial indexing (GIST) for efficient geo-queries
2. WHEN querying recent locations THEN the System SHALL use Redis cache with Geohash keys for sub-millisecond response
3. WHEN the cache misses THEN the System SHALL fall back to PostGIS and populate the cache with 5-minute TTL
4. WHEN location data is updated THEN the System SHALL invalidate relevant Redis cache entries immediately
5. WHEN serializing any entity to JSON and deserializing back THEN the System SHALL produce an equivalent entity
6. WHEN storing user's latest location THEN the System SHALL use Redis GEOADD for geo-radius queries
7. WHEN querying friends within radius THEN the System SHALL use Redis GEORADIUS before falling back to PostGIS

### Requirement 9.1: Redis 缓存策略

**User Story:** As a system, I want efficient caching, so that the system can handle high read throughput.

#### Acceptance Criteria

1. WHEN caching user's latest location THEN the System SHALL use key pattern `loc:user:{userId}` with 5-minute TTL
2. WHEN caching device's latest location THEN the System SHALL use key pattern `loc:device:{deviceId}` with 5-minute TTL
3. WHEN caching friend relationships THEN the System SHALL use key pattern `friends:{userId}` with 1-hour TTL
4. WHEN caching user settings THEN the System SHALL use key pattern `settings:{userId}` with 1-hour TTL
5. WHEN a write operation occurs THEN the System SHALL invalidate cache before database write (cache-aside pattern)
6. WHEN using Redis GEO commands THEN the System SHALL store locations in `geo:users` and `geo:devices` sorted sets
7. WHEN cache memory is limited THEN the System SHALL use LRU eviction policy



### Requirement 10: 高并发位置写入

**User Story:** As a system, I want to handle high-frequency location updates, so that the system remains responsive under load.

#### Acceptance Criteria

1. WHEN receiving location updates THEN the System SHALL use write-behind caching to batch database writes
2. WHEN location update rate exceeds threshold THEN the System SHALL sample locations to reduce storage while preserving trajectory
3. WHEN multiple users report locations simultaneously THEN the System SHALL use connection pooling and prepared statements
4. WHEN database write fails THEN the System SHALL queue the location in Redis for retry
5. WHEN processing batch locations THEN the System SHALL use PostgreSQL COPY for bulk inserts

### Requirement 11: 地理空间计算（PostGIS）

**User Story:** As a user, I want accurate distance and area calculations, so that I can understand spatial relationships.

#### Acceptance Criteria

1. WHEN calculating distance between two points THEN the System SHALL use PostGIS ST_Distance with geography type for accurate results
2. WHEN checking if a point is within a geofence THEN the System SHALL use PostGIS ST_Contains or ST_DWithin
3. WHEN finding nearby entities THEN the System SHALL use PostGIS spatial index with ST_DWithin for efficient queries
4. WHEN storing polygon geofences THEN the System SHALL use PostGIS POLYGON geometry type
5. WHEN calculating trajectory length THEN the System SHALL use PostGIS ST_Length on LINESTRING geometry
