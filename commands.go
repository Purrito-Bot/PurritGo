package main

import (
	"strconv"
	"strings"
)

func handleName(c parsedCommand) {
	if len(c.args) < 1 {
		c.session.ChannelMessageSend(c.message.ChannelID, "Usage: go name <variant>")
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
