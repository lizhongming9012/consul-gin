package util

import "github.com/labstack/gommon/log"

func ShowError(where string, err error) {
	if err != nil {
		log.Printf("%s:%v", where, err)
	}
}
