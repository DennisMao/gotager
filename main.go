//   Copyright 2019 gotager authors
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"gotager/pkg/tagger"
	"io/ioutil"
	"os"
	"os/exec"
)

var (
	nameRegexp string // A regexp filter for the name of structure to tag on.
	prefix     string // Tag prefix.
	tagStyle   string // Tag style,support camel-like and snake-like.
	isBackup   bool   // If true,it will generate a backup on the same path with original go file.
	output     string // Specify a path, tagged file can be stogred in.If set,original file will not be modify and bakup file will not be generated.
)

func init() {
	flag.StringVar(&output, "o", "", "Specify a path, tagged file can be stogred in.If set,original file will not be modify and bakup file will not be generated.")
	flag.StringVar(&nameRegexp, "r", "*", "A regexp filter for the name of structure to tag on.")
	flag.StringVar(&prefix, "p", "json", "Prefix of the tag.'json' by default")
	flag.StringVar(&tagStyle, "s", "snake", "Style of the tag")
	flag.BoolVar(&isBackup, "b", true, "If true,it will generate a backup on the same folder with original go file.")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of gotager:\n")
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), "Example:\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "        gotager -p=json ./example/test.go \n", os.Args[0])
	}

}

func main() {
	flag.Parse()

	var filePath = os.Args[len(os.Args)-1]
	var fOrigin, fTarget *os.File

	fOrigin, err := os.OpenFile(filePath, os.O_SYNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		fmt.Printf("Can not open file: %s ,please check the param! \n", filePath)
		os.Exit(-1)
		return
	}
	defer fOrigin.Close()

	if output != "" {
		fTarget, err := os.OpenFile(output, os.O_SYNC|os.O_RDWR|os.O_CREATE, os.ModePerm)
		if err != nil {
			fmt.Printf("Can not open file: %s ,please check the param! \n", output)
			os.Exit(-1)
			return
		}
		defer fTarget.Close()
	} else {
		fTarget = fOrigin
	}

	src, err := ioutil.ReadAll(fOrigin)
	if err != nil {
		fmt.Printf("Can not read datas from the file: %s ,please check the param! \n", filePath)
		os.Exit(-1)
		return
	}

	buf := bytes.NewBuffer(make([]byte, 0))
	buf.Reset()

	// execute tag process
	opt := &tagger.TagOpt{false, tagStyle}
	t := tagger.New(opt)
	err = t.Tag(src, buf, nameRegexp, prefix)
	if err != nil {
		fmt.Printf("Tag failed,Detail:%s\n", err.Error())
		os.Exit(-1)
		return
	}

	if isBackup && output == "" {
		exec.Command("cp", filePath, filePath+".bak").Run()
	}

	fTarget.Seek(0, 0)
	fTarget.Write(buf.Bytes())
	fmt.Printf("Tag successfully!")
	os.Exit(0)
	return

}
