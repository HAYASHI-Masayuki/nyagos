package frame

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mattn/go-colorable"

	"github.com/zetamatta/go-readline-ny"
	"github.com/zetamatta/go-windows-consoleicon"

	"github.com/zetamatta/nyagos/history"
	"github.com/zetamatta/nyagos/shell"
)

type CmdStreamConsole struct {
	shell.CmdSeeker
	DoPrompt func() (int, error)
	History  *history.Container
	Editor   *readline.Editor
	HistPath string
}

func NewCmdStreamConsole(doPrompt func() (int, error)) *CmdStreamConsole {
	history1 := &history.Container{}
	this := &CmdStreamConsole{
		History: history1,
		Editor: &readline.Editor{
			History: history1,
			Prompt:  doPrompt,
			Writer:  colorable.NewColorableStdout()},
		HistPath: filepath.Join(AppDataDir(), "nyagos.history"),
		CmdSeeker: shell.CmdSeeker{
			PlainHistory: []string{},
			Pointer:      -1,
		},
	}
	history1.Load(this.HistPath)
	history1.Save(this.HistPath)
	return this
}

func (this *CmdStreamConsole) DisableHistory(value bool) bool {
	return this.History.IgnorePush(value)
}

func (this *CmdStreamConsole) ReadLine(ctx context.Context) (context.Context, string, error) {
	if this.Pointer >= 0 {
		if this.Pointer < len(this.PlainHistory) {
			this.Pointer++
			return ctx, this.PlainHistory[this.Pointer-1], nil
		}
		this.Pointer = -1
	}
	var line string
	var err error
	for {
		disabler := colorable.EnableColorsStdout(nil)
		clean, err2 := consoleicon.SetFromExe()
		for {
			line, err = this.Editor.ReadLine(ctx)
			if err != readline.CtrlC {
				break
			}
			fmt.Fprintln(os.Stderr, err.Error())
		}
		if err2 == nil {
			clean(false)
		}
		disabler()
		if err != nil {
			return ctx, line, err
		}
		var isReplaced bool
		line, isReplaced, err = this.History.Replace(line)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			continue
		}
		if isReplaced {
			fmt.Fprintln(os.Stdout, line)
		}
		if line != "" {
			break
		}
	}
	row := history.NewHistoryLine(line)
	this.History.PushLine(row)
	fd, err := os.OpenFile(this.HistPath, os.O_APPEND|os.O_CREATE, 0600)
	if err == nil {
		fmt.Fprintln(fd, row.String())
		fd.Close()
	} else {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	this.PlainHistory = append(this.PlainHistory, line)
	return ctx, line, err
}
