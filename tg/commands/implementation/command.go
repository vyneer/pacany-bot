package implementation

import (
	"context"
	"fmt"

	tgbotapiModels "github.com/go-telegram/bot/models"
	"github.com/vyneer/pacany-bot/db"
)

var (
	parent        = map[string]ParentCommand{}
	parentEnabled = map[string]ParentCommand{}

	interactable      = map[string]Command{}
	interactableOrder = []Command{}

	automatic = map[string]Command{}
)

type ParentCommand interface {
	Name() string
	Description() string
	Initialize()
}

type CommandArgs struct {
	DB      *db.DB
	ChatID  int64
	User    *tgbotapiModels.User
	IsAdmin bool
	Args    []string
}

type CommandResponse struct {
	Text       string
	Reply      bool
	Capitalize bool
}

type Command interface {
	Run(context.Context, CommandArgs) []CommandResponse
	GetName() string
	GetParentName() string
	GetHelp() (string, bool)
	GetDescription() (string, bool)
	IsAdminOnly() bool
}

func CreateParentCommand(cmd ParentCommand) {
	parent[cmd.Name()] = cmd
}

func EnableParentCommand(name string) {
	if cmd, ok := GetParentCommand(name); ok {
		parentEnabled[cmd.Name()] = cmd
	}
}

func CreateInteractableCommand(cmd func() Command) {
	c := cmd()

	interactable[fmt.Sprintf("%s%s", c.GetParentName(), c.GetName())] = c
	interactableOrder = append(interactableOrder, c)
}

func CreateAutomaticCommand(cmd func() Command) {
	c := cmd()

	automatic[fmt.Sprintf("%s%s", c.GetParentName(), c.GetName())] = c
}

func GetParentCommand(name string) (ParentCommand, bool) {
	c, ok := parent[name]
	return c, ok
}

func GetAllParentCommands() map[string]ParentCommand {
	return parent
}

func GetEnabledParentCommands() map[string]ParentCommand {
	return parentEnabled
}

func GetInteractableCommand(command string) Command {
	return interactable[command]
}

func GetInteractableOrder() []Command {
	return interactableOrder
}

func GetAutomaticCommand(command string) Command {
	return automatic[command]
}
