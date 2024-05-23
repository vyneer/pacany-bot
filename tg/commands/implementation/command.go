package implementation

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vyneer/pacany-bot/db"
)

var (
	Interactable      = map[string]Command{}
	InteractableOrder = []Command{}
	Automatic         = map[string]Command{}
)

type CommandArgs struct {
	DB      *db.DB
	ChatID  int64
	User    *tgbotapi.User
	IsAdmin bool
	Args    []string
}

type CommandResponse struct {
	Text       string
	Reply      bool
	Capitalize bool
}

type Command interface {
	Run(context.Context, CommandArgs) CommandResponse
	GetName() string
	GetParentName() string
	GetHelp() (string, bool)
	GetDescription() (string, bool)
	IsAdminOnly() bool
}

func CreateInteractableCommand(cmd func() Command) {
	c := cmd()

	Interactable[fmt.Sprintf("%s%s", c.GetParentName(), c.GetName())] = c
	InteractableOrder = append(InteractableOrder, c)
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
