package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"gokitchenmood/lampen"
	"gokitchenmood/script"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/fatih/color"
)

var uploadalert string
var filetowrite string

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
	if lampen.HardLimit != 0 {
		p := &lampen.Lampen{}
		var broken bool
		for i := 0; i < 10; i++ {
			err := p.Parse(r.FormValue("Lampe"+strconv.Itoa(i)), i)
			if err != nil {
				broken = true
				break
			}
		}
		if broken {
			t, _ := template.ParseFiles("templates/error.html")
			t.Execute(w, p)
		} else {
			err := p.WriteLampValues("moodlights")
			if err != nil {
				log.Fatal(err)
			}
			p.Send()
			t, _ := template.ParseFiles("templates/success.html")
			t.Execute(w, p)
		}
	} else {
		http.Redirect(w, r, "http://www.lemonparty.org", http.StatusFound)
	}
}

func sethandler(w http.ResponseWriter, r *http.Request) {
	if lampen.HardLimit != 0 {
		p := &lampen.Lampen{}
		color := r.FormValue("color")
		for i, _ := range p.Values {
			p.Values[i] = color
		}
		p.WriteLampValues("moodlights")
		p.Send()
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		http.Redirect(w, r, "http://www.lemonparty.org", http.StatusFound)
	}

}

func statichandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

func randomhandler(w http.ResponseWriter, r *http.Request) {
	if lampen.HardLimit != 0 {
		p := &lampen.Lampen{}
		p.SetRandom()
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		http.Redirect(w, r, "http://www.lemonparty.org", http.StatusFound)
	}
}

func recieveHandler(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file") // the FormFile function takes in the POST input id file
	defer file.Close()

	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	fmt.Println(filepath.Ext(header.Filename))

	if filepath.Ext(header.Filename) != ".lua" {
		uploadalert = "<div class=\"alert alert-danger\" role=\"alert\"><b>Oh snap!</b> File doesn't have .lua extension</div>"
		http.Redirect(w, r, "/upload", http.StatusFound)
		return
	}

	out, err := os.Create("uploaded/" + header.Filename)
	if err != nil {
		//fmt.Fprintf(w, "Unable to create file: %s", err.Error())

		errstring := err.Error()
		uploadalert = "<div class=\"alert alert-danger\" role=\"alert\"><b>Oh snap!</b> Unable to create file: " + errstring + "</div>"
		http.Redirect(w, r, "/upload", http.StatusFound)
		return
	}

	defer out.Close()

	// write the content from POST to the file
	_, err = io.Copy(out, file)
	if err != nil {
		//fmt.Fprintln(w, err)

		errstring := err.Error()
		uploadalert = "<div class=\"alert alert-danger\" role=\"alert\"><b>Oh snap!</b> " + errstring + "</div>"
		http.Redirect(w, r, "/upload", http.StatusFound)
		return
	}

	uploadalert = "<div class=\"alert alert-success\" role=\"alert\"><b>Well done!</b> File uploaded</div>"
	http.Redirect(w, r, "/upload", http.StatusFound)

	fmt.Println(color.GreenString("Info:"), "File uploaded successfully:", header.Filename)

	//fmt.Fprintf(w, "File uploaded successfully : ")
	//fmt.Fprintf(w, header.Filename)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/upload.html")
	err := t.Execute(w, template.HTML(uploadalert))
	if err != nil {
		fmt.Println("There was an error:", err)
	}
	uploadalert = ""
}

func scriptHandler(w http.ResponseWriter, r *http.Request) {

	files, _ := ioutil.ReadDir("uploaded/")
	listing := ""
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".lua" {
			listing = listing + "<option>" + template.HTMLEscapeString(f.Name()) + "</option>"
		}
	}

	t, _ := template.ParseFiles("templates/script.html")
	t.Execute(w, template.HTML(listing))

}

func setscriptHandler(w http.ResponseWriter, r *http.Request) {
	script.RunScript(r.FormValue("scripts"))
	http.Redirect(w, r, "/script", http.StatusFound)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s File [-f] \n", os.Args[0])
		return
	}
	lampen.File = false
	filetowrite = os.Args[1]
	if len(os.Args) > 2 {
		if os.Args[2] == "-f" {
			lampen.File = true
		}
	}

	lampen.Port = filetowrite

	lampe := &lampen.Lampen{}
	lampe.SetLampstosavedValues("moodlights")

	rhandler := rest.ResourceHandler{
		EnableRelaxedContentType: true,
		EnableStatusService:      true,
		XPoweredBy:               "phil-api",
	}
	err := rhandler.SetRoutes(
		rest.RouteObjectMethod("GET", "/lamps", lampe, "GetLamps"),
		rest.RouteObjectMethod("POST", "/lamps", lampe, "PostLamps"),
		&rest.Route{"GET", "/.status",
			func(w rest.ResponseWriter, r *rest.Request) {
				w.WriteJson(rhandler.GetStatus())
			},
		},
		//rest.RouteObjectMethod("GET", "/users/:id", &users, "GetUser"),
		//rest.RouteObjectMethod("PUT", "/users/:id", &users, "PutUser"),
		//rest.RouteObjectMethod("DELETE", "/users/:id", &users, "DeleteUser"),
	)
	if err != nil {
		log.Fatal(err)
	}

	lampen.Limit = 100
	lampen.HardLimit = 200

	ticker := time.NewTicker(time.Second * 10)
	go func() {
		for _ = range ticker.C {
			lampen.Limit++
			if lampen.Limit > 100 {
				lampen.Limit = 100
			}
		}
	}()

	ticker2 := time.NewTicker(time.Second * 10)
	go func() {
		for _ = range ticker2.C {
			lampen.HardLimit++
			if lampen.HardLimit > 200 {
				lampen.HardLimit = 200
			}
		}
	}()

	http.HandleFunc("/", handler)
	http.HandleFunc("/save", savehandler)
	http.HandleFunc("/set", sethandler)
	http.HandleFunc("/random", randomhandler)
	http.HandleFunc("/static/", statichandler)
	http.HandleFunc("/receive", recieveHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/script", scriptHandler)
	http.HandleFunc("/setscript", setscriptHandler)
	http.HandleFunc("/templates", statichandler)
	http.Handle("/api/", http.StripPrefix("/api", &rhandler))
	http.ListenAndServe(":8080", nil)
}
