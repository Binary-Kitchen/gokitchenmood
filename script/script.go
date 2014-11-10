package script

import "gokitchenmood/lampen"

var p *lampen.Lampen = &lampen.Lampen{}
var transfer = make(chan bool, 1)
var repeatstopped = true

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
	/*	L := luar.Init()
		defer L.Close()
		// arbitrary Go functions can be registered
		// to be callable from Lua
		luar.Register(L, "", luar.Map{
			"setLamps":   setLamps,
			"sendLamps":  sendLamps,
			"resetLamps": resetLamps,
		})
		//fmt.Println("Cont:", cont, "More:", more)
		cont := true
		more := true
		for cont {
			select {
			case cont, more = <-transfer:
				if !more {
					cont = false
				}
			default:
				L.DoString(string(content))
				//cont = <-transfer
				//fmt.Println("Cont in:", cont)
			}
		}
		repeatstopped = true
		fmt.Println("Bye Bye")
		//return nil*/
}

func RunScript(script string) error {
	/*	transfer <- false
		close(transfer)
		for !repeatstopped {
			time.Sleep(1 * time.Second)
		}
		contents, err := ioutil.ReadFile("uploaded/" + script)
		if err != nil {
			return err
		}
		//fmt.Println(string(contents))
		transfer = make(chan bool, 1)
		transfer <- true
		repeatstopped = false
		go repeatScript(string(contents))
		//transfer <- true*/
	return nil
}
