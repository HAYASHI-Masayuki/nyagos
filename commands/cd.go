package commands

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/zetamatta/go-windows-shortcut"

	"github.com/nyaosorg/nyagos/nodos"
)

var cdHistory = make([]string, 0, 100)
var cdUniq = map[string]int{}

func pushCdHistory() {
	dir, err := os.Getwd()
	if err != nil {
		return
	}
	dirUpper := strings.ToUpper(dir)
	if i, ok := cdUniq[dirUpper]; ok {
		for ; i < len(cdHistory)-1; i++ {
			cdHistory[i] = cdHistory[i+1]
			cdUniq[strings.ToUpper(cdHistory[i])] = i
		}
		cdHistory[i] = dir
		cdUniq[dirUpper] = i
	} else {
		cdUniq[dirUpper] = len(cdHistory)
		cdHistory = append(cdHistory, dir)
	}
}

const (
	errnoChdirFail = 1
	errnoNoHistory = 2
)

func seekCdPath(dir string) string {
	if dir[0] == '.' || strings.ContainsAny(dir, "/\\:") {
		return ""
	}
	cdpath := os.Getenv("CDPATH")
	if cdpath == "" {
		return ""
	}
	for _, cdpath1 := range filepath.SplitList(cdpath) {
		fullpath := filepath.Join(cdpath1, dir)
		stat1, err := os.Stat(fullpath)
		if err == nil && stat1.IsDir() {
			return fullpath
		}
	}
	return ""
}

func cmdCdSub(dir string) (int, error) {
	const fileHead = "file:///"

	if strings.HasPrefix(dir, fileHead) {
		dir = dir[len(fileHead):]
	}
	if strings.HasSuffix(strings.ToLower(dir), ".lnk") {
		newdir, _, err := shortcut.Read(dir)
		if err == nil && newdir != "" {
			dir = newdir
		}
	}
	if dirTmp, err := CorrectCase(dir); err == nil {
		// println(dir, "->", dirTmp)
		dir = dirTmp
	}
	err := nodos.Chdir(dir)
	if err == nil {
		return 0, nil
	}
	if _dir := seekCdPath(dir); _dir != "" {
		if err = nodos.Chdir(_dir); err == nil {
			return 0, nil
		}
	}
	return errnoChdirFail, err
}

func cmdCd(ctx context.Context, cmd Param) (int, error) {
	args := cmd.Args()
	if len(args) >= 2 {
		if args[1] == "-" {
			if len(cdHistory) < 1 {
				return errnoNoHistory, errors.New("cd - : there is no previous directory")

			}
			directory := cdHistory[len(cdHistory)-1]
			pushCdHistory()
			return cmdCdSub(directory)
		} else if args[1] == "--history" {
			dir, err := os.Getwd()
			if err == nil {
				fmt.Fprintln(cmd.Out(), dir)
			} else {
				fmt.Fprintln(cmd.Err(), err.Error())
			}
			for i := len(cdHistory) - 1; i >= 0; i-- {
				fmt.Fprintln(cmd.Out(), cdHistory[i])
			}
			return 0, nil
		} else if args[1] == "-h" || args[1] == "?" {
			i := len(cdHistory) - 10
			if i < 0 {
				i = 0
			}
			for ; i < len(cdHistory); i++ {
				fmt.Fprintf(cmd.Out(), "cd %d => cd \"%s\"\n", i-len(cdHistory), cdHistory[i])
			}
			return 0, nil
		} else if i, err := strconv.ParseInt(args[1], 10, 0); err == nil && i < 0 {
			i += int64(len(cdHistory))
			if i < 0 {
				return errnoNoHistory, fmt.Errorf("cd %s: too old history", args[1])
			}
			directory := cdHistory[i]
			pushCdHistory()
			return cmdCdSub(directory)
		}
		if strings.EqualFold(args[1], "/D") {
			// ignore /D
			args = args[1:]
		}
		pushCdHistory()
		return cmdCdSub(strings.Join(args[1:], " "))
	}
	home := nodos.GetHome()
	if home != "" {
		pushCdHistory()
		return cmdCdSub(home)
	}
	return cmdPwd(ctx, cmd)
}
