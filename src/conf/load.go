package conf

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func setup() *viper.Viper {
	v := viper.New()

	v.SetConfigName("bot-conf")

	v.AddConfigPath("/etc/alpha/")
	v.AddConfigPath("$HOME/")
	v.AddConfigPath(".")

	v.SetEnvPrefix("ALPHA")
	v.AutomaticEnv()

	v.SetDefault("debug", false)
	v.SetDefault("db_dir", "$HOME/alpha/db/")

	v.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("[Config file changed:", e.Name, "]")
	})

	return v
}

// LoadConf provides a standard bot viper struct
func LoadConf() *viper.Viper {
	return setup()
}
