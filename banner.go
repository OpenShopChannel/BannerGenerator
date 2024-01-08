package main

import (
	"bytes"
	_ "embed"
	"encoding/binary"
	"fmt"
	"github.com/wii-tools/arclib"
	"os"
)

const Outpath = "out/%s"

type Banner struct {
	Padding [64]byte
	IMET    IMET
}

func makeBanner(app *App) error {
	b := Banner{}
	b.makeIMET(app.Name)
	arcData, err := b.makeArcFile(app)
	if err != nil {
		return err
	}

	err = b.CalculateIMETMD5()
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, &b)
	if err != nil {
		return err
	}

	buf.Write(arcData)

	err = os.MkdirAll(fmt.Sprintf(Outpath, app.Shop.TitleID), 0777)
	if err != nil {
		return err
	}

	return os.WriteFile(fmt.Sprintf(Outpath+"/00000000.app", app.Shop.TitleID), buf.Bytes(), 0666)
}

func (a *App) makeBannerBin() ([]byte, error) {
	arc, err := arclib.Load(bannerTemplateArc)
	if err != nil {
		return nil, err
	}

	arcDir, err := arc.OpenDir("arc")
	if err != nil {
		return nil, err
	}

	// Open the image directory for writing
	imageDir, err := arcDir.GetDir("timg")
	if err != nil {
		return nil, err
	}

	// First change the badge
	imageDir.WriteFile("C_badge.tpl", badges[a.Category])

	// Write the icon
	icon, err := a.GetIcon()
	if err != nil {
		return nil, err
	}

	imageDir.WriteFile("D_applogo.tpl", icon)

	// Write filename
	name, err := DrawAppName(a.Name)
	if err != nil {
		return nil, err
	}

	imageDir.WriteFile("D_appname.tpl", name)

	// Write version
	version, err := DrawVersion(a.Version)
	if err != nil {
		return nil, err
	}

	imageDir.WriteFile("D_author_version.tpl", version)

	// Write description
	description, err := DrawDescription(a.Description.Short)
	if err != nil {
		return nil, err
	}

	imageDir.WriteFile("D_shortdesc.tpl", description)

	// Finally the release date
	release, err := DrawReleaseDate(a.ReleaseDate)
	if err != nil {
		return nil, err
	}

	imageDir.WriteFile("D_files_reldate.tpl", release)

	return arc.Save()
}
