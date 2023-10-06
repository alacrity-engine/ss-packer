package main

import (
	"flag"
	"fmt"
	"os"

	codec "github.com/alacrity-engine/resource-codec"
	bolt "go.etcd.io/bbolt"
)

var (
	spritesheetsMetaPath string
	resourceFilePath     string
)

func parseFlags() {
	flag.StringVar(&spritesheetsMetaPath, "spritesheets-meta",
		"./spritesheets.yml", "Path to the spritesheets metadata file.")
	flag.StringVar(&resourceFilePath, "out", "./stage.res",
		"Resource file to store animations and spritesheets.")

	flag.Parse()
}

func main() {
	parseFlags()

	data, err := os.ReadFile(spritesheetsMetaPath)
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

	resourceFile, err := bolt.Open(resourceFilePath, 0666, nil)
	handleError(err)
	defer resourceFile.Close()

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

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
