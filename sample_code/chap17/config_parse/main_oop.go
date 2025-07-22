package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Config 结构体定义
type Config struct {
	DB   DBConfig   `mapstructure:"db"`
	SQL  SQLConfig  `mapstructure:"sql"`
	XLSX XLSXConfig `mapstructure:"xlsx"`
}

type DBConfig struct {
	Host     string `mapstructure:"DB_HOST"`
	Port     int    `mapstructure:"DB_PORT"`
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
}

type SQLConfig struct {
	Cmd string `mapstructure:"SQL_CMD"`
}

type XLSXConfig struct {
	RptDir       string `mapstructure:"XLSX_RPT_DIR"`
	RptPrefix    string `mapstructure:"XLSX_RPT_PREFIX"`
	HeaderFill   string `mapstructure:"HEADER_FILL"`
	DataEvenFill string `mapstructure:"DATA_EVEN_FILL"`
	AutoOpen     bool   `mapstructure:"AUTO_OPEN"`
	AutoOpenSecs int    `mapstructure:"AUTO_OPEN_SECS"`
}

// ConfigManager 负责配置管理
type ConfigManager struct {
	viper          *viper.Viper
	currentConfig  Config
	lastConfigHash []byte
}

// NewConfigManager 创建新的配置管理器
func NewConfigManager() *ConfigManager {
	cm := &ConfigManager{
		viper: viper.New(),
	}
	cm.setupViper()
	return cm
}

// 初始化Viper配置
func (cm *ConfigManager) setupViper() {
	cm.viper.SetConfigType("yaml")
	cm.viper.SetConfigName("config")
	cm.viper.AddConfigPath("./conf")
	cm.viper.AddConfigPath(".")
	cm.viper.SetDefault("db.DB_PORT", 3306)
}

// LoadConfig 加载初始配置
func (cm *ConfigManager) LoadConfig() error {
	if err := cm.viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading config file: %v", err)
	}
	if err := cm.viper.Unmarshal(&cm.currentConfig); err != nil {
		return fmt.Errorf("unable to decode config into struct: %v", err)
	}

	// 计算初始哈希
	hash, err := calculateFileHash(cm.viper.ConfigFileUsed())
	if err == nil {
		cm.lastConfigHash = hash
	}
	return nil
}

// WatchConfigChanges 监控配置变更
func (cm *ConfigManager) WatchConfigChanges() {
	cm.viper.OnConfigChange(cm.handleConfigChange)
	cm.viper.WatchConfig()
}

// 处理配置变更事件
func (cm *ConfigManager) handleConfigChange(e fsnotify.Event) {
	newHash, err := calculateFileHash(e.Name)
	if err != nil {
		log.Printf("计算文件哈希失败: %v", err)
		return
	}

	if e.Has(fsnotify.Write) && string(newHash) != string(cm.lastConfigHash) {
		cm.lastConfigHash = newHash
		fmt.Println("检测到有效配置变更:", e.Name)

		var newConfig Config
		if err := cm.viper.ReadInConfig(); err != nil {
			log.Printf("重新加载配置失败: %v", err)
			return
		}
		if err := cm.viper.Unmarshal(&newConfig); err != nil {
			log.Printf("解析新配置失败: %v", err)
			return
		}

		cm.compareAndUpdateConfig(newConfig)
	}
}

// 比较并更新配置
func (cm *ConfigManager) compareAndUpdateConfig(newConfig Config) {
	fmt.Println("==== 配置变更详情 ====")
	compareConfigs(cm.currentConfig, newConfig)

	cm.currentConfig = newConfig
	fmt.Println("==================== 重新载入配置 ====================")
	cm.PrintCurrentConfig()
}

// PrintCurrentConfig 打印当前配置
func (cm *ConfigManager) PrintCurrentConfig() {
	printConfig(cm.currentConfig)
}

// 计算文件哈希
func calculateFileHash(filename string) ([]byte, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	hash := sha256.Sum256(content)
	return hash[:], nil
}

// 打印配置信息（可作为独立函数或移动到ConfigManager中）
func printConfig(cfg Config) {
	fmt.Println("===== 数据库配置 =====")
	fmt.Printf("主机: %s\n", cfg.DB.Host)
	fmt.Printf("端口: %d\n", cfg.DB.Port)
	fmt.Printf("用户名: %s\n", cfg.DB.User)
	fmt.Printf("密码: %s\n", cfg.DB.Password)

	fmt.Println("\n====== SQL 配置 ======")
	fmt.Printf("SQL 命令: \n%s\n", cfg.SQL.Cmd)

	fmt.Println("\n=== Excel 报表配置 ===")
	fmt.Printf("报表目录: %s\n", cfg.XLSX.RptDir)
	fmt.Printf("报表前缀: %s\n", cfg.XLSX.RptPrefix)
	fmt.Printf("表头填充色: %s\n", cfg.XLSX.HeaderFill)
	fmt.Printf("偶数行填充色: %s\n", cfg.XLSX.DataEvenFill)
	fmt.Printf("自动打开: %v\n", cfg.XLSX.AutoOpen)
	fmt.Printf("自动打开延迟(秒): %d\n", cfg.XLSX.AutoOpenSecs)
}

// 比较配置差异（可作为独立函数或移动到ConfigManager中）
func compareConfigs(old, new Config) {
	oldVal := reflect.ValueOf(old)
	newVal := reflect.ValueOf(new)
	typ := oldVal.Type()

	for i := 0; i < oldVal.NumField(); i++ {
		field := typ.Field(i)
		oldField := oldVal.Field(i)
		newField := newVal.Field(i)

		if field.Type.Kind() == reflect.Struct {
			compareStructFields(field.Name, oldField, newField)
			continue
		}

		if !reflect.DeepEqual(oldField.Interface(), newField.Interface()) {
			fmt.Printf("%s 变更: %v → %v\n",
				field.Tag.Get("mapstructure"),
				oldField.Interface(),
				newField.Interface())
		}
	}
}

// 比较结构体字段（可作为独立函数或移动到ConfigManager中）
func compareStructFields(prefix string, oldVal, newVal reflect.Value) {
	typ := oldVal.Type()

	for i := 0; i < oldVal.NumField(); i++ {
		field := typ.Field(i)
		oldField := oldVal.Field(i)
		newField := newVal.Field(i)

		fieldName := prefix + "." + field.Tag.Get("mapstructure")

		if !reflect.DeepEqual(oldField.Interface(), newField.Interface()) {
			fmt.Printf("%s 变更: %v → %v\n",
				fieldName,
				oldField.Interface(),
				newField.Interface())
		}
	}
}

func main() {
	// 创建配置管理器
	configManager := NewConfigManager()

	// 加载初始配置
	if err := configManager.LoadConfig(); err != nil {
		log.Fatal(err)
	}

	// 开始监控配置变更
	configManager.WatchConfigChanges()

	// 打印初始配置
	configManager.PrintCurrentConfig()

	// 阻塞主线程
	select {}
}
