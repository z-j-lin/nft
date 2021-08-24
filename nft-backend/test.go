package test

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"
)

func readCurrentDir() {
	dir := "./tmp/."
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatalf("failed opening directory: %s", err)
	}
	var modTime time.Time
	var mostRecentFi string
	for _, fi := range files {
		fmt.Println(fi.Name(), fi.ModTime())
		if newModtime := fi.ModTime(); newModtime.After(modTime) {
			modTime = fi.ModTime()
			mostRecentFi = fi.Name()
		}
	}
	fmt.Println(mostRecentFi)
}

func main() {
	readCurrentDir()
}
