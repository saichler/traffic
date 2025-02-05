package main

import (
	"github.com/saichler/shared/go/share/logger"
	"github.com/saichler/traffic/go/generator/cmd"
	"os"
)

func main() {
	log := logger.NewLoggerDirectImpl(&logger.FmtLogMethod{})
	cmds := cmd.NewCommands(log)
	cmds.Run(os.Args)
}
