package main

import (
	"os"
	"fmt"
	"log"
	"flag"
	"bufio"
	"path"
	"io/ioutil"
	"github.com/jbarham/cdb.go"
)

func exitOnErr(err os.Error) {
	if err != nil {
		log.Fatal(err)
	}
}

func usage() {
	fmt.Fprint(os.Stderr, "usage: cdbmake f [ftmp]\n")
	os.Exit(2)
}

func main() {
	var tmp *os.File
	var err os.Error

	flag.Parse()
	args := flag.Args()
	if len(args) == 1 {
		dir, _ := path.Split(args[0])
		tmp, err = ioutil.TempFile(dir, "")
		exitOnErr(err)
	} else if len(args) == 2 {
		tmp, err = os.OpenFile(args[1], os.O_RDWR|os.O_CREATE, 0644)
		exitOnErr(err)
	} else {
		usage()
	}

	fname := args[0]
	tmpname := tmp.Name()

	exitOnErr(cdb.Make(tmp, bufio.NewReader(os.Stdin)))
	exitOnErr(tmp.Sync())
	exitOnErr(tmp.Close())
	exitOnErr(os.Rename(tmpname, fname))
}
