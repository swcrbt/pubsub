package dispatcher

type dispatcher struct {
	data chan []byte
}

func New() *dispatcher {
	return &dispatcher{
		data: make(chan []byte, 1000),
	}
}

func (dis *dispatcher) Push (data []byte)  {
	dis.data <- data
}