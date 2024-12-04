package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/charlievieth/fastwalk"
)

func statFile(path string, fileInfo os.FileInfo) {
	fileInfo.Size()
}

type DirStatInfo struct {
	Path  string `json:"Path"`
	IsDir bool   `json:"IsDir"`
	Size  int64  `json:"Size"`
}

func exportJSON(w io.Writer, results []DirStatInfo) {
	json.NewEncoder(w).Encode(results)
}

const usageMsg = `Usage: %[1]s [PATH...]:

%[1]s dir disk usage utility

`

func main() {

	flag.Usage = func() {
		fmt.Fprintf(os.Stdout, usageMsg, filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	flag.Parse()

	args := flag.Args()

	if len(args) == 0 {
		args = append(args, ".")
	}

	conf := fastwalk.Config{
		Follow: false,
	}

	dirs := make([]DirStatInfo, 0)

	//fmt.Println(fastwalk.DefaultNumWorkers())

	walkFn := func(path string, info fs.DirEntry, err error) error {

		if err != nil {
			return nil
		}

		inf, _ := info.Info()

		dirs = append(dirs, DirStatInfo{Path: path, IsDir: info.IsDir(), Size: inf.Size()})

		return nil
	}

	for _, root := range args {

		err := fastwalk.Walk(&conf, root, walkFn)

		if err != nil {
			log.Println(err)
		}

	}

	exportJSON(os.Stdout, dirs)

}
