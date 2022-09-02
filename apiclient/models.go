package apiclient

type RetrieveValueRequest struct {
	Key string
}

type RetrieveValueResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type WriteValueRequest struct {
	Key   string `json:"-"`
	Value string `json:"value"`
}

type WriteValueResponse struct {
}

type DeleteKeyRequest struct {
	Key string `json:"key"`
}

type DeleteKeyResponse struct {
}
