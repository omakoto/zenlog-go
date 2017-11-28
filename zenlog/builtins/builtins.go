package builtins

import (
	"fmt"
	"github.com/omakoto/zenlog-go/zenlog/builtins/history"
	"github.com/omakoto/zenlog-go/zenlog/envs"
	"github.com/omakoto/zenlog-go/zenlog/util"
	"io"
	"os"
	"strconv"
	"strings"
)

func InZenlog() bool {
	return util.Tty() == os.Getenv(envs.ZENLOG_TTY)
}

func FailIfInZenlog() {
	if InZenlog() {
		util.Fatalf("Already in zenlog.")
	}
}

func FailUnlessInZenlog() {
	if !InZenlog() {
		util.Fatalf("Not in zenlog.")
	}
}

func copyStdinToFile(file string) {
	out, err := os.OpenFile(file, os.O_WRONLY, 0)
	util.Check(err, "Unable to open "+file)
	io.Copy(out, os.Stdin)
}

func WriteToLogger() {
	FailUnlessInZenlog()
	copyStdinToFile(os.Getenv(envs.ZENLOG_LOGGER_IN))
}

func WriteToOuter() {
	FailUnlessInZenlog()
	copyStdinToFile(os.Getenv(envs.ZENLOG_OUTER_TTY))
}

func OuterTty() {
	FailUnlessInZenlog()
	fmt.Println(os.Getenv(envs.ZENLOG_OUTER_TTY))
}

func MaybeRunBuiltin(command string, args []string) {
	switch strings.Replace(command, "_", "-", -1) {
	case "in-zenlog":
		util.Exit(InZenlog())

	case "fail-if-in-zenlog":
		FailIfInZenlog()
		os.Exit(0)

	case "fail-unless-in-zenlog":
		FailUnlessInZenlog()
		os.Exit(0)

	case "write-to-logger":
		WriteToLogger()
		os.Exit(0)

	case "write-to-outer":
		WriteToOuter()
		os.Exit(0)

	case "outer-tty":
		OuterTty()
		os.Exit(0)

		// TODO Refactor these commands for testability.
	case "start-command":
		FailUnlessInZenlog()
		if len(args) < 1 {
			util.Fatalf("start-command expects 1 argument.")
		}
		StartCommand("", args[:], util.NewClock())

	case "start-command-with-env":
		FailUnlessInZenlog()
		if len(args) < 2 {
			util.Fatalf("start-command expects 2 arguments.")
		}
		StartCommand(args[0], args[1:], util.NewClock())

	case "stop-log", "end-command":
		FailUnlessInZenlog()

		wantLineNumber := false
		i := 0
		if i < len(args) && args[i] == "-n" {
			i++
			wantLineNumber = true
		}
		exitStatus := -1
		var err error
		if i < len(args) {
			exitStatus, err = strconv.Atoi(args[i])
			util.Check(err, "Exit status must be integer; '%s' given.", args[i])
			i++
		}
		EndCommand(exitStatus, wantLineNumber, util.NewClock())
	case "history":
		FailUnlessInZenlog()
		history.HistoryCommand(args)
	case "current-log":
		FailUnlessInZenlog()
		history.CurrentLogCommand(args)
	case "last-log":
		FailUnlessInZenlog()
		history.LastLogCommand(args)
	}
	return
}
