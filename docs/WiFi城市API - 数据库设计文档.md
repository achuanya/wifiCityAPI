# WiFi城市API - 数据库设计文档

## 1. 数据库概述

WiFi城市API项目使用MySQL数据库作为其主要数据存储解决方案。数据库设计基于关系型模型，包含六个主要表，分别用于存储门店信息、WiFi配置、用户资料、扫码日志、优惠券和优惠券使用记录。

### 1.1 技术栈

- **数据库管理系统**: MySQL 8.0+
- **字符集**: utf8mb4 (支持完整的Unicode字符集，包括表情符号)
- **排序规则**: utf8mb4_unicode_ci (不区分大小写的Unicode排序规则)
- **ORM框架**: GORM v1.30.0 (Go语言的对象关系映射库)
- **驱动**: go-sql-driver/mysql v1.8.1
- **读写分离**: GORM dbresolver插件 v1.6.0
- **连接池管理**: 通过GORM配置

### 1.2 数据库架构

项目采用主从复制架构，实现读写分离，以提高性能和可靠性：
- **主库**: 处理所有写操作
- **从库**: 处理读操作，可配置多个从库实例
- **负载均衡策略**: 随机分配策略

## 2. 表设计详情

### 2.1 门店表 (store)

存储所有WiFi接入点所在的门店信息。

#### 字段定义:

| 字段名 | 类型 | 约束 | 描述 |
|--------|------|------|------|
| store_id | INT | PRIMARY KEY, AUTO_INCREMENT | 门店ID，七位数起步 |
| name | VARCHAR(100) | NOT NULL | 门店名称 |
| country | VARCHAR(64) | | 国家 |
| province | VARCHAR(64) | | 省份 |
| city | VARCHAR(64) | | 城市 |
| district | VARCHAR(64) | | 区/县 |
| address | VARCHAR(255) | | 详细地址 |
| latitude | DECIMAL(10,6) | | 门店纬度 |
| longitude | DECIMAL(10,6) | | 门店经度 |
| phone | VARCHAR(20) | | 联系电话 |
| wifi_count | INT | DEFAULT 0 | 门店WIFI数量 |
| status | TINYINT | DEFAULT 1 | 门店状态，1正常，0停用 |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP | 更新时间 |

#### 索引:
- PRIMARY KEY: `store_id`
- INDEX `idx_location`: (`province`, `city`, `district`) - 用于地理位置查询
- INDEX `idx_status`: (`status`) - 用于状态筛选

#### 自增设置:
- `AUTO_INCREMENT=1000001` - 确保ID从七位数开始

### 2.2 WIFI配置表 (wifi_config)

存储每个门店的WiFi接入点配置信息。

#### 字段定义:

| 字段名 | 类型 | 约束 | 描述 |
|--------|------|------|------|
| wifi_id | INT | PRIMARY KEY, AUTO_INCREMENT | 主键ID |
| store_id | INT | NOT NULL, FOREIGN KEY | 关联的门店ID |
| ssid | VARCHAR(64) | NOT NULL | WIFI名称 |
| password_encrypted | VARCHAR(256) | NOT NULL | 加密后的WIFI密码 |
| encryption_type | ENUM | NOT NULL, DEFAULT 'UNKNOWN' | 加密类型：WPA2/WPA3/WEP/OPEN/UNKNOWN |
| wifi_type | ENUM | NOT NULL, DEFAULT 'CUSTOMER' | WIFI类型：CUSTOMER/STAFF/EVENT/OTHER |
| max_connections | INT | DEFAULT 50 | 最大连接数限制 |
| last_updated | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP | 最后更新时间 |

#### 索引:
- PRIMARY KEY: `wifi_id`
- FOREIGN KEY: `store_id` 引用 `store(store_id)`
- UNIQUE KEY `uniq_store_ssid_type`: (`store_id`, `ssid`, `wifi_type`) - 确保同一门店下相同类型的WiFi名称唯一

### 2.3 用户信息表 (user_profile)

存储微信小程序用户信息。

#### 字段定义:

| 字段名 | 类型 | 约束 | 描述 |
|--------|------|------|------|
| user_union_id | VARCHAR(64) | PRIMARY KEY | 用户UnionID作为主键 |
| open_id | VARCHAR(64) | UNIQUE | 微信OpenID |
| wechat_nickname | VARCHAR(128) | | 微信昵称 |
| wechat_avatar_url | VARCHAR(255) | | 微信头像URL |
| phone_number | VARCHAR(20) | | 手机号（加密存储） |
| phone_country_code | VARCHAR(8) | | 手机号国家区号 |
| gender | TINYINT | | 用户性别（1男，2女，0未知） |
| language | VARCHAR(16) | | 用户语言 |
| country | VARCHAR(64) | | 用户国家 |
| province | VARCHAR(64) | | 用户省份 |
| city | VARCHAR(64) | | 用户城市 |
| first_seen | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 首次记录时间 |
| last_seen | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP | 最近更新时间 |

#### 索引:
- PRIMARY KEY: `user_union_id`
- INDEX `idx_open_id`: (`open_id`) - 用于OpenID查询
- INDEX `idx_phone_number`: (`phone_number`) - 用于手机号查询

### 2.4 扫码日志表 (scan_log)

记录用户扫描WiFi二维码和连接WiFi的行为日志。

#### 字段定义:

| 字段名 | 类型 | 约束 | 描述 |
|--------|------|------|------|
| log_id | BIGINT | PRIMARY KEY, AUTO_INCREMENT | 主键ID |
| store_id | INT | NOT NULL, FOREIGN KEY | 门店ID |
| user_union_id | VARCHAR(64) | FOREIGN KEY | 微信UnionID |
| scan_time | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 扫码时间 |
| device_info | VARCHAR(255) | | 用户设备信息 |
| ip_address | VARCHAR(45) | | 用户IP地址 |
| network_type | ENUM | | 网络类型：WIFI/5G/4G/3G/2G/UNKNOWN |
| location_lat | DECIMAL(10,6) | | 用户扫码纬度 |
| location_lng | DECIMAL(10,6) | | 用户扫码经度 |
| mini_program_version | VARCHAR(32) | | 小程序版本号 |
| success_flag | TINYINT(1) | DEFAULT 0 | 是否成功连接WiFi |
| fail_reason_code | VARCHAR(32) | | 连接失败错误码 |
| fail_reason_message | VARCHAR(255) | | 连接失败详细信息 |
| wifi_ssid | VARCHAR(64) | | 连接的WiFi名称 |
| wifi_mac | VARCHAR(64) | | 连接WiFi的MAC地址 |
| wifi_signal | TINYINT | | WiFi信号强度 |
| qr_code_type | ENUM | | 二维码类型：STORE/EVENT/POSTER/DESK/OTHER |
| qr_code_id | VARCHAR(64) | | 二维码ID |
| system_info | VARCHAR(128) | | 操作系统信息 |
| brand | VARCHAR(64) | | 设备品牌 |
| model | VARCHAR(64) | | 设备型号 |
| page_path | VARCHAR(255) | | 扫码来源页路径 |
| referer | VARCHAR(255) | | 扫码来源URL或分享来源 |
| remark | VARCHAR(255) | | 备注信息 |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 创建时间 |

#### 索引:
- PRIMARY KEY: `log_id`
- FOREIGN KEY: `store_id` 引用 `store(store_id)`
- FOREIGN KEY: `user_union_id` 引用 `user_profile(user_union_id)`
- INDEX `idx_store_scan_time`: (`store_id`, `scan_time`) - 用于门店时间范围查询
- INDEX `idx_user_union_id`: (`user_union_id`) - 用于用户查询
- INDEX `idx_success_store_time`: (`store_id`, `success_flag`, `scan_time`) - 用于成功率统计

#### 分区:
- 生产环境建议按 `scan_time` 进行表分区，以提高大数据量下的查询性能

### 2.5 优惠券表 (coupon)

存储系统中的各类优惠券定义。

#### 字段定义:

| 字段名 | 类型 | 约束 | 描述 |
|--------|------|------|------|
| coupon_id | INT | PRIMARY KEY, AUTO_INCREMENT | 优惠券ID |
| coupon_name | VARCHAR(100) | NOT NULL | 优惠券名称 |
| coupon_code | VARCHAR(32) | UNIQUE | 优惠券兑换码 |
| coupon_type | ENUM | NOT NULL | 优惠券类型：DISCOUNT/CASH/GIFT/SHIPPING |
| value | DECIMAL(10,2) | NOT NULL | 优惠券面值 |
| min_purchase_amount | DECIMAL(10,2) | DEFAULT 0.00 | 最低消费金额 |
| usage_limit_per_user | INT | DEFAULT 1 | 每用户可领取数量限制 |
| total_quantity | INT | DEFAULT 0 | 总发行量，0表示不限制 |
| issued_quantity | INT | DEFAULT 0 | 已发行数量 |
| start_time | TIMESTAMP | NOT NULL | 优惠券生效时间 |
| end_time | TIMESTAMP | NOT NULL | 优惠券过期时间 |
| validity_days | INT | | 领券后有效天数 |
| store_id | INT | FOREIGN KEY | 适用门店ID，NULL表示全平台通用 |
| description | TEXT | | 优惠券详细描述 |
| status | TINYINT | DEFAULT 1 | 优惠券状态：1启用，0禁用 |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP | 更新时间 |

#### 索引:
- PRIMARY KEY: `coupon_id`
- FOREIGN KEY: `store_id` 引用 `store(store_id)`
- INDEX `idx_coupon_time`: (`start_time`, `end_time`) - 用于有效期查询
- INDEX `idx_status`: (`status`) - 用于状态筛选

### 2.6 优惠券发放与使用日志表 (coupon_log)

记录优惠券的发放、领取、使用等操作日志。

#### 字段定义:

| 字段名 | 类型 | 约束 | 描述 |
|--------|------|------|------|
| log_id | BIGINT | PRIMARY KEY, AUTO_INCREMENT | 主键ID |
| coupon_id | INT | NOT NULL, FOREIGN KEY | 优惠券ID |
| user_union_id | VARCHAR(64) | NOT NULL, FOREIGN KEY | 用户UnionID |
| store_id | INT | FOREIGN KEY | 领取/使用门店ID |
| action_type | ENUM | NOT NULL | 行为类型：ISSUE/RECEIVE/USE/EXPIRE/REFUND |
| action_time | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 行为发生时间 |
| order_id | VARCHAR(64) | | 关联的订单ID |
| amount_deducted | DECIMAL(10,2) | | 优惠券抵扣金额 |
| status | TINYINT | DEFAULT 1 | 日志状态：1成功，0失败 |
| remark | VARCHAR(255) | | 备注信息 |

#### 索引:
- PRIMARY KEY: `log_id`
- FOREIGN KEY: `coupon_id` 引用 `coupon(coupon_id)`
- FOREIGN KEY: `user_union_id` 引用 `user_profile(user_union_id)`
- FOREIGN KEY: `store_id` 引用 `store(store_id)`
- INDEX `idx_user_coupon`: (`user_union_id`, `coupon_id`) - 用于查询用户优惠券
- INDEX `idx_action_time`: (`action_time`) - 用于时间范围查询
- INDEX `idx_coupon_action_status`: (`coupon_id`, `action_type`, `status`) - 用于优惠券状态统计

### 2.7 小程序配置表 (app_config)

存储小程序的配置信息。

#### 字段定义:

| 字段名 | 类型 | 约束 | 描述 |
|--------|------|------|------|
| config_id | INT | PRIMARY KEY, AUTO_INCREMENT | 配置ID |
| store_id | INT | NOT NULL, FOREIGN KEY | 门店ID |
| mini_program_id | VARCHAR(64) | NOT NULL | 小程序AppID |
| access_token | VARCHAR(255) | | 调用凭证 |
| token_expiry | TIMESTAMP | | 令牌过期时间 |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP | 更新时间 |

#### 索引:
- PRIMARY KEY: `config_id`
- FOREIGN KEY: `store_id` 引用 `store(store_id)`
- UNIQUE KEY `uniq_store_mini_program`: (`store_id`, `mini_program_id`) - 确保门店的小程序配置唯一

## 3. 表关系

### 3.1 一对多关系
- **Store -> WifiConfig**: 一个门店可以有多个WiFi配置
- **Store -> ScanLog**: 一个门店可以有多个扫码记录
- **Store -> Coupon**: 一个门店可以发布多个优惠券
- **UserProfile -> ScanLog**: 一个用户可以有多个扫码记录

### 3.2 多对多关系
- **User-Coupon**: 通过coupon_log表实现用户和优惠券的多对多关系

## 4. 数据安全

### 4.1 密码加密
系统使用AES-256-GCM加密算法对WiFi密码进行加密存储：
- 密码加密使用AES-256-GCM (Galois/Counter Mode)
- 包含初始化向量(IV)和认证标签(Tag)
- 加密后的数据以Base64编码存储

### 4.2 API安全
- 使用HMAC-SHA256进行API请求签名验证
- 使用时间戳窗口(5分钟)防止重放攻击

## 5. 性能优化

### 5.1 连接池配置
```yaml
settings:
  max_idle_conns: 10      # 最大空闲连接数
  max_open_conns: 100     # 最大打开连接数
  conn_max_idle_time: 10m # 连接最大空闲时间
  conn_max_lifetime: 1h   # 连接最大生命周期
```

### 5.2 读写分离
通过GORM的dbresolver插件实现读写分离：
- 主库处理所有写操作
- 从库处理读操作
- 使用随机策略进行负载均衡

### 5.3 大表优化
- 对scan_log表按时间进行分区，提高查询效率
- 为频繁查询的字段创建合适的索引

## 6. ORM映射

系统使用GORM框架实现对象关系映射，主要特性：
- 自动创建和迁移表结构
- 关联关系处理
- 事务支持
- 钩子函数(Hooks)
- 预加载(Preload)
- 自定义数据类型

示例映射代码：
```go
// Store 对应于 store 表的 GORM 模型
type Store struct {
    StoreID     uint         `gorm:"primaryKey;autoIncrement;comment:门店ID，七位数起步"`
    Name        string       `gorm:"type:varchar(100);not null;comment:门店名称"`
    // ... 其他字段
    WifiConfigs []WifiConfig `gorm:"foreignKey:StoreID"` // 一对多关系
    ScanLogs    []ScanLog    `gorm:"foreignKey:StoreID"` // 一对多关系
    Coupons     []Coupon     `gorm:"foreignKey:StoreID"` // 一对多关系
}
```

## 7. 未来扩展考虑

- **分库分表**: 当数据量增长到一定规模时，考虑实施分库分表策略
- **缓存层**: 引入Redis缓存热点数据，减轻数据库压力
- **全文索引**: 针对搜索需求，可考虑增加全文索引或引入专门的搜索引擎
- **数据归档**: 对历史数据进行归档，提高活跃数据的查询效率 