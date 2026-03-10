package main

import "github.com/TarelX/TCLI/cmd"

var (
	version   = "dev"
	buildDate = "unknown"
)

func main() {
	cmd.Execute(version, buildDate)
}
