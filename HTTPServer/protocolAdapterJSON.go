package HTTPServer

import (
	"io"
	"github.com/x1m3/geotracker/command"
	"encoding/json"
)

type JSONAdapter struct{}

func NewJSONAdapter() *JSONAdapter {
	return &JSONAdapter{}
}

func (e *JSONAdapter) Encode(resp io.Writer, item command.Response) error {
	jsonEncoder := json.NewEncoder(resp)
	return jsonEncoder.Encode(item)
}

func (e *JSONAdapter) Decode(read io.Reader) (command.Request, error) {
	response := make(map[string]interface{})
	decoder := json.NewDecoder(read)
	if decoder.More() {
		if err := decoder.Decode(&response); err != nil {
			return nil, err
		}
		return response, nil
	}
	return nil, nil
}

func (e *JSONAdapter) ContentType() string {
	return "application/json"
}
