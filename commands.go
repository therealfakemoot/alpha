package main

import (
	"errors"
	"fmt"
	"strings"

	dgo "github.com/bwmarrin/discordgo"
	trash "github.com/therealfakemoot/trash-talk"
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
type Command func(args []string, conf Conf, s *dgo.Session, e interface{}) error

// Route blah blah
func Route(input string, conf Conf, cmds map[string]Command, s *dgo.Session, e interface{}) error {
	switch e.(type) {
	case *dgo.MessageCreate:
		args := strings.Split(input, " ")

		if string(args[0][0]) == "!" {
			cmd, ok := cmds[args[0][1:]]
			if !ok {
				return ErrNoCmdFound
			}
			return cmd(args[1:], conf, s, e)
		}
	case *dgo.Ready:
	case *dgo.Connect:
	case *dgo.Resumed:
		s.UpdateStatus(0, conf.Status)
	default:
		return ErrUnexpectedEvent{event: e}
	}
	return nil
}

// Mock is a Command that makes fun of the last message a given user sent.
func Mock(args []string, conf Conf, s *dgo.Session, e interface{}) error {
	msgMap := conf.State["msgMap"].(map[string]*dgo.Message)
	m := e.(*dgo.MessageCreate)
	if len(m.Message.Mentions) == 0 {
		s.ChannelMessageSend(m.ChannelID, "You didn't mention anyone.")
		return ErrIncorrectArgs
	}

	target := m.Message.Mentions[0].ID
	targetMsg, ok := msgMap[target]
	if !ok {
		s.ChannelMessageSend(m.ChannelID, "They haven't said anything yet.")
		return ErrIncorrectArgs
	}

	s.ChannelMessageSend(m.ChannelID, trash.Mock(targetMsg.Content))
	return nil
}

func Complain(args []string, conf Conf, s *dgo.Session, e interface{}) error {
	m := e.(*dgo.MessageCreate)
	s.ChannelMessageSend(m.ChannelID, "Life is hard.")

	return nil
}
