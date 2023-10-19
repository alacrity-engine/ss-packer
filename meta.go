package main

import codec "github.com/alacrity-engine/resource-codec"

// TODO: make spritesheets act
// on a part of a texture, not
// the whole texture (add orig,
// pixel area width and height
// fields).

// SpritesheetMeta is spritesheet metadata
// read from the YAML file.
type SpritesheetMeta struct {
	Name   string `yaml:"name"`
	Width  int    `yaml:"width"`
	Height int    `yaml:"height"`
}

func (meta SpritesheetMeta) ToSpritesheetData() codec.SpritesheetData {
	return codec.SpritesheetData{
		Width:  int32(meta.Width),
		Height: int32(meta.Height),
	}
}
