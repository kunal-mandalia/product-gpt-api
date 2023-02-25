package util

import (
	"io"
	"io/ioutil"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func WriteToDisk(name string, content io.ReadCloser) {
	body, err1 := ioutil.ReadAll(content)
	check(err1)

	err2 := os.WriteFile("./"+name, body, 0644)
	check(err2)
}
