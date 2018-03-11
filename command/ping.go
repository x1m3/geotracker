package command

func Ping(req Request) (Response, error) {
	return "pong", nil
}
