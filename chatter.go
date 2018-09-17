package main

import (
	"golang.org/x/time/rate"
	"strings"
	"time"

	dgo "github.com/bwmarrin/discordgo"
)

var (
	r      = rate.Every(time.Duration(time.Second * 15))
	L      = rate.NewLimiter(r, 1)
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

func Chatter(c string, conf Conf, s *dgo.Session, m *dgo.MessageCreate) {
	if L.Allow() && ContainsAny(c, SWEARS) {
		s.ChannelMessageSend(m.ChannelID, "You've got a dirty mouth, you need a spanking.")
	}
}
