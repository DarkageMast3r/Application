package Utilities

import (
	"os"
)

func StartTest() {
	err := findDir()
	if err != nil {
		panic("Correct Directory could not be found")
	}
	source, err := os.ReadFile("data/backup.db")
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("data/db.sqlite", source, 0777)
	if err != nil {
		panic("File could not be written")
	}
}