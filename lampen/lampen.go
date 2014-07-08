package lampen

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"regexp"
)

var validColor = regexp.MustCompile(`^#([A-Fa-f0-9]{6})|([A-Fa-f0-9]{6})$`)

type Lampen struct {
	Values [10]string
}

func (l *Lampen) Send() error {
	err := errors.New("Not yet implemented")
	return err
}

func (l *Lampen) Parse(input string, number int) error {
	if validColor.MatchString(input) {
		l.Values[number] = input
		return nil
	} else {
		err := errors.New("String does not match Color Regex")
		return err
	}
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
