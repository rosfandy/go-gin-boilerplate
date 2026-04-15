package main

import (
	"fmt"
	"owner-api-proxy/cmd/command"
)

func main() {
	cmd := command.InitCommand()

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		return
	}
}
