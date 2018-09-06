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

	s.AddHandler(func(s *dgo.Session, r *dgo.Ready) {
		s.UpdateStatus(0, conf.Status)
	})

	s.AddHandler(func(s *dgo.Session, r *dgo.Connect) {
		s.UpdateStatus(0, conf.Status)
	})

	s.AddHandler(func(s *dgo.Session, r *dgo.Resumed) {
		s.UpdateStatus(0, conf.Status)
	})

	s.AddHandler(func(s *dgo.Session, m *dgo.MessageCreate) {
		err = Route(m.Content, conf, cmds, s, m)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Error happened")
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
