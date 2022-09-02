package clientcommands

import (
	"encoding/json"
	"io"

	"github.com/bernata/kvstore/apiclient"
)

type CommandContext struct {
	Client       *apiclient.Client
	OutputWriter io.StringWriter
}

func encodeJSON(data interface{}) (string, error) {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
