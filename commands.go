package main

import (
	"errors"
	"fmt"
	"strings"

	dgo "github.com/bwmarrin/discordgo"
	trash "github.com/therealfakemoot/trash-talk"
	bolt "go.etcd.io/bbolt"
)

// ErrUnexpectedEvent is thrown when discordgo gives me an unexpected or unhandled event.
type ErrUnexpectedEvent struct {
	event interface{}
}

func (e ErrUnexpectedEvent) Error() string {
	return fmt.Sprintf("%+v", e.event)
}

var (
	// ErrIncorrectArgs is a custom error type so I can eventually gracefully handle specific errors I guess.
	ErrIncorrectArgs = errors.New("incorrect arguments supplied")
	// ErrNoCmdFound indicates the cmds map doesn't have a matching key.
	ErrNoCmdFound = errors.New("no matching command found")
	// ErrNoCmdGiven indicates the message is not attempting to execute a command.
	ErrNoCmdGiven = errors.New("no command requested")
)

// Command blah blah
type Command func(args []string, db *bolt.DB, s *dgo.Session, e interface{}) error

// Route blah blah
func Route(input string, db *bolt.DB, cmds map[string]Command, s *dgo.Session, e interface{}) error {
	switch e.(type) {
	case *dgo.MessageCreate:
		args := strings.Split(input, " ")

		var prefix string
		db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("core"))
			prefix = string(b.Get([]byte("prefix")))

			return nil
		})

		if string(args[0][0]) == prefix {
			cmd, ok := cmds[args[0][1:]]
			if !ok {
				return ErrNoCmdFound
			}
			return cmd(args[1:], db, s, e)
		}

		return ErrNoCmdGiven
	case *dgo.Ready, *dgo.Connect, *dgo.Resumed:
		var status string
		db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("core"))
			status = string(b.Get([]byte("status")))

			return nil
		})
		err := s.UpdateStatus(0, status)
		if err != nil {
			// These events aren't associated with a channel so I've got to dump errors to stdout
			fmt.Printf("error updating status: %s", err)
		}
	default:
		return ErrUnexpectedEvent{event: e}
	}
	return nil
}

// Mock is a Command that makes fun of the last message a given user sent.
func Mock(args []string, db *bolt.DB, s *dgo.Session, e interface{}) error {
	m := e.(*dgo.MessageCreate)
	if len(m.Message.Mentions) == 0 {
		s.ChannelMessageSend(m.ChannelID, "You didn't mention anyone.")
		return ErrIncorrectArgs
	}

	target := m.Message.Mentions[0].ID

	var targetMsg string
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(m.GuildID))
		messages := b.Bucket([]byte("messages"))

		targetMsg = string(messages.Get([]byte(target)))

		return nil
	})

	_, err := s.ChannelMessageSend(m.ChannelID, trash.Mock(targetMsg))

	return err
}

func Complain(args []string, db *bolt.DB, s *dgo.Session, e interface{}) error {
	m := e.(*dgo.MessageCreate)
	s.ChannelMessageSend(m.ChannelID, "Life is hard.")

	return nil
}
