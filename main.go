package main

import (
	"log"

	dgo "github.com/bwmarrin/discordgo"
)

func updateState(conf Conf, e interface{}) {
	switch e.(type) {
	case *dgo.MessageCreate:
		m := e.(*dgo.MessageCreate)
		msgMap := conf.State["msgMap"].(map[string]*dgo.Message)
		msgMap[m.Author.ID] = m.Message
		conf.State["msgMap"] = msgMap
	case *dgo.Ready:
	default:
		return
	}
}

func main() {
	conf := LoadConfig()
	conf.State = make(map[string]interface{})

	token := "Bot " + conf.Token
	msgMap := make(map[string]*dgo.Message)
	conf.State["msgMap"] = msgMap

	s, err := dgo.New(token)
	if err != nil {
		log.Fatal(err)
	}

	cmds := make(map[string]Command)
	cmds["mock"] = Mock
	cmds["bitch"] = Complain

	s.AddHandler(func(s *dgo.Session, r *dgo.Ready) {
		Route("", conf, cmds, s, r)
	})

	s.AddHandler(func(s *dgo.Session, r *dgo.Connect) {
		Route("", conf, cmds, s, r)
	})

	s.AddHandler(func(s *dgo.Session, r *dgo.Resumed) {
		Route("", conf, cmds, s, r)
	})

	s.AddHandler(func(s *dgo.Session, m *dgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}
		err = Route(m.Content, conf, cmds, s, m)
		if err == ErrNoCmdGiven {
			Chatter(m.Content, conf, s, m)
		} else if err != nil {
			s.ChannelMessageSend(m.ChannelID, err.Error())
		}

		updateState(conf, m)

	})

	err = s.Open()
	defer s.Close()

	if err != nil {
		log.Fatal(err)
	}

	<-make(chan struct{})
}
