package main

import (
	"errors"
	"log"
	"os"
	"regexp"
)

var TitlePattern *regexp.Regexp

func init() {
	var err error
	TitlePattern, err = regexp.Compile(`<\s*title\s*>\s*(.*?)\s*<\s*/title\s*>`)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

const DEBUG bool = true

func savehtml(body []byte) (err error) {
	if DEBUG {
		match := TitlePattern.FindStringSubmatch(string(body))
		if len(match) > 0 {
			htmlFile, err := os.OpenFile(
				match[1]+".html",
				os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
				0666,
			)
			if err != nil {
				return err
			}
			htmlFile.Write(body)
			htmlFile.Close()
		} else {
			return errors.New("can not find title")
		}
	}
	return nil
}
