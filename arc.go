package main

import (
	_ "embed"
	"github.com/wii-tools/arclib"
	"github.com/wii-tools/lzx/lz77"
)

//go:embed temp.u8
var tempArc []byte

//go:embed templates/sound.bin
var sound []byte

func (b *Banner) makeArcFile(app *App) ([]byte, error) {
	arc, err := arclib.Load(tempArc)
	if err != nil {
		return nil, err
	}

	metaDir, err := arc.OpenDir("meta")
	if err != nil {
		return nil, err
	}

	// Create the banner
	_bannerArc, err := app.makeBannerBin()
	if err != nil {
		return nil, err
	}

	b.IMET.FileSizes.BannerSize = uint32(len(_bannerArc))
	compressed, err := lz77.Compress(_bannerArc)
	if err != nil {
		return nil, err
	}

	imd5, err := makeIMD5(compressed)
	if err != nil {
		return nil, err
	}

	metaDir.WriteFile("banner.bin", append(imd5, compressed...))

	// Create the Icon
	_iconArc, err := app.makeIconArchive()
	if err != nil {
		return nil, err
	}

	b.IMET.FileSizes.IconSize = uint32(len(_iconArc))
	compressed, err = lz77.Compress(_iconArc)
	if err != nil {
		return nil, err
	}

	imd5, err = makeIMD5(compressed)
	if err != nil {
		return nil, err
	}

	metaDir.WriteFile("icon.bin", append(imd5, compressed...))

	imd5, err = makeIMD5(sound)
	if err != nil {
		return nil, err
	}

	metaDir.WriteFile("sound.bin", append(imd5, sound...))

	b.IMET.FileSizes.SoundSize = uint32(len(sound))

	return arc.Save()
}
