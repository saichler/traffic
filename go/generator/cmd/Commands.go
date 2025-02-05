package cmd

import (
	"github.com/saichler/shared/go/share/interfaces"
	"reflect"
)

type Commands struct {
	cmds map[string]Command
	log  interfaces.ILogger
}

func NewCommands(log interfaces.ILogger) *Commands {
	cmds := &Commands{}
	cmds.log = log
	cmds.cmds = make(map[string]Command)
	cmds.addCommand(&Start{})
	cmds.addCommand(&Do{})
	return cmds
}

func (this *Commands) addCommand(cmd Command) {
	this.cmds[cmd.Name()] = cmd
}

func (this *Commands) Run(args []string) string {
	if len(args) < 2 {
		return this.log.Error("no command specified").Error()
	}
	cmd := this.cmds[args[1]]
	if cmd == nil {
		return this.log.Error("no " + args[1] + " command was found").Error()
	}
	cmd = reflect.New(reflect.ValueOf(cmd).Elem().Type()).Interface().(Command)
	if len(args) == 2 {
		this.log.Info(cmd.Help())
		return cmd.Help()
	}
	err := fillArgs(cmd, args)
	if err != nil {
		return this.log.Error(err.Error()).Error()
	}
	return cmd.Run(this.log)
}
