package durchreiche

import "errors"

const preamble = "01000000"

type packet struct {
	Source      string
	Destination string
	Lenght      string
	Payload     string
}

func (p *packet) Send(filename string) error {
	err := errors.New("Not yet implemented")
	return err
}
