package main

import (
	"strings"

	dgo "github.com/bwmarrin/discordgo"
)

var (
	SWEARS = []string{"fuck", "shit", "dick", "pussy", "ass"}
)

// ContainsAny searchs a given string for the presence of any of the provided substrings.
func ContainsAny(s string, terms []string) bool {
	for _, t := range terms {
		if strings.Contains(s, t) {
			return true
		}
	}
	return false
}

func Chatter(content string, conf Conf, session *dgo.Session, message *dgo.MessageCreate) {
	if ContainsAny(content, SWEARS) {
		s.ChannelMessageSend("You've got a dirty mouth, you need a spanking.")
	}
}
