package main

import (
	"io"
	"os"
	"strconv"

	"github.com/jbarham/go-cdb"
)

const usage = "usage: cdbget key [skip]"

func main() {

	var (
		key  []byte
		skip int
		err  error
	)

	if len(os.Args) < 2 {
		exitWithMsg(1, usage)
	}
	if len(os.Args) > 1 {
		key = []byte(os.Args[1])
	}
	if len(os.Args) > 2 {
		if skip, err = strconv.Atoi(os.Args[2]); err != nil {
			exitWithMsg(2, "error:", "parsing 'skip'-error", err.Error())
		}
		if skip < 0 {
			exitWithMsg(2, "error:", "skip parameter is invalid (negativ)")
		}
	}

	c := cdb.New(os.Stdin)
	if _, err := cdb.Get(os.Stdout, c, key, skip); err != nil {
		exitWithMsg(3, err.Error())
	}
}

func exitWithMsg(c int, msg ...string) {
	if len(msg) > 0 {
		for _, m := range msg {
			io.WriteString(os.Stderr, m)
		}
		io.WriteString(os.Stderr, "\n")
	}
	os.Exit(c)
}
