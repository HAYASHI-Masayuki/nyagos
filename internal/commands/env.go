package commands

import (
	"context"
	"fmt"
	"os"
	"strings"
)

func array2hash(args []string) ([]string, map[string]string) {
	hash := map[string]string{}
	for i, arg1 := range args {
		equalPos := strings.IndexRune(arg1, '=')
		if equalPos < 0 {
			return args[i:], hash
		}
		key := arg1[:equalPos]
		val := arg1[equalPos+1:]
		hash[key] = val
	}
	return []string{}, hash
}

func cmdEnv(ctx context.Context, cmd Param) (int, error) {
	args, hash := array2hash(cmd.Args()[1:])
	if len(args) <= 0 {
		for _, val := range os.Environ() {
			fmt.Fprintln(cmd.Out(), val)
		}
		return 0, nil
	}
	return cmd.Spawnlpe(ctx, args, args, hash)
}
