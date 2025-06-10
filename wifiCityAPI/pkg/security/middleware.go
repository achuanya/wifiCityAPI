package security

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/wifiCityAPI/config"
)

const (
	// DecryptedBodyKey 是存储在 gin.Context 中的解密后请求体的键
	DecryptedBodyKey = "decryptedBody"
)

// Authenticate 是一个Gin中间件，用于验证请求签名和解密请求体
func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取和验证时间戳
		timestampStr := c.GetHeader("X-Timestamp")
		if timestampStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "缺少时间戳头 (X-Timestamp)"})
			return
		}
		timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "无效的时间戳格式"})
			return
		}
		if time.Now().Unix()-timestamp > int64(config.Cfg.Security.TimestampWindow.Seconds()) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "时间戳已过期，请检查设备时间"})
			return
		}

		// 2. 构建待签名字符串
		// 对查询参数按key进行字典排序
		queryKeys := make([]string, 0, len(c.Request.URL.Query()))
		for k := range c.Request.URL.Query() {
			queryKeys = append(queryKeys, k)
		}
		sort.Strings(queryKeys)
		var sortedQuery strings.Builder
		for i, k := range queryKeys {
			if i > 0 {
				sortedQuery.WriteString("&")
			}
			sortedQuery.WriteString(k)
			sortedQuery.WriteString("=")
			sortedQuery.WriteString(c.Request.URL.Query().Get(k))
		}

		var bodyStr string
		// 如果是 POST, PUT, DELETE 请求，则读取请求体
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "DELETE" {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "无法读取请求体"})
				return
			}
			// 必须将读取的body再写回去，因为 c.Request.Body 是一个只能读取一次的流
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			bodyStr = string(bodyBytes)
		}

		// 构造签名原文
		var stringToSign strings.Builder
		stringToSign.WriteString(timestampStr)
		stringToSign.WriteString(c.Request.Method)
		stringToSign.WriteString(c.Request.URL.Path)
		if sortedQuery.Len() > 0 {
			stringToSign.WriteString("?")
			stringToSign.WriteString(sortedQuery.String())
		}
		if bodyStr != "" {
			stringToSign.WriteString(" ")
			stringToSign.WriteString(bodyStr)
		}

		// 3. 验证签名
		signature := c.GetHeader("X-Signature")
		if signature == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "缺少签名头 (X-Signature)"})
			return
		}
		secretKey := []byte(config.Cfg.Security.APISecret)
		if !ValidateSignature(stringToSign.String(), signature, secretKey) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "签名验证失败"})
			return
		}

		// 4. 如果有请求体，则解密
		if bodyStr != "" {
			var encryptedRequest EncryptedData
			if err := json.Unmarshal([]byte(bodyStr), &encryptedRequest); err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "无效的加密请求体格式"})
				return
			}

			decryptedBody, err := Decrypt(&encryptedRequest, secretKey)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "请求体解密失败", "detail": err.Error()})
				return
			}
			// 将解密后的数据存入 context，并重置请求体，以便后续的 Bind 操作
			c.Set(DecryptedBodyKey, decryptedBody)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(decryptedBody))
		}

		c.Next()
	}
}

// ErrorResponse 定义了标准的错误响应结构体
type ErrorResponse struct {
	Error string `json:"error"`
}

// sendEncryptedError 是一个内部辅助函数，用于发送一个加密后的标准错误响应。
// 这样做可以确保即便是错误信息也不会明文传输。
func sendEncryptedError(c *gin.Context, status int, message string) {
	// 1. 从配置中获取API密钥
	key := []byte(config.Cfg.Security.APISecret)

	// 2. 创建标准错误结构体并序列化为JSON
	errorResponse := ErrorResponse{Error: message}
	jsonBytes, err := json.Marshal(errorResponse)
	if err != nil {
		// 如果连序列化标准错误信息都失败了，记录日志并返回一个未加密的、最基础的错误。
		// 这是极端情况，通常不应该发生。
		log.Printf("CRITICAL: 无法序列化标准错误响应: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误：无法生成错误信息"})
		return
	}

	// 3. 加密序列化后的错误信息
	encryptedData, err := Encrypt(jsonBytes, key)
	if err != nil {
		// 如果加密失败，记录日志并返回一个未加密的、最基础的错误。
		log.Printf("CRITICAL: 无法加密标准错误响应: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误：无法加密错误信息"})
		return
	}

	// 4. 发送加密后的错误信息
	c.JSON(status, encryptedData)
}

// SendEncryptedResponse 是一个统一的响应发送函数。
// 它会将任何给定的数据（无论是成功的结果还是错误信息）加密后发送给客户端。
func SendEncryptedResponse(c *gin.Context, status int, data any) {
	// 1. 从配置中获取API密钥
	key := []byte(config.Cfg.Security.APISecret)

	// 2. 将传入的数据序列化为JSON
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		log.Printf("ERROR: 无法序列化响应数据: %v", err)
		// 如果序列化失败，则发送一个加密的通用服务器错误。
		sendEncryptedError(c, http.StatusInternalServerError, "服务器无法处理您的响应。")
		return
	}

	// 3. 加密JSON数据
	encryptedData, err := Encrypt(jsonBytes, key)
	if err != nil {
		log.Printf("ERROR: 无法加密响应数据: %v", err)
		// 如果加密失败，也发送一个加密的通用服务器错误。
		sendEncryptedError(c, http.StatusInternalServerError, "服务器无法保护您的响应。")
		return
	}

	// 4. 发送最终的加密数据
	c.JSON(status, encryptedData)
}
