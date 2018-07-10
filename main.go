package main

import (
	"flag"
	"time"
	"math/rand"
	"gitlab.com/Kpovoc/chat-steward/src/startup"
	"gitlab.com/Kpovoc/chat-steward/src/plugin"
)

func main() {
	rd := flag.String(
		"resources",
		"resources",
		"directory containing resources needed for application")
	udd := flag.String(
		"userdata",
		"data",
		"directory containing config file and database")
	flag.Parse()
	resourceDir := *rd
	userDataDir := *udd

	// Go ahead and seed rand for plugins that need it.
	rand.Seed(time.Now().UTC().UnixNano())

	// Initialize the application config from file.
	startup.InitializeConfig(userDataDir)

	// Initialize Database
	startup.InitializeDatabase(userDataDir)

	// Initialize application plugins
	plugin.InitPlugins(resourceDir)

	// Launch all bots
	startup.LaunchBots(resourceDir)
}
