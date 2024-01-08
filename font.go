package main

import (
	_ "embed"
	"fmt"
	"github.com/golang/freetype"
	"github.com/wii-tools/libtpl"
	"image"
	"image/draw"
	"image/png"
	"os"
	"time"
)

//go:embed fonts/Kanit-Light.ttf
var kanitLight []byte

//go:embed fonts/Kanit-ExtraLight.ttf
var kanitExtraLight []byte

func DrawAppName(appName string) ([]byte, error) {
	f, err := freetype.ParseFont(kanitLight)
	if err != nil {
		return nil, err
	}

	dst := image.NewRGBA(image.Rect(0, 0, 335, 35))
	draw.Draw(dst, dst.Bounds(), image.Transparent, image.Point{}, draw.Src)

	c := freetype.NewContext()
	c.SetDst(dst)
	c.SetClip(dst.Bounds())
	c.SetSrc(image.White)
	c.SetFont(f)
	c.SetFontSize(35)

	_, err = c.DrawString(appName, freetype.Pt(0, 25))
	if err != nil {
		return nil, err
	}

	return libtpl.ToIA4(dst)
}

func DrawVersion(version string) ([]byte, error) {
	f, err := freetype.ParseFont(kanitExtraLight)
	if err != nil {
		return nil, err
	}

	dst := image.NewRGBA(image.Rect(0, 0, 417, 28))
	draw.Draw(dst, dst.Bounds(), image.Transparent, image.Point{}, draw.Src)

	c := freetype.NewContext()
	c.SetDst(dst)
	c.SetClip(dst.Bounds())
	c.SetSrc(image.White)
	c.SetFont(f)
	c.SetFontSize(35)

	_, err = c.DrawString(fmt.Sprintf("Version: %s", version), freetype.Pt(0, 25))
	if err != nil {
		return nil, err
	}

	return libtpl.ToIA4(dst)
}

func DrawDescription(description string) ([]byte, error) {
	f, err := freetype.ParseFont(kanitExtraLight)
	if err != nil {
		return nil, err
	}

	dst := image.NewRGBA(image.Rect(0, 0, 417, 28))
	draw.Draw(dst, dst.Bounds(), image.Transparent, image.Point{}, draw.Src)

	c := freetype.NewContext()
	c.SetDst(dst)
	c.SetClip(dst.Bounds())
	c.SetSrc(image.White)
	c.SetFont(f)
	c.SetFontSize(30)

	_, err = c.DrawString(description, freetype.Pt(0, 25))
	if err != nil {
		return nil, err
	}

	return libtpl.ToIA4(dst)
}

func DrawReleaseDate(release int64) ([]byte, error) {
	releaseStr := "Not available"
	if release != 0 {
		unix := time.Unix(release, 0)
		releaseStr = unix.Format("2006-01-02")
	}

	f, err := freetype.ParseFont(kanitExtraLight)
	if err != nil {
		return nil, err
	}

	dst := image.NewRGBA(image.Rect(0, 0, 417, 28))
	draw.Draw(dst, dst.Bounds(), image.Transparent, image.Point{}, draw.Src)

	c := freetype.NewContext()
	c.SetDst(dst)
	c.SetClip(dst.Bounds())
	c.SetSrc(image.White)
	c.SetFont(f)
	c.SetFontSize(30)

	_, err = c.DrawString(fmt.Sprintf("Release Date: %s", releaseStr), freetype.Pt(0, 25))
	if err != nil {
		return nil, err
	}

	h, err := os.Create("help.png")
	if err != nil {
		return nil, err
	}

	err = png.Encode(h, dst)
	if err != nil {
		return nil, err
	}

	return libtpl.ToIA4(dst)
}
