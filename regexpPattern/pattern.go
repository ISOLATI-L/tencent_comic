package regexpPattern

import (
	"log"
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
