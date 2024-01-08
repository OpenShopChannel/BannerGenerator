package main

import _ "embed"

//go:embed templates/demos.u8
var bannerTemplateArc []byte

//go:embed templates/badges/utilities.tpl
var utilitiesBadge []byte

//go:embed templates/badges/games.tpl
var gamesBadge []byte

//go:embed templates/badges/demos.tpl
var demosBadge []byte

//go:embed templates/badges/emulators.tpl
var emulatorsBadge []byte

//go:embed templates/badges/media.tpl
var mediaBadge []byte

var badges = map[string][]byte{
	"utilities": utilitiesBadge,
	"games":     gamesBadge,
	"demos":     demosBadge,
	"emulators": emulatorsBadge,
	"media":     mediaBadge,
}
