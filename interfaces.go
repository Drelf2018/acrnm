package acrnm

type Sender interface {
	Send(*Product) error
}

type Senders []Sender

func (s Senders) Send(p *Product) error {
	for _, sender := range s {
		go sender.Send(p)
	}
	return nil
}
