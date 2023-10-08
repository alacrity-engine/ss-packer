package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

	codec "github.com/alacrity-engine/resource-codec"
	"github.com/golang-collections/collections/queue"
	bolt "go.etcd.io/bbolt"
)

var (
	projectPath      string
	resourceFilePath string
)

func parseFlags() {
	flag.StringVar(&projectPath, "project", ".",
		"Path to the project to pack spritesheets for.")
	flag.StringVar(&resourceFilePath, "out", "./stage.res",
		"Resource file to store animations and spritesheets.")

	flag.Parse()
}

func main() {
	parseFlags()

	resourceFile, err := bolt.Open(resourceFilePath, 0666, nil)
	handleError(err)
	defer resourceFile.Close()

	entries, err := os.ReadDir(projectPath)
	handleError(err)

	traverseQueue := queue.New()

	if len(entries) <= 0 {
		return
	}

	for _, entry := range entries {
		traverseQueue.Enqueue(FileTracker{
			EntryPath: ".",
			Entry:     entry,
		})
	}

	for traverseQueue.Len() > 0 {
		fsEntry := traverseQueue.Dequeue().(FileTracker)

		if fsEntry.Entry.IsDir() {
			entries, err = os.ReadDir(path.Join(fsEntry.EntryPath, fsEntry.Entry.Name()))
			handleError(err)

			for _, entry := range entries {
				traverseQueue.Enqueue(FileTracker{
					EntryPath: path.Join(fsEntry.EntryPath, fsEntry.Entry.Name()),
					Entry:     entry,
				})
			}

			continue
		}

		if !strings.HasSuffix(fsEntry.Entry.Name(), ".ss.yml") {
			continue
		}

		data, err := os.ReadFile(path.Join(fsEntry.EntryPath, fsEntry.Entry.Name()))
		handleError(err)
		spritesheetMetas, err := ReadSpritesheetsData(data)
		handleError(err)
		spritesheetDatas := make([]codec.SpritesheetData,
			0, len(spritesheetMetas))

		for i := 0; i < len(spritesheetMetas); i++ {
			spritesheetMeta := spritesheetMetas[i]
			spritesheetDatas = append(spritesheetDatas,
				spritesheetMeta.ToSpritesheetData())
		}

		for i := 0; i < len(spritesheetDatas); i++ {
			spritesheetMeta := spritesheetMetas[i]
			spritesheetData := spritesheetDatas[i]

			err = resourceFile.Update(func(tx *bolt.Tx) error {
				buck, err := tx.CreateBucketIfNotExists([]byte("spritesheets"))

				if err != nil {
					return err
				}

				if buck == nil {
					return fmt.Errorf("no spritesheets bucket present")
				}

				textureBytes, err := spritesheetData.ToBytes()

				if err != nil {
					return err
				}

				err = buck.Put([]byte(spritesheetMeta.Name), textureBytes)

				if err != nil {
					return err
				}

				return nil
			})
			handleError(err)
		}
	}
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
