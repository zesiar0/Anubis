package main

import "Anubis/anubis"

var (
	configPath = "/Users/a3bz/GolandProjects/Anubis/pkg/config/config.yml"
)

func main() {
	anubis.Run(configPath)
}
