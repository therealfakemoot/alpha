package main

import (
	"fmt"
	dgo "github.com/bwmarrin/discordgo"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"github.com/therealfakemoot/alpha/routes"
	"log"
	"os/user"
)

// LoadConfig instantiates a Viper object with config info required for the bot to work.
func LoadConfig() *viper.Viper {
	v := viper.New()

	v.SetEnvPrefix("ALPHA")
	v.AutomaticEnv()
	v.SetConfigName(".alpha")
	v.AddConfigPath("/etc/alpha")

	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	v.AddConfigPath(user.HomeDir)

	err = v.ReadInConfig()
	if err != nil {
		fmt.Errorf("Fatal error config file: %s \n", err)
	}

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})

	return v
}

func main() {
	conf := LoadConfig()
	fmt.Printf("%+v\n", conf.AllSettings())

	token := "Bot " + conf.GetString("token")

	msgMap := make(map[string]*dgo.Message)
	conf.Set("msgMap", msgMap)

	s, err := dgo.New(token)
	if err != nil {
		log.Fatal(err)
	}

	r := routes.All()

	s.AddHandler(func(_ *dgo.Session, m *dgo.MessageCreate) {
		r.FindAndExecute(s, conf.GetString("prefix"), s.State.User.ID, m.Message)

		msgMap := conf.Get("msgMap").(map[string]*dgo.Message)
		msgMap[m.Author.ID] = m.Message
		conf.Set("msgMap", msgMap)
	})
	err = s.Open()

	if err != nil {
		log.Fatal(err)
	}

	<-make(chan struct{})
}
