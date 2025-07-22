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

// Config 结构体定义，通过结构体标签 mapstructure 映射配置文件中的键名到结构体字段。
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

// 打印配置信息
func printConfig(cfg Config) {
	// 打印配置信息
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

// 用反射(reflect)来简化配置集比较代码，避免手动编写每个字段的比较逻辑。
func compareConfigs(old, new Config) {

	fmt.Println("==== 配置变更详情 ====")

	oldVal := reflect.ValueOf(old)
	newVal := reflect.ValueOf(new)
	typ := oldVal.Type()

	for i := 0; i < oldVal.NumField(); i++ {
		field := typ.Field(i)
		oldField := oldVal.Field(i)
		newField := newVal.Field(i)

		// 递归比较结构体字段
		if field.Type.Kind() == reflect.Struct {
			compareStructFields(field.Name, oldField, newField)
			continue
		}

		// 比较基本类型字段
		if !reflect.DeepEqual(oldField.Interface(), newField.Interface()) {
			fmt.Printf("%s 变更: %v → %v\n",
				field.Tag.Get("mapstructure"),
				oldField.Interface(),
				newField.Interface())
		}
	}
}

// 用于递归比较嵌套结构体字段的辅助函数。
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

// 计算文件哈希值，用于检测配置文件的变更。
func calculateFileHash(filename string) ([]byte, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	hash := sha256.Sum256(content)
	return hash[:], nil
}

func main() {
	// 1. 初始化 Viper 实例
	v := viper.New()          // 初始化 Viper 实例
	v.SetConfigType("yaml")   // 设置配置文件类型为 YAML
	v.SetConfigName("config") // 设置配置文件名（不含路径、不含后缀）
	// 在当前目录（.）以及 conf 子目录中搜索配置文件 config.yaml
	v.AddConfigPath("./conf")
	v.AddConfigPath(".")

	// 2. 首先设置默认值
	v.SetDefault("db.DB_PORT", 3306)

	// 3. 初始加载配置
	var currentConfig Config
	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
	if err := v.Unmarshal(&currentConfig); err != nil {
		log.Fatalf("Unable to decode config into struct: %v", err)
	}

	// 3. 设置配置变更回调（在初始加载之后）
	var (
		lastConfigHash []byte
	)
	v.OnConfigChange(func(e fsnotify.Event) {

		newHash, err := calculateFileHash(e.Name) // 计算配置文件的哈希值
		if err != nil {
			log.Printf("计算文件哈希失败: %v", err)
			return
		}

		// 只处理 fsnotify.Write（文件写入），并且文件 hash 值改变，才处理配置变更事件。
		if e.Has(fsnotify.Write) && string(newHash) != string(lastConfigHash) {

			lastConfigHash = newHash
			fmt.Println("检测到有效配置变更:", e.Name, "操作类型:", e.Op)

			var newConfig Config
			if err := v.ReadInConfig(); err != nil {
				log.Printf("重新加载配置失败: %v", err)
				return
			}
			if err := v.Unmarshal(&newConfig); err != nil {
				log.Printf("解析新配置失败: %v", err)
				return
			}

			// 比较配置差异
			compareConfigs(currentConfig, newConfig)

			// 更新当前配置
			currentConfig = newConfig
			fmt.Println("==================== 重新载入配置 ====================")
			printConfig(currentConfig)
		}
	})
	// 4. 监控配置文件变更
	v.WatchConfig() // 开始监控配置文件的更改

	// 5. 打印初始配置
	printConfig(currentConfig)

	select {} // 阻塞主线程，防止退出，等待配置变更事件
}
