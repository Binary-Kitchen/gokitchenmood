package lampen

import (
	"encoding/json"
	"errors"
	"gokitchenmood/durchreiche"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/ant0ine/go-json-rest/rest"
)

const controlleradress byte = 0x10 //0x10
const clientadress byte = 0xFE     //0xFE
const payloadlength byte = 0x1E    //30 da immer 3 byte pro Lampe a 10 Lampen
var validColor = regexp.MustCompile(`^#([A-Fa-f0-9]{6})|([A-Fa-f0-9]{6})$`)
var File bool
var Port string

type Lampen struct {
	Values [10]string
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
	p.Send(Port, File)
	//err := errors.New("wa")
}

func (l *Lampen) Parse(input string, number int) error {
	var reterr error
	if validColor.MatchString(input) {
		l.Values[number] = input
		reterr = nil
	} else {
		err := errors.New("String does not match Color Regex")
		reterr = err
	}
	return reterr
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

func (l *Lampen) SetLampstosavedValues(filename string) error {
	err := l.LoadLampValues(filename)
	if err == nil {
		l.Send()
		return nil
	} else {
		return err
	}
}

func (l *Lampen) GetAllLamps(w rest.ResponseWriter, r *rest.Request) {
	err := l.LoadLampValues("moodlights")
	if err == nil {
		w.WriteJson(&l)
	} else {
		rest.Error(w, err.Error(), http.StatusInternalServerError)

	}
}

func (l *Lampen) PostLamps(w rest.ResponseWriter, r *rest.Request) {
	err := r.DecodeJsonPayload(&l)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	l.Send()
}
