package cmd

import (
	"github.com/saichler/shared/go/share/interfaces"
	"os"
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

func (this *Commands) Run() error {
	if len(os.Args) < 2 {
		return this.log.Error("no command specified")
	}
	cmd := this.cmds[os.Args[1]]
	if cmd == nil {
		return this.log.Error("no " + os.Args[1] + " command was found")
	}
	if len(os.Args) == 2 {
		this.log.Info(cmd.Help())
		return nil
	}
	err := fillArgs(cmd)
	if err != nil {
		return this.log.Error(err.Error())
	}
	cmd.Run(this.log)
	return nil
}
