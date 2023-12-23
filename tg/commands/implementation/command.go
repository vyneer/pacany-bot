package implementation

import (
	"context"
	"fmt"

	"github.com/vyneer/tg-tagbot/db"
)

var Map = map[string]Command{}

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
	GetDescription() string
}

func CreateCommand(cmd func() Command) {
	c := cmd()

	if len(c.GetParentName()) == 0 {
		Map[c.GetName()] = c
		return
	}

	Map[fmt.Sprintf("%s %s", c.GetParentName(), c.GetName())] = c
}

func GetCommand(command, subcommand string) Command {
	key := fmt.Sprintf("%s %s", command, subcommand)

	if len(subcommand) == 0 {
		key = command
	}

	return Map[key]
}
