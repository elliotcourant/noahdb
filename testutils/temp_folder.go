package testutils

import (
	"fmt"
	"io/ioutil"
	"os"
)

func TempFile(fileName string) (string, func()) {
	path, remove := TempFolder()
	return fmt.Sprintf("%s/%s", path, fileName), remove
}

func TempFolder() (string, func()) {
	folder, err := ioutil.TempDir("", "noahdb")
	if err != nil {
		panic(err)
	}
	return folder, func() {
		if err = os.RemoveAll(folder); err != nil {
			panic(err)
		}
	}
}
