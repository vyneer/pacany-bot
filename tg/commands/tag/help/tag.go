package help

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/vyneer/pacani-bot/tg/commands/implementation"
)

const (
	name              string = "help"
	parentName        string = "tag"
	help              string = "Print help for the tag commands"
	helpOrder         int    = -1
	shape             string = "/taghelp"
	descriptionOrder  int    = 1
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

func (c *Command) GetDescription() (string, int) {
	if !showInCommandList {
		return "", descriptionOrder
	}
	return fmt.Sprintf("%s - %s", help, shape), descriptionOrder
}

func (c *Command) Run(_ context.Context, _ implementation.CommandArgs) implementation.CommandResponse {
	helpMap := map[int]string{}

	for _, v := range implementation.Interactable {
		if v.GetParentName() != parentName {
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
