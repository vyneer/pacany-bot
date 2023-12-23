package tag

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/vyneer/tg-tagbot/tg/commands/implementation"
)

const (
	name              string = "tag"
	parentName        string = "help"
	help              string = "Print help for the tag commands"
	helpOrder         int    = -1
	shape             string = "/helptag"
	showInCommandList bool   = true
)

type Command struct{}

func New() implementation.Command {
	return &Command{}
}

func (c *Command) GetName() string {
	return name
}

func (c *Command) GetParentName() string {
	return parentName
}

func (c *Command) GetHelp() (string, int) {
	return fmt.Sprintf("%s - %s", shape, help), helpOrder
}

func (c *Command) GetDescription() string {
	if !showInCommandList {
		return ""
	}
	return fmt.Sprintf("%s - %s", help, shape)
}

func (c *Command) Run(_ context.Context, _ implementation.CommandArgs) implementation.CommandResponse {
	helpMap := map[int]string{}

	for _, v := range implementation.Interactable {
		if v.GetParentName() != "tag" {
			continue
		}
		helpString, order := v.GetHelp()
		if order != -1 {
			helpMap[order] = helpString
		}
	}

	keys := make([]int, 0, len(helpMap))
	for k := range helpMap {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	helpSlice := make([]string, 0, len(helpMap))
	for _, i := range keys {
		helpSlice = append(helpSlice, helpMap[i])
	}

	return implementation.CommandResponse{
		Text:  strings.Join(helpSlice, "\n\n"),
		Reply: true,
	}
}
