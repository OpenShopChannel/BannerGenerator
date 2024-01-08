package main

import (
	_ "embed"
	"github.com/wii-tools/arclib"
)

//go:embed templates/icon.u8
var iconArc []byte

func (a *App) makeIconArchive() ([]byte, error) {
	arc, err := arclib.Load(iconArc)
	if err != nil {
		return nil, err
	}

	arcDir, err := arc.OpenDir("arc")
	if err != nil {
		return nil, err
	}

	// Image files
	imageDir, err := arcDir.GetDir("timg")
	if err != nil {
		return nil, err
	}

	// Write the icon
	icon, err := a.GetIcon()
	if err != nil {
		return nil, err
	}

	imageDir.WriteFile("D_applogo.tpl", icon)

	return arc.Save()
}
