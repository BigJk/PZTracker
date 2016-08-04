package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"

	"encoding/json"

	"golang.org/x/net/websocket"
)

// Event Stuff

type dataEvent struct {
	Name string
	Data map[string]string
}

type basicEvent struct {
	Name string
}

type tracker struct {
}

func (t tracker) Write(p []byte) (n int, err error) {
	s := string(p)

	if trackerEventRegex.MatchString(s) {
		eventName := trackerEventRegex.FindStringSubmatch(s)

		if strings.Contains(s, "(") {
			values := trackerDataRegex.FindAllStringSubmatch(s, -1)
			event := dataEvent{eventName[1], make(map[string]string)}

			for i := 0; i < len(values); i++ {
				event.Data[values[i][1]] = values[i][2]
			}

			sendToWs(event)

			if conf.Debug {
				log.Println(event)
			}
		} else if strings.HasPrefix(s, "Press any key") {
			os.Exit(0)
		} else {
			event := basicEvent{eventName[1]}

			sendToWs(event)

			if conf.Debug {
				log.Println(event)
			}
		}

	}
	return len(p), nil
}

// Main Stuff

type config struct {
	StartupFile string
	Parameters  string
	WebPort     int
	Debug       bool
}

var trackerEventRegex = regexp.MustCompile("Tracker:([a-zA-Z]+)")
var trackerDataRegex = regexp.MustCompile("([a-zA-Z]+)\\(([^)]*)\\)")
var trk = tracker{}

var loading sync.WaitGroup
var connection *websocket.Conn
var jsScript string
var conf config

func main() {
	loading.Add(2)

	// Awesome print, because... cool
	fmt.Println(`  ______ _______ _______                   __               
 |   __ \__     |_     _|.----.---.-.----.|  |--.-----.----.
 |    __/     __| |   |  |   _|  _  |  __||    <|  -__|   _| by BigJk
 |___|  |_______| |___|  |__| |___._|____||__|__|_____|__|
`)

	checkConfig()

	// Checks if it exists
	if _, err := os.Stat(conf.StartupFile); os.IsNotExist(err) {
		fmt.Println(" -", conf.StartupFile, "not existing. Can't start Project Zomboid.")
		return
	}

	loadJsScript()

	go web()
	go waitForLoading()

	fmt.Println(" + Starting " + conf.StartupFile + "...")
	cmd := exec.Command(conf.StartupFile, strings.Split(conf.Parameters, " ")...)
	cmd.Stdout = trk
	loading.Done()
	cmd.Run()
}

func loadJsScript() {
	b, err := ioutil.ReadFile("inject-min.js")
	if err == nil {
		jsScript = string(b)
		jsScript = strings.Replace(jsScript, "%PORT%", fmt.Sprint(conf.WebPort), -1)
	}
}

func checkConfig() {
	if _, err := os.Stat("pztracker_config.json"); os.IsNotExist(err) {
		setupConfig()
	} else {
		b, _ := ioutil.ReadFile("pztracker_config.json")
		json.Unmarshal(b, &conf)
	}
}

func setupConfig() {
	bit64 := strings.Contains(runtime.GOARCH, "64")
	if bit64 {
		conf = config{"ProjectZomboid64.exe", "-Xmx1024m -Xms1024m -- -nosteam", 9090, false}
	} else {
		conf = config{"ProjectZomboid32.exe", "-Xmx768m -Xms768m -- -nosteam", 9090, false}
	}
	b, _ := json.MarshalIndent(conf, "", "	")
	ioutil.WriteFile("pztracker_config.json", b, 0777)
	fmt.Println(" + Missing config file created. (pztracker_config.json)")
}

func waitForLoading() {
	loading.Wait()
	if len(jsScript) > 0 {
		fmt.Println("\n > Insert into development console: $.getScript('http://127.0.0.1:" + fmt.Sprint(conf.WebPort) + "/js');")
	}
	fmt.Println("\n__/ PZTRACKER LOG \\__________________________________________________")
}

// Websocket Stuff

func sendToWs(data interface{}) {
	if connection == nil {
		return
	}
	b, _ := json.Marshal(data)
	connection.Write(b)
}

func wsHandler(ws *websocket.Conn) {
	log.Println("Websocket connected")
	if connection != nil {
		ws.Close()
	}
	connection = ws

	for {
		msg := make([]byte, 1)
		_, err := ws.Read(msg)
		if err != nil {
			break
		}
		time.Sleep(100)
	}

	log.Println("Websocket disconnected")
	connection = nil
}

func jsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, jsScript)
}

func web() {
	fmt.Println(" + Starting Webservice...")

	http.Handle("/websocket", websocket.Handler(wsHandler))
	http.HandleFunc("/js", jsHandler)

	loading.Done()
	err := http.ListenAndServe(":"+fmt.Sprint(conf.WebPort), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
