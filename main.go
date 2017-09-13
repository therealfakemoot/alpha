package main

import (
	"fmt"
	dgo "github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
	conf "github.com/therealfakemoot/alpha/src/conf"
	exc "github.com/therealfakemoot/alpha/src/exchange"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func messageCreate(s *dgo.Session, m *dgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "!exchange") {
		args := strings.Split(m.Content, " ")
		if len(args) != 3 {
			s.ChannelMessageSend(m.ChannelID, "Doing it wrong")
		}

		from := strings.ToUpper(args[1])
		to := strings.ToUpper(args[2])

		apiResp := exc.HistoMinute(0, from, to)

		s.ChannelMessageSendEmbed(m.ChannelID, apiResp.Embed(false))

	}

}

func guildCreate(s *dgo.Session, event *dgo.GuildCreate) {

	if event.Guild.Unavailable {
		return
	}

	for _, channel := range event.Guild.Channels {
		if channel.ID == event.Guild.ID {
			_, _ = s.ChannelMessageSend(channel.ID, "Alpha, reporting for duty.")
			return
		}
	}
}

func runBot(v *viper.Viper) {
	d, err := dgo.New("Bot " + v.GetString("TOKEN_DISCORD"))

	d.LogLevel = dgo.LogDebug

	d.AddHandler(messageCreate)
	d.AddHandler(guildCreate)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	err = d.Open()

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	d.Close()

}

func main() {
	v := conf.LoadConf()
	v.ReadInConfig()
	runBot(v)
}
