package handler

type Handler interface {
	Authenticate(string) error
	CreateTinyUrl()
	GetOriginalUrl()
	//DeleteUrl(string, string) middleware.Responder
}
