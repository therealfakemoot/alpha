package discord

import (
	dgo "github.com/bwmarrin/discordgo"
	"time"
)

// MessagePolicy is used to control how a message is handled after being sent.
//
// Delete: a time.Duration after which the message will delete itself.
//   A Duration of 0 means tat the message will not be deleted.
type MessagePolicy struct {
	Delete time.Duration
}

// Convenience MessagePolicy values to cover common cases.
var (
	OneSecondPolicy    = MessagePolicy{Delete: 1 * time.Second}
	FiveSecondPolicy   = MessagePolicy{Delete: 5 * time.Second}
	TenSecondPolicy    = MessagePolicy{Delete: 10 * time.Second}
	ThirtySecondPolicy = MessagePolicy{Delete: 30 * time.Second}
	OneMinutePolicy    = MessagePolicy{Delete: 1 * time.Minute}
	FiveMinutesPolicy  = MessagePolicy{Delete: 5 * time.Minute}
)

// Session interface will allow me to mock out my discordgo session to facilitate testing.
type Session interface {
	ChannelMessageSend(string, string) (*dgo.Message, error)
	ChannelMessageSendEmbed(string, *dgo.MessageEmbed) (*dgo.Message, error)
}

// NewEnvoy creates a value that
func NewEnvoy(s *Session, d time.Duration, p MessagePolicy) *Envoy {
	return &Envoy{Session: s, Duration: d, Policy: p}
}

// Envoy is a 'factory' type that is constructed with a MessagePolicy and create messages accordingly.
//
// The name "envoy" was chosen to represent that this value is an independent actor following predetermined orders.
// It can be used for batched sends or repetitive sends ( periodic channel updates, or messages that will always self-desturct, etc ).
type Envoy struct {
	Session  *Session
	Duration time.Duration
	Policy   MessagePolicy
}

// SendMessage will send a message to a given channel, if possible.
//
// SendMessage will respect the MessagePolicy given to the Envoy at creation.
func (e *Envoy) SendMessage(content string) error {
	return nil
}

// SendEmbed will send rich embeds according to its MessagePolicy.
func (e *Envoy) SendEmbed(embed *dgo.MessageEmbed) error {
	return nil
}
