package main

import (
	"flag"
	"gitlab.com/Kpovoc/chat-steward/internal/app/startup"
	"math/rand"
	"time"
)

func main() {
	rd := flag.String(
		"resources",
		".",
		"directory containing resources needed for application")
	udd := flag.String(
		"userdata",
		"configs",
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
	startup.InitializePlugins(resourceDir)

	// Launch all bots
	startup.LaunchBots(resourceDir)
}
