package lampen

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

type Lampen struct {
	ValueLampe0x1 string
	ValueLampe0x2 string
	ValueLampe0x3 string
	ValueLampe0x4 string
	ValueLampe0x5 string
	ValueLampe0x6 string
	ValueLampe0x7 string
	ValueLampe0x8 string
	ValueLampe0x9 string
	ValueLampe0xA string
}

func (l *Lampen) Send() error {
	err := errors.New("Not yet implemented")
	return err
}

func (l *Lampen) LoadLampValues(filename string) error {
	filename = filename + ".json"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &l)
	if err != nil {
		return err
	}
	return nil
}

func (l *Lampen) WriteLampValues(filename string) error {
	b, err := json.MarshalIndent(&l, "", "    ")
	if err != nil {
		return err
	}
	ioutil.WriteFile(filename+".json", b, 0644)
	if err != nil {
		return err
	}
	return nil
}
