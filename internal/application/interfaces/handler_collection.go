package interfaces

type HandlerCollection map[string]Handler

func (hc HandlerCollection) Add(handler Handler) {
	_, ok := hc[handler.Name()]
	if ok {
		panic(handler.Name() + " already exists")
	}

	hc[handler.Name()] = handler
}
