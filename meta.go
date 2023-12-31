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
	Name   string   `yaml:"name"`
	Width  int      `yaml:"width"`
	Height int      `yaml:"height"`
	Orig   OrigMeta `taml:"orig"`
	Area   AreaMeta `yaml:"area"`
}

type OrigMeta struct {
	X int `yaml:"x"`
	Y int `yaml:"y"`
}

type AreaMeta struct {
	PixelWidth  int `yaml:"pixelWidth"`
	PixelHeight int `yaml:"pixelHeight"`
}

func (meta SpritesheetMeta) ToSpritesheetData() codec.SpritesheetData {
	return codec.SpritesheetData{
		Width:  int32(meta.Width),
		Height: int32(meta.Height),
		Orig: codec.OrigData{
			X: int32(meta.Orig.X),
			Y: int32(meta.Orig.Y),
		},
		Area: codec.AreaData{
			PixelWidth:  int32(meta.Area.PixelWidth),
			PixelHeight: int32(meta.Area.PixelHeight),
		},
	}
}
