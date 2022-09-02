package clientcommands

import (
	"context"

	"github.com/bernata/kvstore/apiclient"
)

type GetCommand struct {
	Key string `help:"key to retrieve value for" required:""`
}

func (c *GetCommand) Run(commandContext *CommandContext) error {
	client := commandContext.Client

	response, err := client.RetrieveValue(context.Background(), &apiclient.RetrieveValueRequest{Key: c.Key})
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
