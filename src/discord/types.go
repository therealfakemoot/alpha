package discord

import (
	dgo "github.com/bwmarrin/discordgo"
)

// Embedder describes a value that can create a discord rich bembed of itself.
type Embedder interface {
	Embed(bool) dgo.MessageEmbed
}
