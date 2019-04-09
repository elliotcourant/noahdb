package static

import (
	"github.com/elliotcourant/statik/fs"
)

//go:generate make -C ../../../ embedded

func GetEmbeddedFile(fileName string) ([]byte, error) {
	fileSystem, err := fs.New()
	if err != nil {
		panic(err)
	}
	file, err := fileSystem.Open(fileName)
	if err != nil {
		return nil, err
	}
	stats, _ := file.Stat()
	size := stats.Size()
	bytes := make([]byte, size)
	_, _ = file.Read(bytes)
	return bytes, nil
}
