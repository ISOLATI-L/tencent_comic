package regexpPattern

import (
	"log"
	"regexp"
)

var TitlePattern *regexp.Regexp
var ChapterPattern *regexp.Regexp

func init() {
	var err error
	TitlePattern, err = regexp.Compile(`<\s*title\s*>\s*(.*?)\s*<\s*/title\s*>`)
	if err != nil {
		log.Fatalln(err.Error())
	}
	ChapterPattern, err = regexp.Compile(`//*[@id="chapter"]/div[2]/ol[1]/li`)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
