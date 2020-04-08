package main

import (
	"log"
	"os"
	"path/filepath"
)

func GetCurrentPath() string {
	exec, err := os.Executable()
	if err != nil {
		log.Fatal("Failed to get an executable")
	}

	return filepath.Dir(exec)
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
