package main

import (
	"fmt"
	"log"
)

func logInfo(name string, args ...interface{}) {
	nArgs := []interface{}{fmt.Sprintf("[%s]", name)}
	nArgs = append(nArgs, args...)
	if !toStdout {
		log.Println(nArgs...)
	}
}
