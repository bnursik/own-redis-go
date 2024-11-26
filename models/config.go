package models

import "flag"

var (
	PortFlag = flag.String("port", "8080", "Port for th HTTP server")
	HelpFlag = flag.Bool("help", false, "provides usage information")
)
