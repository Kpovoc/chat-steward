package core

import "fmt"

func ExampleF_parseContent() {
	pluginName, args := parseContent("!Test this thing out!")
	fmt.Println(pluginName)
	for _, str := range args {
		fmt.Println(str)
	}
	// Output:
	// Test
	// this
	// thing
	// out!
}