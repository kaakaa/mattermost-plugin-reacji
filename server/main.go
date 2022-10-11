package main

import (
	mmplugin "github.com/mattermost/mattermost-server/v6/plugin"

	"github.com/kaakaa/mattermost-plugin-reacji/server/plugin"
)

func main() {
	mmplugin.ClientMain(&plugin.Plugin{})
}
