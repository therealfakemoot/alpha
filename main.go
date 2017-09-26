package main

import (
	"fmt"
	dgo "github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
	conf "github.com/therealfakemoot/alpha/src/conf"
	disc "github.com/therealfakemoot/alpha/src/discord"
	exc "github.com/therealfakemoot/alpha/src/exchange"
	tick "github.com/therealfakemoot/alpha/src/tick"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var lastMessage *dgo.Message

func messageCreate(s *dgo.Session, m *dgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	args := strings.Split(m.Content, " ")
	if strings.HasPrefix(m.Content, "!exchange") {
		if len(args) != 3 {
			disc.NewMessage("Doing it wrong.", s, disc.FiveSecondPolicy)
			return
		}

		lastMessage, _ = s.ChannelMessageSend(m.ChannelID, "Doing it wrong")
		var i = 3
		f := func(t *tick.Timer) {
			i--
			fmt.Println("TICK")
			if i == 0 {
				t.Done()
			}
		}

		c := func(t *tick.Timer) {
			s.ChannelMessageDelete(lastMessage.ChannelID, lastMessage.ID)
		}

		tick.NewTimer(3*time.Second, f, c)
		return
	}

	from := strings.ToUpper(args[1])
	to := strings.ToUpper(args[2])

	apiResp := exc.HistoMinute(0, from, to)
	apiEmbed := apiResp.Embed(false)
	lastPriceMessage, err := s.ChannelMessageSendEmbed(m.ChannelID, apiEmbed)
	if err != nil {

		fmt.Println(err)
		return
	}

	var i = 0

	tf := func(tt *tick.Timer) {
		if i > 4 {
			tt.Done()
			return
		}
		tsField := &dgo.MessageEmbedField{}
		tsField.Name = "Self destruct timer"
		tsField.Value = string(5 - i)
		tsField.Inline = false
		me := dgo.NewMessageEdit(lastPriceMessage.ChannelID, lastPriceMessage.ID)
		apiEmbed.Fields[2] = tsField
		me.SetEmbed(apiEmbed)
		i++
	}

	cf := func(to *tick.Timer) {
		s.ChannelMessageDelete(lastPriceMessage.ChannelID, lastPriceMessage.ID)
	}

	tick.NewTimer(5*time.Second, tf, cf)

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
