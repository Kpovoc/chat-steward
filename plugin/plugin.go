package plugin

func GetPluginResponse(pluginName string, args []string) string {
	response := ""

	switch pluginName {
	case "8ball":
		fallthrough
	case "eightball":
		response = EightBall(args)
	}

	return response
}