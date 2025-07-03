package Utilities

import (
	"os"
)

func FindDir() error {
	for {
		entries, err := os.ReadDir(".")
		if err != nil {
			return err
		}
		for _, entry := range entries {
			if entry.Name() == "Templates" {
				return nil
			}
		}
		err = os.Chdir("..")
		if err != nil {
			return err
		}
	}
}