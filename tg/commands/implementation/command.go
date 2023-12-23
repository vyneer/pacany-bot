package implementation

import (
	"context"
	"fmt"

	"github.com/vyneer/tg-tagbot/db"
)

var (
	Interactable = map[string]Command{}
	Automatic    = map[string]Command{}
)

type CommandArgs struct {
	DB     *db.DB
	ChatID int64
	Args   []string
}

type CommandResponse struct {
	Text  string
	Reply bool
}

type Command interface {
	Run(context.Context, CommandArgs) CommandResponse
	GetName() string
	GetParentName() string
	GetHelp() (string, int)
	GetDescription() (string, int)
}

func CreateInteractableCommand(cmd func() Command) {
	c := cmd()

	Interactable[fmt.Sprintf("%s%s", c.GetParentName(), c.GetName())] = c
}

func CreateAutomaticCommand(cmd func() Command) {
	c := cmd()

	Automatic[fmt.Sprintf("%s%s", c.GetParentName(), c.GetName())] = c
}

func GetInteractableCommand(command string) Command {
	return Interactable[command]
}

func GetAutomaticCommand(command string) Command {
	return Automatic[command]
}
