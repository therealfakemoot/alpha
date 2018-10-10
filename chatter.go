package main

import (
	"math/rand"
	"strings"
	"time"

	"golang.org/x/time/rate"

	dgo "github.com/bwmarrin/discordgo"
)

type ErrNoPromptMatch struct{}

func (e ErrNoPromptMatch) Error() string {
	return "No matching prompt text."
}

type ResponseMapping struct {
	Terms     []string
	Responses []string
}

func (r ResponseMapping) Reply(t string) (string, error) {
	if containsAny(t, r.Terms) {
		return r.Responses[rand.Intn(len(r.Responses))], nil
	}
	return "", ErrNoPromptMatch{}
}

var (
	r = rate.Every(time.Duration(time.Second * 90))
	L = rate.NewLimiter(r, 1)
)

var Swears ResponseMapping = ResponseMapping{
	Terms:     []string{"fuck", "shit", "dick", "pussy", "ass", "cunt", "whore", "bitch", "cock"},
	Responses: []string{"You need a spanking.", "You kiss your mother with that mouth?", "Excuse me, this is a *Christian server*. Thank you.", "Ooooooh, I'm telling mom on you!"},
}

// ContainsAny searchs a given string for the presence of any of the provided substrings.
func containsAny(s string, terms []string) bool {
	for _, t := range terms {
		if strings.Contains(s, t) {
			return true
		}
	}
	return false
}

func Chatter(c string, conf Conf, s *dgo.Session, m *dgo.MessageCreate) {
	r, err := Swears.Reply(c)

	if err != nil {
		return
	}

	if L.Allow() {
		s.ChannelMessageSend(m.ChannelID, r)
	}
}
