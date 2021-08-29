package main

import (
	"errors"
	"os"
	"tencent_comic/regexpPattern"
)

func savehtml(body []byte) (err error) {
	match := regexpPattern.TitlePattern.FindStringSubmatch(string(body))
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
	return nil
}
