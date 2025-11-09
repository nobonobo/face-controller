package main

import (
	"encoding/json"
	"fmt"
	"time"

	"go.bug.st/serial"
)

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

type JoyStickService struct {
	port    serial.Port
	encoder *json.Encoder
	decoder *json.Decoder
	id      int
}

func NewJoyStickService(port string) (*JoyStickService, error) {
	p, err := serial.Open(port, &serial.Mode{
		BaudRate: 12000000,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	})
	if err != nil {
		return nil, err
	}
	p.SetReadTimeout(3 * time.Second)
	encoder := json.NewEncoder(p)
	decoder := json.NewDecoder(p)
	return &JoyStickService{
		port:    p,
		encoder: encoder,
		decoder: decoder,
		id:      0,
	}, nil
}

func (js *JoyStickService) Close() error {
	return js.port.Close()
}

func (js *JoyStickService) call(method string, params map[string]any) (any, error) {
	js.id++
	if err := js.encoder.Encode(&Request{
		ID:      js.id,
		JsonRpc: "2.0",
		Method:  method,
		Params:  params,
	}); err != nil {
		return nil, err
	}
	var resp Response
	if err := js.decoder.Decode(&resp); err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, fmt.Errorf("%v", resp)
	}
	return resp.Result, nil
}

func (js *JoyStickService) Button(index int) (bool, error) {
	res, err := js.call("Button", map[string]any{
		"index": index,
	})
	if err != nil {
		return false, err
	}
	return res.(bool), nil
}

func (js *JoyStickService) SetButton(index int, push bool) error {
	if _, err := js.call("SetButton", map[string]any{
		"index": index,
		"push":  push,
	}); err != nil {
		return err
	}
	return nil
}

func (js *JoyStickService) Hat(index int) (uint8, error) {
	res, err := js.call("Hat", map[string]any{
		"index": index,
	})
	if err != nil {
		return 0, err
	}
	return uint8(res.(float64)), nil
}

func (js *JoyStickService) SetHat(index int, dir uint8) error {
	if _, err := js.call("SetHat", map[string]any{
		"index": index,
		"dir":   dir,
	}); err != nil {
		return err
	}
	return nil
}

func (js *JoyStickService) Axis(index int) (int, error) {
	res, err := js.call("Axis", map[string]any{
		"index": index,
	})
	if err != nil {
		return 0, err
	}
	return int(res.(float64)), nil
}

func (js *JoyStickService) SetAxis(index int, v int) error {
	if _, err := js.call("SetAxis", map[string]any{
		"index": index,
		"value": v,
	}); err != nil {
		return err
	}
	return nil
}

func (js *JoyStickService) SendState() error {
	if _, err := js.call("SendState", nil); err != nil {
		return err
	}
	return nil
}
