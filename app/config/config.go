package config

import (
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

// 全局配置变量
var Cfg *Config

// Config 定义了应用的配置结构
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Security SecurityConfig `yaml:"security"`
}

// ServerConfig 定义了服务器相关的配置
type ServerConfig struct {
	Port     string `yaml:"port"`
	Domain   string `yaml:"domain"`
	UseHTTPS bool   `yaml:"use_https"`
}

// DatabaseConfig 定义了数据库连接配置
type DatabaseConfig struct {
	Master   DBSource   `yaml:"master"`
	Slaves   []DBSource `yaml:"slaves"`
	Settings DBSettings `yaml:"settings"`
}

// DBSource 定义了单个数据源的连接信息
type DBSource struct {
	DSN string `yaml:"dsn"`
}

// DBSettings 定义了连接池的配置
type DBSettings struct {
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
	ConnMaxIdleTime time.Duration `yaml:"conn_max_idle_time"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
}

// SecurityConfig 定义了安全相关的配置
type SecurityConfig struct {
	APISecret       string        `yaml:"api_secret"`
	TimestampWindow time.Duration `yaml:"timestamp_window"`
}

// init 在包被导入时自动执行，用于加载配置
func init() {
	// 在测试环境中运行时，可能不需要加载配置文件
	if os.Getenv("GIN_MODE") == "test" {
		// 为测试环境设置默认配置
		Cfg = &Config{
			Server: ServerConfig{
				Port:     "8080",
				Domain:   "localhost",
				UseHTTPS: false,
			},
			Database: DatabaseConfig{
				Master:   DBSource{DSN: "user:pass@tcp(127.0.0.1:3306)/test_db?charset=utf8mb4&parseTime=True&loc=Local"},
				Slaves:   []DBSource{{DSN: "user:pass@tcp(127.0.0.1:3306)/test_db?charset=utf8mb4&parseTime=True&loc=Local"}},
				Settings: DBSettings{MaxIdleConns: 1, MaxOpenConns: 2, ConnMaxIdleTime: time.Minute, ConnMaxLifetime: time.Hour},
			},
			Security: SecurityConfig{APISecret: "1234567890123456", TimestampWindow: 300 * time.Second},
		}
		return
	}

	// 尝试从环境变量获取配置文件路径
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config/config.yaml"
	}

	// 加载配置文件
	if err := LoadConfig(configPath); err != nil {
		log.Fatalf("无法加载配置文件: %v", err)
	}
}

// LoadConfig 从指定路径加载配置文件并解析到全局 Cfg 变量
func LoadConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return err
	}

	// 将解析后的配置赋值给全局变量
	Cfg = &config

	// 将秒转换为 time.Duration
	Cfg.Security.TimestampWindow = Cfg.Security.TimestampWindow * time.Second

	// 允许从环境变量覆盖域名配置
	if domain := os.Getenv("API_DOMAIN"); domain != "" {
		Cfg.Server.Domain = domain
	}

	// 允许从环境变量覆盖HTTPS配置
	if useHTTPS := os.Getenv("USE_HTTPS"); useHTTPS == "true" {
		Cfg.Server.UseHTTPS = true
	}

	return nil
}
