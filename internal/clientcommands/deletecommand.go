package clientcommands

import (
	"context"

	"github.com/bernata/kvstore/apiclient"
)

type DeleteCommand struct {
	Key string `help:"key to delete" required:""`
}

func (c *DeleteCommand) Run(commandContext *CommandContext) error {
	client := commandContext.Client

	response, err := client.DeleteKey(context.Background(), &apiclient.DeleteKeyRequest{Key: c.Key})
	if err != nil {
		return err
	}

	s, err := encodeJSON(response)
	if err != nil {
		return err
	}

	_, _ = commandContext.OutputWriter.WriteString(s + "\n")
	return nil
}
