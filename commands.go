package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	commMap = make(map[string]command)
)

//Small wrapper function to reduce clutter
func l(s string) string {
	return strings.ToLower(s)
}

func parseCommand(s *discordgo.Session, m *discordgo.MessageCreate, message string) {
	fmt.Println(message)
	msglist := strings.Fields(message)
	if len(msglist) == 0 {
		return
	}
	command := l(func() string {
		if strings.HasPrefix(message, " ") {
			return " " + msglist[0]
		}
		return msglist[0]
	}())

	submatch := emojiRegex.FindStringSubmatch(msglist[0])
	if len(submatch) > 0 {
		commMap["bigmoji"].Exec(s, m, append([]string{""}, msglist...))
		return
	}

	if fromMap, ok := commMap[command]; ok && command == l(fromMap.Name) {
		userPerms, err := s.State.UserChannelPermissions(m.Author.ID, m.ChannelID)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Error verifying permissions :(")
			errorLog.Println("error getting user permissions", err)
			return
		}

		hasPerms := userPerms&fromMap.PermsRequired > 0
		if (fromMap.RequiresPerms && hasPerms) || !fromMap.RequiresPerms {
			fromMap.Exec(s, m, msglist)
			return
		}
		s.ChannelMessageSend(m.ChannelID, "You don't have the correct permissions to run this!")
	}
}

func (c command) add() command {
	commMap[l(c.Name)] = c
	return c
}

func newCommand(name string, permissions int, noah, needsPerms bool, f func(*discordgo.Session, *discordgo.MessageCreate, []string)) command {
	return command{
		Name:          name,
		PermsRequired: permissions,
		RequiresPerms: needsPerms,
		NoahOnly:      noah,
		Exec:          f,
	}
}

func (c command) help() string {
	return c.Help
}

func (c command) sendHelp(s *discordgo.Session, id string) {
	s.ChannelMessageSend(id, c.help())
}

func (c command) setHelp(help string) command {
	c.Help = help
	return c
}
