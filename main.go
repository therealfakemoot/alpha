package main

import (
	"log"

	dgo "github.com/bwmarrin/discordgo"
	bolt "go.etcd.io/bbolt"
)

func main() {
	db, err := bolt.Open("state.db", 0600, nil)
	if err != nil {
		log.Fatalf("unable to open state db: %s", err)
	}
	defer db.Close()

	LoadConfig(db)

	var token string

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("core"))
		token = string(b.Get([]byte("token")))

		return nil
	})

	s, err := dgo.New("Bot " + token)
	if err != nil {
		log.Fatal(err)
	}

	cmds := make(map[string]Command)
	cmds["mock"] = Mock
	cmds["bitch"] = Complain

	s.AddHandler(func(s *dgo.Session, gc *dgo.GuildCreate) {
		db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte(gc.ID))
			if err != nil {
				return err
			}

			_, err = b.CreateBucketIfNotExists([]byte("messages"))

			if err != nil {
				return err
			}

			return nil
		})
	})

	s.AddHandler(func(s *dgo.Session, r *dgo.Ready) {
		Route("", db, cmds, s, r)
	})

	s.AddHandler(func(s *dgo.Session, r *dgo.Connect) {
		Route("", db, cmds, s, r)
	})

	s.AddHandler(func(s *dgo.Session, r *dgo.Resumed) {
		Route("", db, cmds, s, r)
	})

	s.AddHandler(func(s *dgo.Session, m *dgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}
		err = Route(m.Content, db, cmds, s, m)
		if err == ErrNoCmdGiven {
			Chatter(m.Content, db, s, m)
		} else if err != nil {
			s.ChannelMessageSend(m.ChannelID, err.Error())
		}

		updateState(db, s, m)

	})

	err = s.Open()
	defer s.Close()

	if err != nil {
		log.Fatal(err)
	}

	<-make(chan struct{})
}
