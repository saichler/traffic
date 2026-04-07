package main

import (
	"github.com/saichler/l8traffic/go/generator/cmd"
	"github.com/saichler/l8utils/go/utils/logger"
	"os"
)

func main() {
	log := logger.NewLoggerDirectImpl(&logger.FmtLogMethod{})
	cmds := cmd.NewCommands(log)
	cmds.Run(os.Args)
}
