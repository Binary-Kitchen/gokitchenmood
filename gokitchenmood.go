package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/binary-kitchen/gokitchenmood/lampen"
	"github.com/carbocation/interpose"
	"github.com/gorilla/mux"
)

var uploadalert string
var filetowrite string
var limit int
var mu = &sync.Mutex{}
var url = "http://127.0.0.1/api/"

func handler(w http.ResponseWriter, r *http.Request) {
	title := "moodlights"
	p := &lampen.Lampen{}
	err := p.LoadLampValues(title)
	if err != nil {
		log.Printf("Error loading Config File")
		for i := range p.Values {
			p.Values[i] = "000000"

		}
	}
	t, _ := template.ParseFiles("templates/template.html")
	t.Execute(w, p)

}

func savehandler(w http.ResponseWriter, r *http.Request) {
	p := &lampen.Lampen{}
	for i := 0; i < 10; i++ {
		err := p.Parse(r.FormValue("Lampe"+strconv.Itoa(i)), i)
		if err != nil {
			t, _ := template.ParseFiles("templates/error.html")
			t.Execute(w, p)
			return
		}
	}

	b, err := json.Marshal(p)
	if err != nil {
		fmt.Println(err)
		return
	}

	req, err := http.NewRequest("POST", url+"lampen", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Status Code not ok")
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("response Body:", string(body))
		return
	}

	t, _ := template.ParseFiles("templates/success.html")
	t.Execute(w, p)
}

func sethandler(w http.ResponseWriter, r *http.Request) {
	p := &lampen.Lampen{}
	color := r.FormValue("color")
	for i, _ := range p.Values {
		p.Values[i] = color
	}
	b, err := json.Marshal(p)
	if err != nil {
		fmt.Println(err)
		return
	}

	req, err := http.NewRequest("POST", url+"lampen", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Status Code not ok")
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("response Body:", string(body))
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)

}

func randomhandler(w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest("POST", url+"lampen/random", nil)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Status Code not ok")
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("response Body:", string(body))
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)

}

func isauthorized(header string) bool {
	switch header {
	case "xaver":
		return true
	default:
		return false
	}
}

func authMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			isauth := isauthorized(r.Header.Get("X-BinaryKitchen-Login"))
			if isauth {
				next.ServeHTTP(w, r)
			} else {
				mu.Lock()
				if limit > 0 {
					limit = limit - 1
					fmt.Println("Limit:", limit)
					next.ServeHTTP(w, r)
				} else {
					http.Error(w, "Rate Limit Reached", http.StatusForbidden)
				}
				mu.Unlock()
			}
		})
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s File [-f] \n", os.Args[0])
		return
	}
	limit = 10
	ticker := time.NewTicker(time.Minute)
	go func() {
		for _ = range ticker.C {
			mu.Lock()
			if limit < 10 {
				limit = 10
			}
			mu.Unlock()
		}
	}()
	lampen.Setup()
	lampen.File = false
	filetowrite = os.Args[1]
	if len(os.Args) > 2 {
		if os.Args[2] == "-f" {
			lampen.File = true
		}
	}

	lampen.Port = filetowrite

	lampen.Lampe = lampen.Lampen{}
	lampen.Lampe.SetLampstosavedValues("moodlights")

	middle := interpose.New()

	r := mux.NewRouter()
	r = r.StrictSlash(true)
	r.HandleFunc("/", handler)
	r.HandleFunc("/save", savehandler)
	r.HandleFunc("/set", sethandler)
	r.HandleFunc("/random", randomhandler)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	api := mux.NewRouter().PathPrefix("/api").Subrouter()
	api = api.StrictSlash(true)
	api.HandleFunc("/lampen", lampen.GetLampsHandler).Methods("GET")
	api.HandleFunc("/lampen", lampen.PostLampsHandler).Methods("POST")
	api.HandleFunc("/lampen/random", lampen.PostLampsRandomHandler).Methods("POST")

	middle.Use(authMiddleware())
	middle.UseHandler(api)
	r.PathPrefix("/api").Handler(middle)

	go func() {
		log.Fatal(http.ListenAndServe(":80", r)) // dual stack
	}()
	log.Fatal(http.ListenAndServeTLS(":443", "cert.pem", "private.key", r))
}
