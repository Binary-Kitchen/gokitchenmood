package lampen

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/binary-kitchen/gokitchenmood/durchreiche"
)

const controlleradress byte = 0x10 //0x10
const clientadress byte = 0xFE     //0xFE
const payloadlength byte = 0x1E    //30 da immer 3 byte pro Lampe a 10 Lampen
const gamma = 2.8

var validColor = regexp.MustCompile(`^#([A-Fa-f0-9]{6})|([A-Fa-f0-9]{6})$`)
var File bool
var Port string
var correction [256]int

var Lampe Lampen

type Lampen struct {
	Values [10]string
}

func Setup() {
	for i := range correction {
		correction[i] = (int)(math.Pow(float64(i)/float64(255), gamma)*255 + 0.5)
		//fmt.Println(correction[i])
	}
}

func (l *Lampen) Send() {
	p := &durchreiche.Packet{}
	for k, s := range l.Values {
		//	fmt.Println(s)
		newS := strings.Replace(s, "#", "", -1)
		for j := 0; j < 6; j += 2 {
			i, _ := strconv.ParseUint(newS[j:j+2], 16, 8)
			b := byte(correction[i])
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
	if number > 10 {
		err := errors.New("ID is not in Range")
		return err
	}

	if !validColor.MatchString(input) {
		err := errors.New("String does not match Color Regex")
		return err
	}

	l.Values[number] = input
	return nil

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
	if err != nil {
		return err
	}
	l.Send()
	return nil
}

func strtohex(color int64) string {
	news := strconv.FormatInt(color, 16)
	if len(news) < 2 {
		news = "0" + news
	}
	return news

}

func (l *Lampen) SetRandom() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i, _ := range l.Values {
		colorr := r.Int63n(255)
		colorg := r.Int63n(255)
		colorb := r.Int63n(255)
		l.Values[i] = strtohex(colorr)
		l.Values[i] = l.Values[i] + strtohex(colorg)
		l.Values[i] = l.Values[i] + strtohex(colorb)
	}
	l.WriteLampValues("moodlights")
	l.Send()
}

func GetLampsHandler(w http.ResponseWriter, r *http.Request) {
	j, err := json.Marshal(&Lampe)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func PostLampsHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&Lampe)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = Lampe.WriteLampValues("moodlights")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	Lampe.Send()
}

func PostLampsRandomHandler(w http.ResponseWriter, r *http.Request) {
	Lampe.SetRandom()
}
