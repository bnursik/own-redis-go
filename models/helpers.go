package models

import "fmt"

func HelpMessage() {
	msg := `Own Redis

	Usage:
	  own-redis [--port <N>]
	  own-redis --help
	
	Options:
	  --help       Show this screen.
	  --port N     Port number.`
	fmt.Println(msg)
}
