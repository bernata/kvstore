package main

import (
	"fmt"
	"os"

	"github.com/bernata/kvstore/apiclient"

	"github.com/alecthomas/kong"
	"github.com/bernata/kvstore/internal/clientcommands"
)

type CommandLine struct {
	URL  string `help:"endpoint to query" short:"u" default:"http://localhost"`
	Port int    `help:"Port to listen on" short:"p" default:"8282"`

	Get    clientcommands.GetCommand    `cmd:"" help:"Get key"`
	Delete clientcommands.DeleteCommand `cmd:"" help:"Delete key-value"`
	Write  clientcommands.WriteCommand  `cmd:"" help:"Write key-value"`
}

func main() {
	var commandLine CommandLine
	command := kong.Parse(&commandLine)

	client, err := newClient(&commandLine)
	if err != nil {
		command.Errorf("error creating client %s", err)
		os.Exit(255)
	}

	err = command.Run(&clientcommands.CommandContext{
		Client:       client,
		OutputWriter: os.Stdout,
	})
	command.FatalIfErrorf(err)
}

func newClient(commandLine *CommandLine) (*apiclient.Client, error) {
	u := fmt.Sprintf("%s:%d", commandLine.URL, commandLine.Port)
	return apiclient.NewClient(apiclient.BaseURL(u)), nil
}
