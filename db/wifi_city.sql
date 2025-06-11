-- wifi_city
CREATE DATABASE IF NOT EXISTS wifi_city DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE wifi_city;

-- 门店表 store
CREATE TABLE store (
    store_id INT PRIMARY KEY AUTO_INCREMENT COMMENT '门店ID，七位数起步',
    name VARCHAR(100) NOT NULL COMMENT '门店名称',
    country VARCHAR(64) COMMENT '国家',
    province VARCHAR(64) COMMENT '省份',
    city VARCHAR(64) COMMENT '城市',
    district VARCHAR(64) COMMENT '区/县',
    address VARCHAR(255) COMMENT '详细地址',
    latitude DECIMAL(10,6) COMMENT '门店纬度',
    longitude DECIMAL(10,6) COMMENT '门店经度',
    phone VARCHAR(20) COMMENT '联系电话',
    wifi_count INT DEFAULT 0 COMMENT '门店WIFI数量',
    status TINYINT DEFAULT 1 COMMENT '门店状态，1正常，0停用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_location (province, city, district),
    INDEX idx_status (status)
) AUTO_INCREMENT=1000001 COMMENT='门店表';

-- WIFI配置表 wifi_config
CREATE TABLE wifi_config (
    wifi_id INT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    store_id INT NOT NULL,
    ssid VARCHAR(64) NOT NULL COMMENT 'WIFI名称',
    password_encrypted VARCHAR(256) NOT NULL COMMENT '加密后的WIFI密码',
    encryption_type ENUM('WPA2', 'WPA3', 'WEP', 'OPEN', 'UNKNOWN') DEFAULT 'UNKNOWN' NOT NULL COMMENT '加密类型',
    wifi_type ENUM('CUSTOMER', 'STAFF', 'EVENT', 'OTHER') DEFAULT 'CUSTOMER' NOT NULL COMMENT 'WIFI类型：CUSTOMER顾客WIFI，STAFF员工WIFI，EVENT活动WIFI等',
    max_connections INT DEFAULT 50 COMMENT '最大连接数限制',
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (store_id) REFERENCES store(store_id),
    UNIQUE KEY uniq_store_ssid_type (store_id, ssid, wifi_type) -- 确保门店下每种类型的Wi-Fi名称唯一
) COMMENT='WIFI配置表';

-- 用户信息表 user_profile
CREATE TABLE user_profile (
    user_union_id VARCHAR(64) PRIMARY KEY COMMENT '用户UnionID作为主键',
    open_id VARCHAR(64) UNIQUE COMMENT '微信OpenID',
    wechat_nickname VARCHAR(128) COMMENT '微信昵称',
    wechat_avatar_url VARCHAR(255) COMMENT '微信头像URL',
    phone_number VARCHAR(20) COMMENT '手机号（纯数字，建议加密存储）',
    phone_country_code VARCHAR(8) COMMENT '手机号国家区号，例如86',
    gender TINYINT COMMENT '用户性别（1男，2女，0未知）',
    language VARCHAR(16) COMMENT '用户语言，如zh_CN',
    country VARCHAR(64) COMMENT '用户国家',
    province VARCHAR(64) COMMENT '用户省份',
    city VARCHAR(64) COMMENT '用户城市',
    first_seen TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '首次记录时间',
    last_seen TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最近更新时间',
    INDEX idx_open_id (open_id),
    INDEX idx_phone_number (phone_number)
) COMMENT='用户信息表';

-- 扫码日志表 scan_log
CREATE TABLE scan_log (
    log_id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    store_id INT NOT NULL COMMENT '门店ID，外键',
    user_union_id VARCHAR(64) COMMENT '微信UnionID，关联user_profile表',
    
    scan_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '扫码时间',
    device_info VARCHAR(255) COMMENT '用户设备信息（机型、系统等）',
    ip_address VARCHAR(45) COMMENT '用户IP地址',
    
    network_type ENUM('WIFI', '5G', '4G', '3G', '2G', 'UNKNOWN') COMMENT '用户扫码时网络类型',
    
    location_lat DECIMAL(10,6) COMMENT '用户扫码纬度',
    location_lng DECIMAL(10,6) COMMENT '用户扫码经度',
    mini_program_version VARCHAR(32) COMMENT '小程序版本号',
    success_flag TINYINT(1) DEFAULT 0 COMMENT '是否成功连接WiFi（0失败，1成功）',
    fail_reason_code VARCHAR(32) COMMENT '连接失败错误码，例如密码错误、连接超时、设备不支持等',
    fail_reason_message VARCHAR(255) COMMENT '连接失败详细信息或用户提示语',
    wifi_ssid VARCHAR(64) COMMENT '连接的WiFi名称',
    wifi_mac VARCHAR(64) COMMENT '连接WiFi的MAC地址',
    wifi_signal TINYINT COMMENT 'WiFi信号强度',

    qr_code_type ENUM('STORE', 'EVENT', 'POSTER', 'DESK', 'OTHER') COMMENT '二维码类型',
    qr_code_id VARCHAR(64) COMMENT '二维码ID',
    
    -- 设备及环境信息
    system_info VARCHAR(128) COMMENT '操作系统信息',
    brand VARCHAR(64) COMMENT '设备品牌',
    model VARCHAR(64) COMMENT '设备型号',
    
    -- 业务分析可用
    page_path VARCHAR(255) COMMENT '扫码来源页路径',
    referer VARCHAR(255) COMMENT '扫码来源URL或分享来源',
    remark VARCHAR(255) COMMENT '备注信息',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    
    FOREIGN KEY (store_id) REFERENCES store(store_id),
    FOREIGN KEY (user_union_id) REFERENCES user_profile(user_union_id), 
    
    INDEX idx_store_scan_time (store_id, scan_time),
    INDEX idx_user_union_id (user_union_id),
    INDEX idx_success_store_time (store_id, success_flag, scan_time) 
) COMMENT='扫码日志表';
-- 生产环境对此表按 scan_time 进行分区

-- 优惠券表 coupon
CREATE TABLE coupon (
    coupon_id INT PRIMARY KEY AUTO_INCREMENT COMMENT '优惠券ID',
    coupon_name VARCHAR(100) NOT NULL COMMENT '优惠券名称',
    coupon_code VARCHAR(32) UNIQUE COMMENT '优惠券兑换码，可空，如果需要兑换码则填写',
    coupon_type ENUM('DISCOUNT', 'CASH', 'GIFT', 'SHIPPING') NOT NULL COMMENT '优惠券类型：DISCOUNT折扣券, CASH现金券, GIFT礼品券, SHIPPING运费券',
    value DECIMAL(10, 2) NOT NULL COMMENT '优惠券面值：如折扣率（0.8代表8折），或现金金额（10代表10元）',
    min_purchase_amount DECIMAL(10, 2) DEFAULT 0.00 COMMENT '最低消费金额，达到此金额才可使用',
    usage_limit_per_user INT DEFAULT 1 COMMENT '每个用户可领取的最大数量，0表示不限制',
    total_quantity INT DEFAULT 0 COMMENT '优惠券总发行量，0表示不限制',
    issued_quantity INT DEFAULT 0 COMMENT '已发行数量',
    start_time TIMESTAMP NOT NULL COMMENT '优惠券生效时间',
    end_time TIMESTAMP NOT NULL COMMENT '优惠券过期时间',
    validity_days INT COMMENT '领券后有效天数（与start_time/end_time二选一，如果设置了此项，则按领取时间+天数计算有效期）',
    store_id INT COMMENT '适用门店ID，可空，表示全平台通用；非空表示仅限特定门店',
    description TEXT COMMENT '优惠券详细描述或使用说明',
    status TINYINT DEFAULT 1 COMMENT '优惠券状态：1启用，0禁用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    FOREIGN KEY (store_id) REFERENCES store(store_id),
    INDEX idx_coupon_time (start_time, end_time),
    INDEX idx_status (status)
) COMMENT='优惠券表';

-- 优惠券发放与使用日志表 coupon_log
CREATE TABLE coupon_log (
    log_id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    coupon_id INT NOT NULL COMMENT '优惠券ID，关联coupon表',
    user_union_id VARCHAR(64) NOT NULL COMMENT '用户UnionID，关联user_profile表',
    store_id INT COMMENT '领取/使用门店ID，可空，如果优惠券是全平台通用',
    action_type ENUM('ISSUE', 'RECEIVE', 'USE', 'EXPIRE', 'REFUND') NOT NULL COMMENT '行为类型：ISSUE发放, RECEIVE领取, USE使用, EXPIRE过期, REFUND退券',
    action_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '行为发生时间',
    order_id VARCHAR(64) COMMENT '关联的订单ID，如果优惠券用于支付',
    amount_deducted DECIMAL(10, 2) COMMENT '优惠券抵扣金额',
    status TINYINT DEFAULT 1 COMMENT '日志状态：1成功，0失败（例如领取失败，使用失败等）',
    remark VARCHAR(255) COMMENT '备注信息，如失败原因',
    FOREIGN KEY (coupon_id) REFERENCES coupon(coupon_id),
    FOREIGN KEY (user_union_id) REFERENCES user_profile(user_union_id),
    FOREIGN KEY (store_id) REFERENCES store(store_id),
    INDEX idx_user_coupon (user_union_id, coupon_id),
    INDEX idx_action_time (action_time),
    INDEX idx_coupon_action_status (coupon_id, action_type, status)
) COMMENT='优惠券发放与使用日志表';

-- 小程序配置表 app_config
CREATE TABLE app_config (
    config_id INT PRIMARY KEY AUTO_INCREMENT,
    store_id INT NOT NULL,
    mini_program_id VARCHAR(64) NOT NULL COMMENT '小程序AppID',
    access_token VARCHAR(255) COMMENT '调用凭证',
    token_expiry TIMESTAMP COMMENT '令牌过期时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (store_id) REFERENCES store(store_id),
    UNIQUE KEY uniq_store_mini_program (store_id, mini_program_id)
) COMMENT='小程序配置表';