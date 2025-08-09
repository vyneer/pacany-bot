package implementation

import (
	"context"
	"regexp"

	tgbotapiModels "github.com/go-telegram/bot/models"
	"github.com/vyneer/pacany-bot/db"
)

var (
	parent        = map[string]ParentCommand{}
	parentEnabled = map[string]ParentCommand{}
)

type ParentCommand interface {
	Name() string
	Description() string
	IsDisableable() bool
	Initialize() []Command
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
}

type InteractableCommand interface {
	Command
	GetName() string
	GetParentName() string
	GetHelp() (string, bool)
	GetDescription() (string, bool)
	IsAdminOnly() bool
}

type AutomaticCommand interface {
	Command
	GetIdentifier() string
	GetMatcher() *regexp.Regexp
}

func CreateParentCommand(cmd ParentCommand) {
	parent[cmd.Name()] = cmd
}

func EnableParentCommand(name string) {
	if cmd, ok := GetParentCommand(name); ok {
		parentEnabled[cmd.Name()] = cmd
	}
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
