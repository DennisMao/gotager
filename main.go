package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"gotager/pkg/tagger"
)

var (
	filePath    string // Specified a file need to add tags for
	nameRegexp  string // A regexp filter for the name of structure to tag on.
	prefix      string // Tag prefix.
	tagStyle    string // Tag style,support camel-like and snake-like.
	isBackup    bool   // If true,it will generate a backup on the same path with original go file.
	isOverwrite bool   // If ture,it will overwrite the tag with specified prefix when it exists.
)

func init() {
	flag.StringVar(&nameRegexp, "r", "*", "A regexp filter for the name of structure to tag on.")
	flag.StringVar(&prefix, "p", "json", "Prefix of the tag.'json' by default")
	flag.StringVar(&tagStyle, "s", "camel", "Style of the tag")
	flag.BoolVar(&isBackup, "b", true, "If true,it will generate a backup on the same folder with original go file.")
	flag.BoolVar(&isOverwrite, "o", true, "If ture,it will overwrite the tag with specified prefix when it exists.")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of gotager:\n")
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), "Example:\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "        gotager -p=json ./example/test.go \n", os.Args[0])
	}

}

func main() {
	flag.Parse()

	filePath = os.Args[len(os.Args)-1]

	f, err := os.OpenFile(filePath, os.O_SYNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		fmt.Printf("Can not open file: %s ,please check the param! \n", filePath)
		os.Exit(-1)
		return
	}
	defer f.Close()

	src, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Printf("Can not read datas from the file: %s ,please check the param! \n", filePath)
		os.Exit(-1)
		return
	}

	buf := bytes.NewBuffer(make([]byte, 0))
	buf.Reset()

	fmt.Println(filePath)

	opt := &tagger.TagOpt{isOverwrite, tagStyle}
	t := tagger.New(opt)
	err = t.Tag(src, buf, nameRegexp, prefix)
	if err != nil {
		fmt.Printf("Tag failed,Detail:%s\n", err.Error())
		os.Exit(-1)
		return
	}

	if isBackup {
		exec.Command("cp", filePath, filePath+".bak").Run()
	}

	f.Seek(0, 0)
	f.Write(buf.Bytes())
	fmt.Printf("Tag successfully!")
	os.Exit(0)
	return

}
