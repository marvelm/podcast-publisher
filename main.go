package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/jbub/podcasts"
	"github.com/gabriel-vasile/mimetype"
)

func hash(str string) string {
	sum := md5.Sum([]byte(str))
	return hex.EncodeToString(sum[:])
}

func addPathsToUrl(baseUrl *url.URL, paths ...string) (*url.URL, error) {
	u, err := url.Parse(baseUrl.String())
	if err != nil {
		return nil, err
	}

	pathsConcat := append([]string{u.Path}, paths...)
	finalPath := path.Join(pathsConcat...)
	u.Path = finalPath

	return u, nil
}

func main() {
	folder := flag.String("folder", "", "Directory to check for files")
	urlFlag := flag.String("url", "", "Base url like https://my-podcasts.com/somepodcast")
	flag.Parse()

	if *folder == "" {
		panic("folder must be defined")
	}

	if *urlFlag == "" {
		panic("url must be defined")
	}

	parsedUrl, err := url.Parse(*urlFlag)
	if err != nil {
		panic(err)
	}

	var files []string
	err = filepath.Walk(*folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	podcast := &podcasts.Podcast{
		Title: "My podcast",
	}

	currentWorkingDirectory := filepath.Dir(os.Args[0])
	for _, filePath := range files {
		fileLocation, err := filepath.Rel(currentWorkingDirectory, filePath)
		if err != nil {
			panic(err)
		}

		downloadUrl, err := addPathsToUrl(parsedUrl, fileLocation)
		if err != nil {
			panic(err)
		}

		mime, _, err := mimetype.DetectFile(filePath)
		if err != nil {
			panic(err)
		}

		filename := filepath.Base(filePath)

		podcast.AddItem(&podcasts.Item{
			Title: filename,
			GUID:  hash(filePath),
			Enclosure: &podcasts.Enclosure{
				URL: downloadUrl.String(),
				Type: mime,
			},
		})
	}

	feed, err := podcast.Feed()
	if err != nil {
		panic(err)
	}

	err = feed.Write(os.Stdout)
	if err != nil {
		panic(err)
	}
}
