package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/purrito-bot/purrigo/voice"
)

func handleName(c parsedCommand) {
	if len(c.args) < 1 {
		c.session.ChannelMessageSend(c.message.ChannelID, "Usage: go name <variant> <count?>")
		return
	}
	if c.args[0] == "variants" {
		message := fmt.Sprintf("I can generate the following types of name:\n%s", strings.Join(generator.Variants(), "\n"))
		c.session.ChannelMessageSend(c.message.ChannelID, message)
		return
	}
	count := 1
	if len(c.args) >= 2 {
		count, _ = strconv.Atoi(c.args[1])
	}
	names := make([]string, count)
	for i := 0; i < count; i++ {
		name, err := generator.GenerateName(c.args[0])
		if err != nil {
			c.session.ChannelMessageSend(c.message.ChannelID, err.Error())
			return
		}
		names = append(names, name)
	}

	c.session.ChannelMessageSend(c.message.ChannelID, strings.Join(names, "\n"))
}

func handleShow(c parsedCommand) {
	urls := []string{}
	for _, v := range c.message.Mentions {
		urls = append(urls, v.AvatarURL(""))
	}
	c.session.ChannelMessageSend(c.message.ChannelID, strings.Join(urls, "\n"))
}

func handleSpeak(c parsedCommand) {
	// Find the channel that the message came from.
	ch, err := c.session.State.Channel(c.message.ChannelID)
	if err != nil {
		// Could not find channel.
		return
	}

	// Find the guild for that channel.
	g, err := c.session.State.Guild(ch.GuildID)
	if err != nil {
		// Could not find guild.
		return
	}

	// Look for the message sender in that guild's current voice states.
	for _, vs := range g.VoiceStates {
		if vs.UserID == c.message.Author.ID {
			err = voice.PlaySound(c.session, g.ID, vs.ChannelID, meowBuffer)
			if err != nil {
				fmt.Println("Error playing sound:", err)
			}

			return
		}
	}
}
