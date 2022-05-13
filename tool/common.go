package tool

import "log"

func CkErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
