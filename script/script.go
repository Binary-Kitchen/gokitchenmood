package script

import (
	"fmt"
	"gokitchenmood/lampen"
	"io/ioutil"
	"github.com/stevedonovan/luar"
)

var p *lampen.Lampen = &lampen.Lampen{}
var run bool = true

func setLamps(id int, value string) error {
	return p.Parse(value, id)
}

func sendLamps() {
	p.Send()
}

func resetLamps() {
	for i, _ := range p.Values {
		p.Values[i] = "000000"
	}
}

func repeatScript(content string) {
	L := luar.Init()
	defer L.Close()

	// arbitrary Go functions can be registered
	// to be callable from Lua
	luar.Register(L, "", luar.Map{
		"setLamps":   setLamps,
		"sendLamps":  sendLamps,
		"resetLamps": resetLamps,
	})
	for {
		fmt.Println("blub")
		res := L.DoString(string(content))
		if res != nil {
			run = false
		}
		if run == false {
			fmt.Println("waswfdqwef")
			break
			run = true
		}
	}
	//return nil
}

func RunScript(script string) error {
	for !run {
	}
	fmt.Println("===============================================================================================")
	contents, err := ioutil.ReadFile("uploaded/" + script)
	if err != nil {
		return err
	}
	fmt.Println(string(contents))

	go repeatScript(string(contents))

	return nil
}
