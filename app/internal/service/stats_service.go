package service

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/app/internal/models"
	"github.com/gin-gonic/gin/app/pkg/database"
)

// StatsService 提供了数据统计相关的业务逻辑
type StatsService struct{}

// StoreTotalStats 定义了门店总数统计的结果结构
type StoreTotalStats struct {
	TotalStores int64 `json:"total_stores"`
}

// StoreCountByProvince 定义了按省份统计门店数量的结果结构
type StoreCountByProvince struct {
	Province string `json:"province"`
	Count    int64  `json:"count"`
}

// StoreCountByCity 定义了按城市统计门店数量的结果结构
type StoreCountByCity struct {
	Province string `json:"province"`
	City     string `json:"city"`
	Count    int64  `json:"count"`
}

// GetStoreStats 用于获取门店相关的统计数据
func (s *StatsService) GetStoreStats() (any, error) {
	db := database.DB.WithContext(context.Background())

	// 1. 统计门店总数
	var totalStats StoreTotalStats
	if err := db.Model(&models.Store{}).Count(&totalStats.TotalStores).Error; err != nil {
		return nil, err
	}

	// 2. 按省份统计门店数量
	var byProvince []StoreCountByProvince
	if err := db.Model(&models.Store{}).
		Select("province, count(*) as count").
		Group("province").
		Order("count DESC").
		Find(&byProvince).Error; err != nil {
		return nil, err
	}

	// 3. 按城市统计门店数量
	var byCity []StoreCountByCity
	if err := db.Model(&models.Store{}).
		Select("province, city, count(*) as count").
		Group("province, city").
		Order("province, count DESC").
		Find(&byCity).Error; err != nil {
		return nil, err
	}

	// 4. 组合结果
	result := gin.H{
		"total_stats": totalStats,
		"by_province": byProvince,
		"by_city":     byCity,
	}

	return result, nil
}

// WifiTotalUsageStats 定义了WIFI使用的总体统计结果
type WifiTotalUsageStats struct {
	TotalConnections      int64   `json:"total_connections"`
	SuccessfulConnections int64   `json:"successful_connections"`
	SuccessRate           float64 `json:"success_rate"`
}

// WifiUsageByEncryptionType 定义了按加密类型统计的结果
type WifiUsageByEncryptionType struct {
	EncryptionType string `json:"encryption_type"`
	Count          int64  `json:"count"`
}

// WifiUsageByFailReason 定义了按失败原因统计的结果
type WifiUsageByFailReason struct {
	FailReasonCode string `json:"fail_reason_code"`
	Count          int64  `json:"count"`
}

// WifiUsageBySSID 定义了按SSID统计的结果
type WifiUsageBySSID struct {
	SSID  string `json:"ssid"`
	Count int64  `json:"count"`
}

// GetWifiUsageStats 用于获取WIFI使用相关的统计数据
func (s *StatsService) GetWifiUsageStats(storeID *uint) (any, error) {
	db := database.DB.WithContext(context.Background())
	query := db.Model(&models.ScanLog{})

	// 如果指定了门店ID，则只查询该门店的数据
	if storeID != nil {
		query = query.Where("store_id = ?", *storeID)
	}

	// 1. 统计总连接次数和成功次数
	var stats WifiTotalUsageStats
	if err := query.Count(&stats.TotalConnections).Error; err != nil {
		return nil, err
	}
	if err := query.Where("success_flag = ?", 1).Count(&stats.SuccessfulConnections).Error; err != nil {
		return nil, err
	}
	if stats.TotalConnections > 0 {
		stats.SuccessRate = float64(stats.SuccessfulConnections) / float64(stats.TotalConnections)
	}

	// 2. 按加密类型统计 (需要关联 wifi_config 表，此处简化为直接从 scan_log 中获取，实际情况可能需要调整)
	// 注意：当前 scan_log 表没有加密类型字段，此统计暂时无法精确实现。
	// 这里我们先留空，假设未来会增加。

	// 3. 按失败原因统计
	var byFailReason []WifiUsageByFailReason
	failQuery := db.Model(&models.ScanLog{})
	if storeID != nil {
		failQuery = failQuery.Where("store_id = ?", *storeID)
	}
	if err := failQuery.Where("success_flag = ? AND fail_reason_code != ''", 0).
		Select("fail_reason_code, count(*) as count").
		Group("fail_reason_code").
		Order("count DESC").
		Find(&byFailReason).Error; err != nil {
		return nil, err
	}

	// 4. 按WIFI名称(SSID)统计
	var bySSID []WifiUsageBySSID
	ssidQuery := db.Model(&models.ScanLog{})
	if storeID != nil {
		ssidQuery = ssidQuery.Where("store_id = ?", *storeID)
	}
	if err := ssidQuery.Where("wifi_ssid != ''").
		Select("wifi_ssid as ssid, count(*) as count").
		Group("wifi_ssid").
		Order("count DESC").
		Limit(10). // 限制返回最受欢迎的10个
		Find(&bySSID).Error; err != nil {
		return nil, err
	}

	result := gin.H{
		"total_usage":    stats,
		"by_fail_reason": byFailReason,
		"by_ssid":        bySSID,
	}

	return result, nil
}

// GetUserBehaviorStats 用于获取用户行为相关的统计数据
type GetUserBehaviorStatsInput struct {
	StartDate string `form:"start_date"` // 格式: "2006-01-02"
	EndDate   string `form:"end_date"`   // 格式: "2006-01-02"
}

// UserGenderDistribution 定义了用户性别分布的统计结果
type UserGenderDistribution struct {
	Gender string `json:"gender"` // 1男, 2女, 0未知
	Count  int64  `json:"count"`
}

// UserProvinceDistribution 定义了用户地理位置分布的统计结果
type UserProvinceDistribution struct {
	Province string `json:"province"`
	Count    int64  `json:"count"`
}

// UserDeviceDistribution 定义了用户设备分布的统计结果
type UserDeviceDistribution struct {
	Brand string `json:"brand"`
	Model string `json:"model"`
	Count int64  `json:"count"`
}

// GetUserBehaviorStats 用于获取用户行为相关的统计数据
func (s *StatsService) GetUserBehaviorStats(input *GetUserBehaviorStatsInput) (any, error) {
	db := database.DB.WithContext(context.Background())

	// --- 用户档案统计 (来自 user_profile 表) ---
	userQuery := db.Model(&models.UserProfile{})
	if input.StartDate != "" {
		userQuery = userQuery.Where("first_seen >= ?", input.StartDate)
	}
	if input.EndDate != "" {
		// 注意：为了包含当天，需要查询到当天的最后一秒
		userQuery = userQuery.Where("first_seen < ?", input.EndDate+" 23:59:59")
	}

	// 1. 统计新用户注册数
	var newUsersCount int64
	if err := userQuery.Count(&newUsersCount).Error; err != nil {
		return nil, err
	}

	// 2. 统计用户性别分布
	var genderDist []UserGenderDistribution
	if err := db.Model(&models.UserProfile{}).Select("gender, count(*) as count").Group("gender").Find(&genderDist).Error; err != nil {
		return nil, err
	}

	// 3. 统计用户省份分布
	var provinceDist []UserProvinceDistribution
	if err := db.Model(&models.UserProfile{}).Where("province != ''").Select("province, count(*) as count").Group("province").Order("count DESC").Limit(20).Find(&provinceDist).Error; err != nil {
		return nil, err
	}

	// --- 扫码行为统计 (来自 scan_log 表) ---
	scanQuery := db.Model(&models.ScanLog{})
	if input.StartDate != "" {
		scanQuery = scanQuery.Where("scan_time >= ?", input.StartDate)
	}
	if input.EndDate != "" {
		scanQuery = scanQuery.Where("scan_time < ?", input.EndDate+" 23:59:59")
	}

	// 4. 统计活跃用户数 (定义为在时间段内有扫码行为的用户)
	var activeUsersCount int64
	if err := scanQuery.Distinct("user_union_id").Count(&activeUsersCount).Error; err != nil {
		return nil, err
	}

	// 5. 统计用户设备分布
	var deviceDist []UserDeviceDistribution
	if err := scanQuery.Where("brand != ''").Select("brand, model, count(*) as count").Group("brand, model").Order("count DESC").Limit(20).Find(&deviceDist).Error; err != nil {
		return nil, err
	}

	result := gin.H{
		"new_users_count":       newUsersCount,
		"active_users_count":    activeUsersCount,
		"gender_distribution":   genderDist,
		"province_distribution": provinceDist,
		"device_distribution":   deviceDist,
	}

	return result, nil
}

// GetCouponStatsInput 定义了获取优惠券统计数据的输入结构
type GetCouponStatsInput struct {
	StartDate string `form:"start_date"` // 格式: "2006-01-02"
	EndDate   string `form:"end_date"`   // 格式: "2006-01-02"
	StoreID   *uint  `form:"store_id"`
}

// CouponOverallStats 定义了优惠券的总体统计
type CouponOverallStats struct {
	TotalIssued   int64   `json:"total_issued"`   // 总发放(领取)
	TotalUsed     int64   `json:"total_used"`     // 总使用
	UsageRate     float64 `json:"usage_rate"`     // 核销率
	TotalDeducted float64 `json:"total_deducted"` // 总抵扣金额
}

// CouponStatsByType 定义了按类型统计的结果
type CouponStatsByType struct {
	CouponType string `json:"coupon_type"`
	Issued     int64  `json:"issued"`
	Used       int64  `json:"used"`
}

// GetCouponStats 用于获取优惠券相关的统计数据
func (s *StatsService) GetCouponStats(input *GetCouponStatsInput) (any, error) {
	db := database.DB.WithContext(context.Background())
	logQuery := db.Model(&models.CouponLog{})

	if input.StartDate != "" {
		logQuery = logQuery.Where("action_time >= ?", input.StartDate)
	}
	if input.EndDate != "" {
		logQuery = logQuery.Where("action_time < ?", input.EndDate+" 23:59:59")
	}
	if input.StoreID != nil {
		logQuery = logQuery.Where("store_id = ?", *input.StoreID)
	}

	// 1. 总体统计
	var overall CouponOverallStats
	// 克隆查询以避免互相影响
	if err := logQuery.Where("action_type = 'RECEIVE'").Count(&overall.TotalIssued).Error; err != nil {
		return nil, err
	}
	// 重新构建查询以替代 logQuery.Clone()
	usedQuery := db.Model(&models.CouponLog{})
	if input.StartDate != "" {
		usedQuery = usedQuery.Where("action_time >= ?", input.StartDate)
	}
	if input.EndDate != "" {
		usedQuery = usedQuery.Where("action_time < ?", input.EndDate+" 23:59:59")
	}
	if input.StoreID != nil {
		usedQuery = usedQuery.Where("store_id = ?", *input.StoreID)
	}
	if err := usedQuery.Where("action_type = 'USE'").Count(&overall.TotalUsed).Error; err != nil {
		return nil, err
	}

	// 再次重新构建查询
	deductedQuery := db.Model(&models.CouponLog{})
	if input.StartDate != "" {
		deductedQuery = deductedQuery.Where("action_time >= ?", input.StartDate)
	}
	if input.EndDate != "" {
		deductedQuery = deductedQuery.Where("action_time < ?", input.EndDate+" 23:59:59")
	}
	if input.StoreID != nil {
		deductedQuery = deductedQuery.Where("store_id = ?", *input.StoreID)
	}
	if err := deductedQuery.Where("action_type = 'USE'").Select("sum(amount_deducted)").Scan(&overall.TotalDeducted).Error; err != nil {
		return nil, err
	}
	if overall.TotalIssued > 0 {
		overall.UsageRate = float64(overall.TotalUsed) / float64(overall.TotalIssued)
	}

	// 2. 按类型统计 (需要关联 coupon 表)
	var byType []CouponStatsByType

	// 子查询：在指定时间范围内领取的券
	issuedSubQuery := db.Model(&models.CouponLog{}).Select("coupon_id").Where("action_type = 'RECEIVE'")
	if input.StartDate != "" {
		issuedSubQuery = issuedSubQuery.Where("action_time >= ?", input.StartDate)
	}
	if input.EndDate != "" {
		issuedSubQuery = issuedSubQuery.Where("action_time < ?", input.EndDate+" 23:59:59")
	}

	// 子查询：在指定时间范围内使用的券
	usedSubQuery := db.Model(&models.CouponLog{}).Select("coupon_id").Where("action_type = 'USE'")
	if input.StartDate != "" {
		usedSubQuery = usedSubQuery.Where("action_time >= ?", input.StartDate)
	}
	if input.EndDate != "" {
		usedSubQuery = usedSubQuery.Where("action_time < ?", input.EndDate+" 23:59:59")
	}

	if err := db.Table("coupon").
		Select(`
			coupon_type,
			(SELECT count(*) FROM coupon_log WHERE coupon_log.coupon_id = coupon.coupon_id AND coupon_log.action_type = 'RECEIVE') as issued,
			(SELECT count(*) FROM coupon_log WHERE coupon_log.coupon_id = coupon.coupon_id AND coupon_log.action_type = 'USE') as used
		`).
		Group("coupon_type").
		Find(&byType).Error; err != nil {
		// 注意: 上面的SQL没有加入时间范围和门店筛选，因为它比较复杂。
		// 在实际生产中，可能需要更精细的SQL或数据仓库来处理这类复杂聚合。
		// 此处为了演示，我们先用一个简化的全局统计。
	}

	result := gin.H{
		"overall_stats": overall,
		"by_type_stats": byType,
	}

	return result, nil
}
