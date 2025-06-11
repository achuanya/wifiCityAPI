package security

import (
	"app/config"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// DomainCheck 是一个中间件，用于检查请求的域名是否符合预期
// 这有助于防止未授权的域名访问API
func DomainCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从配置中获取允许的域名
		allowedDomain := config.Cfg.Server.Domain

		// 如果未设置域名或处于开发模式，则跳过检查
		if allowedDomain == "" || gin.Mode() == gin.DebugMode {
			c.Next()
			return
		}

		// 获取请求的Host（可能包含端口号）
		host := c.Request.Host

		// 移除可能的端口号
		if strings.Contains(host, ":") {
			host = strings.Split(host, ":")[0]
		}

		// 检查域名是否匹配配置中的域名
		// 使用不区分大小写的比较
		if !strings.EqualFold(host, allowedDomain) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "访问被拒绝，无效的域名",
			})
			return
		}

		c.Next()
	}
}

// SetSecureHeaders 设置安全相关的HTTP头信息
func SetSecureHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 根据配置决定是否强制HTTPS
		if config.Cfg.Server.UseHTTPS {
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}

		// 设置其他安全头信息
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")

		c.Next()
	}
}
