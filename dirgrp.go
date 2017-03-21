package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
)

const (
	srcDir = "/Volumes/My Passport for Mac"
	//srcDir   = "/Users/troy/Downloads"
	maxFiles = 5000
)

func main() {
	files, err := ioutil.ReadDir(srcDir)
	if err != nil {
		log.Fatal(err)
	}

	state := &destInfo{}

	for _, file := range files {

		if !state.isExtFile(file) {
			continue
		}
		// loop til we get a path we can use
		for state.PathIndex == 0 || state.FileIndex >= maxFiles {
			state.PathIndex++
			//log.Println(state)
			err := state.mkdir()
			if err != nil {
				log.Fatal("state.mkdir", err)
			}
		}
		state.FileIndex++
		err := os.Rename(path.Join(srcDir, file.Name()), state.Path(file.Name()))
		if err != nil {
			log.Fatal("os.rename", err)
		}

		fmt.Println(state, file.Name())
	}
}

type destInfo struct {
	FileIndex int
	PathIndex int
}

func (d *destInfo) Path(joins ...string) string {
	joins = append([]string{srcDir, "img" + strconv.Itoa(d.PathIndex)}, joins...)
	return path.Join(joins...)
}

func (d *destInfo) mkdir() error {
	newPath := d.Path()
	fmt.Println("mk", newPath)
	err := os.MkdirAll(newPath, 0777)
	if err != nil {
		return err
	}

	files, err := ioutil.ReadDir(newPath)
	if err != nil {
		log.Fatal(err)
	}
	// set FileIndex to number of matching files
	index := 0
	for _, file := range files {
		if d.isExtFile(file) {
			index++
		}
	}
	d.FileIndex = index

	return nil
}

func (d *destInfo) isExtFile(file os.FileInfo) bool {
	if file.IsDir() {
		return false
	}
	ext := path.Ext(file.Name())
	if ext == ".jpg" || ext == ".jpeg" || ext == ".JPG" || ext == ".JPEG" {
		return true
	}
	return false
}

func (d *destInfo) String() string {
	return fmt.Sprintf("%d:%d", d.PathIndex, d.FileIndex)
}
