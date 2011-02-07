package main

import (
	"os"
	"bufio"
	"github.com/jbarham/cdb.go"
)

func main() {
	bin, bout := bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout)
	err := cdb.Dump(bout, bin)
	bout.Flush()
	if err != nil {
		os.Exit(111)
	}
}
