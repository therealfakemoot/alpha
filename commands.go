package main

import (
	"errors"
	"strings"

	dgo "github.com/bwmarrin/discordgo"
	trash "github.com/therealfakemoot/trash-talk"
)

var (
	// IncorrectArgs is a custom error type so I can eventually gracefully handle specific errors I guess.
	IncorrectArgs = errors.New("incorrect arguments supplied")
)

// Command blah blah
type Command func(args []string, conf Conf, s *dgo.Session, e interface{}) error

// Route blah blah
func Route(input string, conf Conf, cmds map[string]Command, s *dgo.Session, e interface{}) error {
	args := strings.Split(input, " ")

	if string(args[0][0]) == "!" {
		cmd, ok := cmds[args[0][1:]]
		if !ok {
			return IncorrectArgs
		}
		return cmd(args[1:], conf, s, e)
	}
	return nil
}

// Mock is a Command that makes fun of the last message a given user sent.
func Mock(args []string, conf Conf, s *dgo.Session, e interface{}) error {
	msgMap := conf.State["msgMap"].(map[string]*dgo.Message)
	m := e.(*dgo.MessageCreate)
	if len(m.Message.Mentions) == 0 {
		s.ChannelMessageSend(m.ChannelID, "You didn't mention anyone.")
		return IncorrectArgs
	}

	target := m.Message.Mentions[0].ID
	targetMsg, ok := msgMap[target]
	if !ok {
		s.ChannelMessageSend(m.ChannelID, "They haven't said anything yet.")
		return IncorrectArgs
	}

	s.ChannelMessageSend(m.ChannelID, trash.Mock(targetMsg.Content))
	return nil
}
