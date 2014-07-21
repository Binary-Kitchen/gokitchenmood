package lampen

import (
	"encoding/json"
	"errors"
	"gokitchenmood/durchreiche"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

const controlleradress byte = 0x10 //0x10
const clientadress byte = 0xFE     //0xFE
const payloadlength byte = 0x1E    //30 da immer 3 byte pro Lampe a 10 Lampen
var validColor = regexp.MustCompile(`^#([A-Fa-f0-9]{6})|([A-Fa-f0-9]{6})$`)
var File bool

type Lampen struct {
	Values [10]string
	Port   string
}

func (l *Lampen) Send() {
	p := &durchreiche.Packet{}
	for k, s := range l.Values {
		news := strings.Replace(s, "#", "", -1)
		for j := 0; j < 6; j += 2 {
			i, _ := strconv.ParseUint(news[j:j+2], 16, 8)
			b := byte(i)
			p.Payload[k*3+(j/2)] = b
		}

	}
	//fmt.Println(p.Payload)
	//fmt.Println(len(payload))
	//p.Payload = payload
	p.Source = clientadress
	p.Destination = controlleradress
	p.Length = payloadlength
	p.Send(l.Port, File)
	//err := errors.New("wa")
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
