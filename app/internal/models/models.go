package models

import (
	"time"
)

// Store 对应于 store 表的 GORM 模型
type Store struct {
	StoreID     uint         `gorm:"primaryKey;autoIncrement;comment:门店ID，七位数起步"`
	Name        string       `gorm:"type:varchar(100);not null;comment:门店名称"`
	Country     string       `gorm:"type:varchar(64);comment:国家"`
	Province    string       `gorm:"type:varchar(64);comment:省份"`
	City        string       `gorm:"type:varchar(64);comment:城市"`
	District    string       `gorm:"type:varchar(64);comment:区/县"`
	Address     string       `gorm:"type:varchar(255);comment:详细地址"`
	Latitude    float64      `gorm:"type:decimal(10,6);comment:门店纬度"`
	Longitude   float64      `gorm:"type:decimal(10,6);comment:门店经度"`
	Phone       string       `gorm:"type:varchar(20);comment:联系电话"`
	WifiCount   int          `gorm:"default:0;comment:门店WIFI数量"`
	Status      int8         `gorm:"type:tinyint;default:1;comment:门店状态，1正常，0停用"`
	CreatedAt   time.Time    `gorm:"comment:创建时间"`
	UpdatedAt   time.Time    `gorm:"comment:更新时间"`
	WifiConfigs []WifiConfig `gorm:"foreignKey:StoreID"` // 一对多关系
	ScanLogs    []ScanLog    `gorm:"foreignKey:StoreID"` // 一对多关系
	Coupons     []Coupon     `gorm:"foreignKey:StoreID"` // 一对多关系
}

func (Store) TableName() string {
	return "store"
}

// WifiConfig 对应于 wifi_config 表的 GORM 模型
type WifiConfig struct {
	WifiID            uint      `gorm:"primaryKey;autoIncrement;comment:主键ID"`
	StoreID           uint      `gorm:"not null;comment:门店ID"`
	SSID              string    `gorm:"type:varchar(64);not null;comment:WIFI名称"`
	PasswordEncrypted string    `gorm:"type:varchar(256);not null;comment:加密后的WIFI密码"`
	EncryptionType    string    `gorm:"type:enum('WPA2','WPA3','WEP','OPEN','UNKNOWN');default:'UNKNOWN';not null;comment:加密类型"`
	WifiType          string    `gorm:"type:enum('CUSTOMER','STAFF','EVENT','OTHER');default:'CUSTOMER';not null;comment:WIFI类型"`
	MaxConnections    int       `gorm:"default:50;comment:最大连接数限制"`
	LastUpdated       time.Time `gorm:"autoUpdateTime;comment:最后更新时间"`
}

func (WifiConfig) TableName() string {
	return "wifi_config"
}

// UserProfile 对应于 user_profile 表的 GORM 模型
type UserProfile struct {
	UserUnionID      string    `gorm:"primaryKey;type:varchar(64);comment:用户UnionID作为主键"`
	OpenID           string    `gorm:"type:varchar(64);unique;comment:微信OpenID"`
	WechatNickname   string    `gorm:"type:varchar(128);comment:微信昵称"`
	WechatAvatarURL  string    `gorm:"type:varchar(255);comment:微信头像URL"`
	PhoneNumber      string    `gorm:"type:varchar(20);comment:手机号（纯数字，建议加密存储）"`
	PhoneCountryCode string    `gorm:"type:varchar(8);comment:手机号国家区号"`
	Gender           int8      `gorm:"type:tinyint;comment:用户性别（1男，2女，0未知）"`
	Language         string    `gorm:"type:varchar(16);comment:用户语言"`
	Country          string    `gorm:"type:varchar(64);comment:用户国家"`
	Province         string    `gorm:"type:varchar(64);comment:用户省份"`
	City             string    `gorm:"type:varchar(64);comment:用户城市"`
	FirstSeen        time.Time `gorm:"autoCreateTime;comment:首次记录时间"`
	LastSeen         time.Time `gorm:"autoUpdateTime;comment:最近更新时间"`
	ScanLogs         []ScanLog `gorm:"foreignKey:UserUnionID"` // 一对多关系
}

func (UserProfile) TableName() string {
	return "user_profile"
}

// ScanLog 对应于 scan_log 表的 GORM 模型
type ScanLog struct {
	LogID              uint64    `gorm:"primaryKey;autoIncrement;comment:主键ID"`
	StoreID            uint      `gorm:"not null;comment:门店ID"`
	UserUnionID        string    `gorm:"type:varchar(64);comment:微信UnionID"`
	ScanTime           time.Time `gorm:"autoCreateTime;comment:扫码时间"`
	DeviceInfo         string    `gorm:"type:varchar(255);comment:用户设备信息"`
	IPAddress          string    `gorm:"type:varchar(45);comment:用户IP地址"`
	NetworkType        string    `gorm:"type:enum('WIFI','5G','4G','3G','2G','UNKNOWN');comment:用户扫码时网络类型"`
	LocationLat        float64   `gorm:"type:decimal(10,6);comment:用户扫码纬度"`
	LocationLng        float64   `gorm:"type:decimal(10,6);comment:用户扫码经度"`
	MiniProgramVersion string    `gorm:"type:varchar(32);comment:小程序版本号"`
	SuccessFlag        bool      `gorm:"type:tinyint(1);default:0;comment:是否成功连接WiFi"`
	FailReasonCode     string    `gorm:"type:varchar(32);comment:连接失败错误码"`
	FailReasonMessage  string    `gorm:"type:varchar(255);comment:连接失败详细信息"`
	WifiSSID           string    `gorm:"type:varchar(64);comment:连接的WiFi名称"`
	WifiMac            string    `gorm:"type:varchar(64);comment:连接WiFi的MAC地址"`
	WifiSignal         int8      `gorm:"type:tinyint;comment:WiFi信号强度"`
	QrCodeType         string    `gorm:"type:enum('STORE','EVENT','POSTER','DESK','OTHER');comment:二维码类型"`
	QrCodeID           string    `gorm:"type:varchar(64);comment:二维码ID"`
	SystemInfo         string    `gorm:"type:varchar(128);comment:操作系统信息"`
	Brand              string    `gorm:"type:varchar(64);comment:设备品牌"`
	Model              string    `gorm:"type:varchar(64);comment:设备型号"`
	PagePath           string    `gorm:"type:varchar(255);comment:扫码来源页路径"`
	Referer            string    `gorm:"type:varchar(255);comment:扫码来源URL或分享来源"`
	Remark             string    `gorm:"type:varchar(255);comment:备注信息"`
	CreatedAt          time.Time `gorm:"comment:创建时间"`
}

func (ScanLog) TableName() string {
	return "scan_log"
}

// Coupon 对应于 coupon 表的 GORM 模型
type Coupon struct {
	CouponID          uint      `gorm:"primaryKey;autoIncrement;comment:优惠券ID"`
	CouponName        string    `gorm:"type:varchar(100);not null;comment:优惠券名称"`
	CouponCode        string    `gorm:"type:varchar(32);unique;comment:优惠券兑换码"`
	CouponType        string    `gorm:"type:enum('DISCOUNT','CASH','GIFT','SHIPPING');not null;comment:优惠券类型"`
	Value             float64   `gorm:"type:decimal(10,2);not null;comment:优惠券面值"`
	MinPurchaseAmount float64   `gorm:"type:decimal(10,2);default:0.00;comment:最低消费金额"`
	UsageLimitPerUser int       `gorm:"default:1;comment:每个用户可领取的最大数量"`
	TotalQuantity     int       `gorm:"default:0;comment:优惠券总发行量"`
	IssuedQuantity    int       `gorm:"default:0;comment:已发行数量"`
	StartTime         time.Time `gorm:"not null;comment:优惠券生效时间"`
	EndTime           time.Time `gorm:"not null;comment:优惠券过期时间"`
	ValidityDays      int       `gorm:"comment:领券后有效天数"`
	StoreID           *uint     `gorm:"comment:适用门店ID"` // 使用指针以接受 NULL 值
	Description       string    `gorm:"type:text;comment:优惠券详细描述"`
	Status            int8      `gorm:"type:tinyint;default:1;comment:优惠券状态"`
	CreatedAt         time.Time `gorm:"comment:创建时间"`
	UpdatedAt         time.Time `gorm:"comment:更新时间"`
}

func (Coupon) TableName() string {
	return "coupon"
}

// CouponLog 对应于 coupon_log 表的 GORM 模型
type CouponLog struct {
	LogID          uint64    `gorm:"primaryKey;autoIncrement;comment:主键ID"`
	CouponID       uint      `gorm:"not null;comment:优惠券ID"`
	UserUnionID    string    `gorm:"type:varchar(64);not null;comment:用户UnionID"`
	StoreID        *uint     `gorm:"comment:领取/使用门店ID"` // 使用指针以接受 NULL 值
	ActionType     string    `gorm:"type:enum('ISSUE','RECEIVE','USE','EXPIRE','REFUND');not null;comment:行为类型"`
	ActionTime     time.Time `gorm:"autoCreateTime;comment:行为发生时间"`
	OrderID        string    `gorm:"type:varchar(64);comment:关联的订单ID"`
	AmountDeducted float64   `gorm:"type:decimal(10,2);comment:优惠券抵扣金额"`
	Status         int8      `gorm:"type:tinyint;default:1;comment:日志状态"`
	Remark         string    `gorm:"type:varchar(255);comment:备注信息"`
}

func (CouponLog) TableName() string {
	return "coupon_log"
}

// AppConfig 对应于 app_config 表的 GORM 模型
type AppConfig struct {
	ConfigID      uint      `gorm:"primaryKey;autoIncrement"`
	StoreID       uint      `gorm:"not null"`
	MiniProgramID string    `gorm:"type:varchar(64);not null;comment:小程序AppID"`
	AccessToken   string    `gorm:"type:varchar(255);comment:调用凭证"`
	TokenExpiry   time.Time `gorm:"comment:令牌过期时间"`
	CreatedAt     time.Time `gorm:"comment:创建时间"`
	UpdatedAt     time.Time `gorm:"comment:更新时间"`
}

func (AppConfig) TableName() string {
	return "app_config"
}
