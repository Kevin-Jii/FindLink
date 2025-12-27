-- Location Sharing Backend Database Migration (MySQL 8.0+)

-- 注意：MySQL 不需要像 PostGIS 那样显式开启 Extension，
-- 但必须使用 MySQL 8.0 或更高版本以支持 SRID 4326 和 ST_Distance_Sphere 等计算。

-- User Locations Table
CREATE TABLE IF NOT EXISTS user_locations (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    -- MySQL 8.0 使用 SRID 4326 对应 GPS 坐标，必须指定 NOT NULL 才能建立空间索引
    location POINT NOT NULL SRID 4326, 
    accuracy FLOAT,
    altitude FLOAT,
    speed FLOAT,
    bearing FLOAT,
    battery_level INT,
    location_mode VARCHAR(20) DEFAULT 'foreground',
    is_low_accuracy BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_user_locations_user_id (user_id),
    INDEX idx_user_locations_created_at (created_at DESC),
    SPATIAL INDEX idx_user_locations_location (location)
) ENGINE=InnoDB;

-- Device Locations Table
CREATE TABLE IF NOT EXISTS device_locations (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    device_id VARCHAR(64) NOT NULL,
    location POINT NOT NULL SRID 4326,
    accuracy FLOAT,
    battery_level INT,
    connection_status VARCHAR(20) DEFAULT 'unknown',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_device_locations_device_id (device_id),
    INDEX idx_device_locations_created_at (created_at DESC),
    SPATIAL INDEX idx_device_locations_location (location)
) ENGINE=InnoDB;

-- Friend Relationships Table
CREATE TABLE IF NOT EXISTS friends (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    friend_id BIGINT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    sharing_status VARCHAR(20) DEFAULT 'sharing',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uniq_user_friend (user_id, friend_id),
    INDEX idx_friends_user_id (user_id),
    INDEX idx_friends_friend_id (friend_id),
    INDEX idx_friends_status (status)
) ENGINE=InnoDB;

-- Friend Requests Table
CREATE TABLE IF NOT EXISTS friend_requests (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    from_user_id BIGINT NOT NULL,
    to_user_id BIGINT NOT NULL,
    message VARCHAR(255),
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uniq_request (from_user_id, to_user_id),
    INDEX idx_friend_requests_from_user (from_user_id),
    INDEX idx_friend_requests_to_user (to_user_id),
    INDEX idx_friend_requests_status (status)
) ENGINE=InnoDB;

-- Devices Table
CREATE TABLE IF NOT EXISTS devices (
    id VARCHAR(64) PRIMARY KEY,
    user_id BIGINT NOT NULL,
    name VARCHAR(100),
    type VARCHAR(50) DEFAULT 'unknown',
    battery_level INT DEFAULT 0,
    connection_status VARCHAR(20) DEFAULT 'unknown',
    notification_enabled BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_devices_user_id (user_id)
) ENGINE=InnoDB;

-- Geofences Table
CREATE TABLE IF NOT EXISTS geofences (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    name VARCHAR(100) NOT NULL,
    center POINT NOT NULL SRID 4326,
    radius_meters FLOAT NOT NULL,
    notify_on_enter BOOLEAN DEFAULT TRUE,
    notify_on_exit BOOLEAN DEFAULT TRUE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_geofences_user_id (user_id),
    SPATIAL INDEX idx_geofences_center (center)
) ENGINE=InnoDB;

-- Geofence Events Table
CREATE TABLE IF NOT EXISTS geofence_events (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    geofence_id BIGINT NOT NULL,
    entity_type VARCHAR(20) NOT NULL,
    entity_id BIGINT NOT NULL,
    event_type VARCHAR(20) NOT NULL,
    -- 这里允许 NULL，所以不强制加 SPATIAL INDEX，如果需要加则必须设为 NOT NULL
    location POINT SRID 4326, 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_geofence_events_geofence (geofence_id),
    INDEX idx_geofence_events_entity (entity_type, entity_id),
    INDEX idx_geofence_events_created_at (created_at DESC)
) ENGINE=InnoDB;

-- User Settings Table
CREATE TABLE IF NOT EXISTS user_settings (
    user_id BIGINT PRIMARY KEY,
    share_location BOOLEAN DEFAULT TRUE,
    ghost_mode BOOLEAN DEFAULT FALSE,
    smart_alerts BOOLEAN DEFAULT TRUE,
    sos_alerts BOOLEAN DEFAULT TRUE,
    map_style VARCHAR(20) DEFAULT 'dark',
    distance_unit VARCHAR(10) DEFAULT 'km',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB;

-- ---------------------------------------------------------
-- Locations batch insert procedure
-- MySQL 不支持在 Function 中进行 INSERT 操作，必须使用 PROCEDURE。
-- 使用 JSON_TABLE (MySQL 8.0 特性) 解析 JSON 数组，效率极高。
-- ---------------------------------------------------------

DELIMITER //

CREATE PROCEDURE insert_user_locations_batch(
    IN p_user_id BIGINT,
    IN p_locations JSON,
    OUT p_affected_rows INT
)
BEGIN
    INSERT INTO user_locations (
        user_id, location, accuracy, altitude, speed, bearing, 
        battery_level, location_mode, is_low_accuracy
    )
    SELECT 
        p_user_id,
        -- 使用 ST_SRID 和 POINT 构建地理对象
        ST_SRID(POINT(jt.longitude, jt.latitude), 4326),
        jt.accuracy,
        jt.altitude,
        jt.speed,
        jt.bearing,
        jt.battery_level,
        COALESCE(jt.location_mode, 'foreground'),
        COALESCE(jt.is_low_accuracy, 0)
    FROM JSON_TABLE(
        p_locations, 
        '$[*]' COLUMNS (
            longitude DOUBLE PATH '$.longitude',
            latitude DOUBLE PATH '$.latitude',
            accuracy FLOAT PATH '$.accuracy',
            altitude FLOAT PATH '$.altitude',
            speed FLOAT PATH '$.speed',
            bearing FLOAT PATH '$.bearing',
            battery_level INT PATH '$.battery_level',
            location_mode VARCHAR(20) PATH '$.location_mode',
            is_low_accuracy BOOLEAN PATH '$.is_low_accuracy'
        )
    ) AS jt;

    -- 获取插入的行数
    SET p_affected_rows = ROW_COUNT();
END //

DELIMITER ;