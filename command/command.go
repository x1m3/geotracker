package command

type Request map[string]interface{}
type Response interface{}

type Command interface {
	Call(Request) (Response, error)
}
