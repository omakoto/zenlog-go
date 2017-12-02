package main

import (
	"github.com/omakoto/zenlog-go/zenlog"
	"github.com/omakoto/zenlog-go/zenlog/builtins"
	"github.com/omakoto/zenlog-go/zenlog/util"
	"os"
)

func restart() {
	util.Say("Restarting zenlog...")
	util.MustExec(util.StringSlice(util.FindZenlogBin()))
}

func realMain() int {
	command, args := util.GetSubcommand()

	if command == "" {
		builtins.FailIfInZenlog()
		status, resurrect := zenlog.StartZenlog(args)
		if resurrect {
			restart()
		}
		return status
	}
	builtins.MaybeRunBuiltin(command, args)
	MaybeRunExternalCommand(command, args)

	util.Fatalf("Unknown subcommand: '%s'", command)
	return 0
}

func main() {
	os.Exit(util.RunWithRescue(realMain))
}
