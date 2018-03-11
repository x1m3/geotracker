package command

type Request map[string]interface{}
type Response interface{}

type Command func(req Request) (Response, error)
