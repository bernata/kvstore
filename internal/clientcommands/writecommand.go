package clientcommands

import (
	"context"

	"github.com/bernata/kvstore/apiclient"
)

type WriteCommand struct {
	Key   string `help:"key to write value for" required:""`
	Value string `help:"value to write" required:""`
}

func (c *WriteCommand) Run(commandContext *CommandContext) error {
	client := commandContext.Client

	response, err := client.WriteValue(context.Background(), &apiclient.WriteValueRequest{Key: c.Key, Value: c.Value})
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
