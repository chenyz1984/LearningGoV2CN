package main

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// Config 结构体定义
type Config struct {
	DB struct {
		Host     string `mapstructure:"DB_HOST"`
		Port     int    `mapstructure:"DB_PORT"`
		User     string `mapstructure:"DB_USER"`
		Password string `mapstructure:"DB_PASSWORD"`
	} `mapstructure:"db"`

	SQL struct {
		Cmd string `mapstructure:"SQL_CMD"`
	} `mapstructure:"sql"`

	XLSX struct {
		RptDir       string `mapstructure:"XLSX_RPT_DIR"`
		RptPrefix    string `mapstructure:"XLSX_RPT_PREFIX"`
		HeaderFill   string `mapstructure:"HEADER_FILL"`
		DataEvenFill string `mapstructure:"DATA_EVEN_FILL"`
		AutoOpen     bool   `mapstructure:"AUTO_OPEN"`
		AutoOpenSecs int    `mapstructure:"AUTO_OPEN_SECS"`
	} `mapstructure:"xlsx"`
}

func main() {
	v := viper.New()          // 初始化 Viper 实例
	v.SetConfigType("yaml")   // 设置配置文件类型为 YAML
	v.SetConfigName("config") // 设置配置文件名（不含路径、不含后缀）
	v.AddConfigPath("./conf")      // 设置配置文件搜索路径为当前目录
	v.AddConfigPath(".")      // 设置配置文件搜索路径为当前目录

	// 查找并读取配置文件
	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	// 将配置反序列化到结构体
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		log.Fatalf("Unable to decode config into struct: %v", err)
	}

	// 打印配置信息
	fmt.Println("=== 数据库配置 ===")
	fmt.Printf("主机: %s\n", cfg.DB.Host)
	fmt.Printf("端口: %d\n", cfg.DB.Port)
	fmt.Printf("用户名: %s\n", cfg.DB.User)
	fmt.Printf("密码: %s\n", cfg.DB.Password)

	fmt.Println("\n=== SQL 配置 ===")
	fmt.Printf("SQL 命令: \n%s\n", cfg.SQL.Cmd)

	fmt.Println("\n=== Excel 报表配置 ===")
	fmt.Printf("报表目录: %s\n", cfg.XLSX.RptDir)
	fmt.Printf("报表前缀: %s\n", cfg.XLSX.RptPrefix)
	fmt.Printf("表头填充色: %s\n", cfg.XLSX.HeaderFill)
	fmt.Printf("偶数行填充色: %s\n", cfg.XLSX.DataEvenFill)
	fmt.Printf("自动打开: %v\n", cfg.XLSX.AutoOpen)
	fmt.Printf("自动打开延迟(秒): %d\n", cfg.XLSX.AutoOpenSecs)
}
