package command

type Ping struct {}

func NewPing() *Ping {
	return &Ping{}
}

func (c *Ping) Call(req Request) (Response, error) {
	return "pong", nil
}
