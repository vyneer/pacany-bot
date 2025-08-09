package implementation

import (
	"context"
	"regexp"

	tgbotapiModels "github.com/go-telegram/bot/models"
	"github.com/vyneer/pacany-bot/db"
)

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
