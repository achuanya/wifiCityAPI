package router

import (
	v1 "app/internal/api/v1"
	"app/pkg/security"

	"github.com/gin-gonic/gin"
)

// SetupRouter 配置并返回一个 Gin 引擎
func SetupRouter() *gin.Engine {
	// 使用 gin.Default() 创建一个带有默认中间件（Logger 和 Recovery）的路由引擎
	r := gin.Default()

	// 添加全局中间件
	// 1. 域名检查中间件
	r.Use(security.DomainCheck())

	// 2. 安全头信息中间件
	r.Use(security.SetSecureHeaders())

	// 3. 可以添加其他全局中间件，例如 CORS 等
	// r.Use(middlewares.Cors())

	// API V1 路由组
	apiV1 := r.Group("/api/v1")

	// 对 /api/v1 路由组下所有接口启用签名和加密校验
	apiV1.Use(security.Authenticate())
	{
		// 统一在此处实例化所有 Handler
		storeHandler := v1.NewStoreHandler()
		wifiHandler := v1.NewWifiConfigHandler()
		userHandler := v1.NewUserProfileHandler()
		scanLogHandler := v1.NewScanLogHandler()
		couponHandler := v1.NewCouponHandler()
		couponLogHandler := v1.NewCouponLogHandler()
		statsHandler := v1.NewStatsHandler()

		// 门店相关路由
		stores := apiV1.Group("/stores")
		{
			stores.POST("/with-wifi", storeHandler.CreateStoreWithWifi)
			stores.POST("/", storeHandler.CreateStore)
			stores.GET("/:storeId", storeHandler.GetStore)
			stores.GET("/", storeHandler.GetStores)
			stores.PUT("/:storeId", storeHandler.UpdateStore)
			stores.DELETE("/:storeId", storeHandler.DeleteStore)
			// 细分更新接口
			stores.PATCH("/:storeId/status", storeHandler.UpdateStoreStatus)     // 更新门店状态
			stores.PATCH("/:storeId/phone", storeHandler.UpdateStorePhone)       // 更新门店电话
			stores.PATCH("/:storeId/location", storeHandler.UpdateStoreLocation) // 更新门店位置
			// 关联路由：查询门店下的WIFI
			stores.GET("/:storeId/wifis", wifiHandler.GetWifiConfigsByStore)
			// 关联路由：查询门店的每日扫码量
			stores.GET("/:storeId/scans/daily-count", scanLogHandler.GetDailyScanCountByStore)
		}

		// WIFI 配置相关路由
		wifis := apiV1.Group("/wifis")
		{
			wifis.POST("", wifiHandler.CreateWifiConfig)       // 新增WIFI配置
			wifis.GET("/:id", wifiHandler.GetWifiConfig)       // 查询单个WIFI配置详情
			wifis.PUT("/:id", wifiHandler.UpdateWifiConfig)    // 更新WIFI配置
			wifis.DELETE("/:id", wifiHandler.DeleteWifiConfig) // 删除WIFI配置
		}

		// 用户信息相关路由
		users := apiV1.Group("/users")
		{
			users.POST("/", userHandler.CreateOrUpdateUser)            // 创建或更新用户档案
			users.GET("/", userHandler.GetUser)                        // 根据 UnionID, OpenID 或手机号获取用户详情
			users.POST("/bind-phone", userHandler.BindPhoneNumber)     // 用户绑定手机号
			users.POST("/unbind-phone", userHandler.UnbindPhoneNumber) // 用户解绑手机号
			users.GET("/scan-history", userHandler.GetUserScanHistory) // 查询用户扫码门店历史
		}

		// 扫码日志相关路由
		scanLogs := apiV1.Group("/scan-logs")
		{
			scanLogs.POST("/", scanLogHandler.CreateScanLog)                   // 记录用户扫码连接日志
			scanLogs.GET("/", scanLogHandler.GetScanLogs)                      // 查询扫码日志列表
			scanLogs.PUT("/:logId/result", scanLogHandler.UpdateScanLogResult) // 更新扫码日志连接结果
			scanLogs.GET("/stats/daily-count/:storeId", scanLogHandler.GetDailyScanCountByStore)
			scanLogs.GET("/failed", scanLogHandler.GetFailedScanLogs) // 查询扫码连接失败日志
			scanLogs.GET("/user", scanLogHandler.GetUserScanLogs)     // 查询指定用户的扫码历史
		}

		// 优惠券路由
		coupons := apiV1.Group("/coupons")
		{
			coupons.POST("/", couponHandler.CreateCoupon)
			coupons.POST("/batch", couponHandler.CreateBatchCoupons)
			coupons.GET("/available-for-user", couponHandler.GetAvailableCouponsForUser)
			coupons.GET("/store", couponHandler.GetCouponsByStore) // 查询门店可用优惠券列表
			coupons.GET("/:id", couponHandler.GetCoupon)
			coupons.GET("/", couponHandler.GetCoupons)
			coupons.PUT("/:id", couponHandler.UpdateCoupon)
			coupons.DELETE("/:id", couponHandler.DeleteCoupon)
			// 细分的优惠券更新接口
			coupons.PATCH("/:id/validity", couponHandler.UpdateCouponValidity) // 更新有效期
			coupons.PATCH("/:id/limit", couponHandler.UpdateCouponLimit)       // 更新使用限制
			coupons.PATCH("/:id/quantity", couponHandler.UpdateCouponQuantity) // 更新发行量
			coupons.PATCH("/:id/store", couponHandler.UpdateCouponStore)       // 更新适用门店
			coupons.PATCH("/:id/status", couponHandler.UpdateCouponStatus)     // 更新优惠券状态
		}

		// 优惠券日志路由
		couponLogs := apiV1.Group("/coupon-logs")
		{
			couponLogs.POST("/", couponLogHandler.CreateCouponLog)
			couponLogs.GET("/", couponLogHandler.GetCouponLogs)
			couponLogs.GET("/claim", couponLogHandler.GetCouponClaimLogs) // 查询优惠券领取记录
			couponLogs.GET("/use", couponLogHandler.GetCouponUseLogs)     // 查询优惠券核销使用记录
		}

		// 数据统计与报表路由
		stats := apiV1.Group("/stats")
		{
			stats.GET("/stores", statsHandler.GetStoreStats)
			stats.GET("/wifi-usage", statsHandler.GetWifiUsageStats)
			stats.GET("/user-behavior", statsHandler.GetUserBehaviorStats)
			stats.GET("/coupons", statsHandler.GetCouponStats)
			stats.GET("/popular-wifi", statsHandler.GetPopularWifi)                    // 最受欢迎WIFI统计
			stats.GET("/scan-time-distribution", statsHandler.GetScanTimeDistribution) // 扫码时段分布统计
		}

		// WIFI配置路由
		wifi := apiV1.Group("/wifi-configs")
		{
			wifi.POST("/", wifiHandler.CreateWifiConfig)
			wifi.POST("/batch", wifiHandler.CreateBatchWifiConfigs)
			wifi.GET("/type", wifiHandler.GetWifiConfigsByStoreAndType) // 查询门店特定类型的WIFI配置
			wifi.GET("/:id", wifiHandler.GetWifiConfig)
			wifi.GET("/store/:storeId", wifiHandler.GetWifiConfigsByStore)
			wifi.PUT("/:id", wifiHandler.UpdateWifiConfig)
			wifi.DELETE("/:id", wifiHandler.DeleteWifiConfig)         // 删除单个WIFI配置
			wifi.DELETE("/batch", wifiHandler.DeleteBatchWifiConfigs) // 批量删除WIFI配置
		}

		// 您可以在此继续添加其他资源的路由
	}

	return r
}
