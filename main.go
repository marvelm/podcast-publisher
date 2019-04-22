package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/jbub/podcasts"
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

func removeFirstDirectoryIfNecessary(str string) string {
	parts := strings.Split(str, "/")
	if len(parts) <= 1 {
		return str
	}

	return strings.Join(parts[1:], "/")
}

func main() {
	folder := flag.String("folder", "", "Directory to check for files")
	urlFlag := flag.String("url", "", "Base url like https://my-podcasts.com/somepodcast")
	title := flag.String("title", "My podcast", "Title of the podcast")
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

	podcast := &podcasts.Podcast{
		Title: *title,
	}

	currentWorkingDirectory := filepath.Dir(os.Args[0])

	_ = filepath.Walk(*folder, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			panic(err)
		}

		if info.IsDir() {
			return nil
		}

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

		podcast.AddItem(&podcasts.Item{
			Title: strings.Replace(removeFirstDirectoryIfNecessary(fileLocation), "/", " ", -1),
			PubDate: &podcasts.PubDate{info.ModTime()},
			GUID:    hash(filePath),
			Enclosure: &podcasts.Enclosure{
				URL:  downloadUrl.String(),
				Type: mime,
			},
		})
		return nil
	})

	feed, err := podcast.Feed()
	if err != nil {
		panic(err)
	}

	err = feed.Write(os.Stdout)
	if err != nil {
		panic(err)
	}
}
