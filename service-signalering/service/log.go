package service

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
)

func Log(val ...any) {
	file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()
	fmt.Fprintln(file, val...)
}

func LogInfo(val ...any) {
	val = append([]any{"[INFO] "}, val)
	Log(val)
}

func LogWarning(val ...any) {
	val = append([]any{"[WARNING] "}, val)
	val = append(val, "\nStacktrace: ")
	val = append(val, string(debug.Stack()))
	Log(val)
}

func LogError(val ...any) {
	val = append([]any{"[ERROR] "}, val)
	val = append(val, "\nStacktrace: ")
	val = append(val, string(debug.Stack()))
	Log(val)
}
