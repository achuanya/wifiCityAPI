package router

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/gin-gonic/gin/wifiCityAPI/internal/api/v1"
	"github.com/gin-gonic/gin/wifiCityAPI/pkg/security"
)

// SetupRouter 配置并返回一个 Gin 引擎
func SetupRouter() *gin.Engine {
	// 使用 gin.Default() 创建一个带有默认中间件（Logger 和 Recovery）的路由引擎
	r := gin.Default()

	// 可以在这里添加一些全局中间件，例如 CORS 等
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
			users.POST("/", userHandler.CreateOrUpdateUser) // 创建或更新用户档案
			users.GET("/", userHandler.GetUser)             // 根据 UnionID, OpenID 或手机号获取用户详情
		}

		// 扫码日志相关路由
		scanLogs := apiV1.Group("/scan-logs")
		{
			scanLogs.POST("/", scanLogHandler.CreateScanLog)                   // 记录用户扫码连接日志
			scanLogs.GET("/", scanLogHandler.GetScanLogs)                      // 查询扫码日志列表
			scanLogs.PUT("/:logId/result", scanLogHandler.UpdateScanLogResult) // 更新扫码日志连接结果
			scanLogs.GET("/stats/daily-count/:storeId", scanLogHandler.GetDailyScanCountByStore)
		}

		// 优惠券路由
		coupons := apiV1.Group("/coupons")
		{
			coupons.POST("/", couponHandler.CreateCoupon)
			coupons.POST("/batch", couponHandler.CreateBatchCoupons)
			coupons.GET("/available-for-user", couponHandler.GetAvailableCouponsForUser)
			coupons.GET("/:id", couponHandler.GetCoupon)
			coupons.GET("/", couponHandler.GetCoupons)
			coupons.PUT("/:id", couponHandler.UpdateCoupon)
			coupons.DELETE("/:id", couponHandler.DeleteCoupon)
		}

		// 优惠券日志路由
		couponLogs := apiV1.Group("/coupon-logs")
		{
			couponLogs.POST("/", couponLogHandler.CreateCouponLog)
			couponLogs.GET("/", couponLogHandler.GetCouponLogs)
		}

		// 数据统计与报表路由
		stats := apiV1.Group("/stats")
		{
			stats.GET("/stores", statsHandler.GetStoreStats)
			stats.GET("/wifi-usage", statsHandler.GetWifiUsageStats)
			stats.GET("/user-behavior", statsHandler.GetUserBehaviorStats)
			stats.GET("/coupons", statsHandler.GetCouponStats)
		}

		// WIFI配置路由
		wifi := apiV1.Group("/wifi-configs")
		{
			wifi.POST("/", wifiHandler.CreateWifiConfig)
			wifi.POST("/batch", wifiHandler.CreateBatchWifiConfigs)
			wifi.GET("/:id", wifiHandler.GetWifiConfig)
			wifi.GET("/store/:storeId", wifiHandler.GetWifiConfigsByStore)
			wifi.PUT("/:id", wifiHandler.UpdateWifiConfig)
		}

		// 您可以在此继续添加其他资源的路由
	}

	return r
}
