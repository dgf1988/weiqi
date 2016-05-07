package mux

type Handler interface {
	ServeHTTP(*Http)
}

type HandlerFunc func(*Http)

func (hf HandlerFunc) ServeHTTP(h *Http) {
	hf(h)
}
