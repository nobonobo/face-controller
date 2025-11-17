package service

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"machine/usb/hid/joystick"

	"github.com/nobonobo/gamepad-emulator/jsonrpc"
)

type method func(params map[string]any) (any, error)

type JoyStick struct {
	methods map[string]method
}

func New() *JoyStick {
	j := &JoyStick{}
	j.methods = map[string]method{
		"Button": func(params map[string]any) (any, error) {
			arg, ok := params["index"]
			if !ok {
				return nil, fmt.Errorf("missing argument: index")
			}
			v, ok := arg.(float64)
			if !ok {
				return nil, fmt.Errorf("invalid argument: index")
			}
			return js.Button(int(v)), nil
		},
		"SetButton": func(params map[string]any) (any, error) {
			arg1, ok := params["index"]
			if !ok {
				return nil, fmt.Errorf("missing argument: index")
			}
			v1, ok := arg1.(float64)
			if !ok {
				return nil, fmt.Errorf("invalid argument: index")
			}
			arg2, ok := params["push"]
			if !ok {
				return nil, fmt.Errorf("missing argument: push")
			}
			v2, ok := arg2.(bool)
			if !ok {
				return nil, fmt.Errorf("invalid argument: push")
			}
			js.SetButton(int(v1), v2)
			return true, nil
		},
		"Hat": func(params map[string]any) (any, error) {
			arg, ok := params["index"]
			if !ok {
				return nil, fmt.Errorf("missing argument: index")
			}
			v, ok := arg.(float64)
			if !ok {
				return nil, fmt.Errorf("invalid argument: index")
			}
			return js.Hat(int(v)), nil
		},
		"SetHat": func(params map[string]any) (any, error) {
			arg1, ok := params["index"]
			if !ok {
				return nil, fmt.Errorf("missing argument: index")
			}
			v1, ok := arg1.(float64)
			if !ok {
				return nil, fmt.Errorf("invalid argument: index")
			}
			arg2, ok := params["dir"]
			if !ok {
				return nil, fmt.Errorf("missing argument: dir")
			}
			v2, ok := arg2.(float64)
			if !ok {
				return nil, fmt.Errorf("invalid argument: dir")
			}
			js.SetHat(int(v1), joystick.HatDirection(uint8(v2)))
			return true, nil
		},
		"Axis": func(params map[string]any) (any, error) {
			arg, ok := params["index"]
			if !ok {
				return nil, fmt.Errorf("missing argument: index")
			}
			v, ok := arg.(float64)
			if !ok {
				return nil, fmt.Errorf("invalid argument: index")
			}
			return js.Axis(int(v)), nil
		},
		"SetAxis": func(params map[string]any) (any, error) {
			arg1, ok := params["index"]
			if !ok {
				return nil, fmt.Errorf("missing argument: index")
			}
			v1, ok := arg1.(float64)
			if !ok {
				return nil, fmt.Errorf("invalid argument: index")
			}
			arg2, ok := params["value"]
			if !ok {
				return nil, fmt.Errorf("missing argument: value")
			}
			v2, ok := arg2.(float64)
			if !ok {
				return nil, fmt.Errorf("invalid argument: value")
			}
			js.SetAxis(int(v1), int(v2))
			return true, nil
		},
		"SendState": func(params map[string]any) (any, error) {
			js.SendState()
			return true, nil
		},
	}
	return j
}

func (j *JoyStick) Run(conn io.ReadWriteCloser) error {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		req := new(jsonrpc.Request)
		if err := req.UnmarshalJSON([]byte(line)); err != nil {
			return fmt.Errorf("req unmarshal failed: %w", err)
		}
		resp := &jsonrpc.Response{
			ID: req.ID,
		}
		r, err := j.methods[req.Method](req.Params)
		if err != nil {
			resp.Error = &jsonrpc.Error{
				Code:    -32603,
				Message: err.Error(),
			}
		} else {
			resp.Result = r
		}
		respBytes, err := resp.MarshalJSON()
		if err != nil {
			return fmt.Errorf("resp marshal failed: %w", err)
		}
		if _, err := conn.Write(append(respBytes, '\n')); err != nil {
			return fmt.Errorf("write failed: %w", err)
		}
	}
	return nil
}
