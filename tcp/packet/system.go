package packet

type Handler func(data []byte)
type ErrorHandler func(err error)

type System struct {
	Handlers     map[byte]Handler
	ErrorHandler ErrorHandler
}

func NewPacketSystem() *System {
	s := System{
		Handlers: map[byte]Handler{},
	}

	return &s
}

func (s *System) SetErrorHandler(handler ErrorHandler) {
	s.ErrorHandler = handler
}
