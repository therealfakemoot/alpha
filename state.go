package main

import (
	"fmt"

	dgo "github.com/bwmarrin/discordgo"
	bolt "go.etcd.io/bbolt"
)

func updateState(db *bolt.DB, s *dgo.Session, e interface{}) {
	switch e.(type) {
	case *dgo.MessageCreate:
		m := e.(*dgo.MessageCreate)
		db.Update(func(tx *bolt.Tx) error {
			guildBucket := tx.Bucket([]byte(m.GuildID))
			userMessages := guildBucket.Bucket([]byte("messages"))

			err := userMessages.Put([]byte(m.Author.ID), []byte(m.Content))

			if err != nil {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("error adding message to db: %s", err))
			}
			return err
		})
	case *dgo.Ready:
	default:
		return
	}
}
