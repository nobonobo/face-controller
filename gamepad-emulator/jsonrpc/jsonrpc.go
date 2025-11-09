package jsonrpc

//go:generate easyjson -all jsonrpc.go

type Request struct {
	ID      int            `json:"id"`
	JsonRpc string         `json:"jsonrpc"`
	Method  string         `json:"method"`
	Params  map[string]any `json:"params,omitempty"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Response struct {
	ID     int `json:"id"`
	Result any `json:"result,omitempty"`
	Error  any `json:"error,omitempty"`
}
