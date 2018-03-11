package server

import (
	"io"
	"github.com/x1m3/geotracker/command"
)

type ProtocolAdapter interface {
	// Encodes a command.Response with specific format and writes it into resp.
	Encode(resp io.Writer, item command.Response) error
	// Reads bytes from read and returns a command.Request
	Decode(read io.Reader) (command.Request, error)
	// Return a mime content-type of the protocol, like application/json or application/xml
	ContentType() string
}
