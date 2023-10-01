package main

import "gopkg.in/yaml.v2"

func ReadSpritesheetsData(data []byte) ([]SpritesheetMeta, error) {
	spritesheets := make([]SpritesheetMeta, 0)
	err := yaml.Unmarshal(data, &spritesheets)

	if err != nil {
		return nil, err
	}

	return spritesheets, nil
}
