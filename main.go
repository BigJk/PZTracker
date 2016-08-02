package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
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

			if debug {
				fmt.Println(event)
			}
		} else {
			event := basicEvent{eventName[1]}

			sendToWs(event)

			if debug {
				fmt.Println(event)
			}
		}

	}
	return len(p), nil
}

// Main Stuff

var debug = false
var trackerEventRegex = regexp.MustCompile("Tracker:([a-zA-Z]+)")
var trackerDataRegex = regexp.MustCompile("([a-zA-Z]+)\\(([^)]*)\\)")
var trk = tracker{}

var connection *websocket.Conn

func main() {
	// Awesome print, because... cool
	fmt.Println(`  ______ _______ _______                   __               
 |   __ \__     |_     _|.----.---.-.----.|  |--.-----.----.
 |    __/     __| |   |  |   _|  _  |  __||    <|  -__|   _| by BigJk
 |___|  |_______| |___|  |__| |___._|____||__|__|_____|__|
`)

	// Check for 64 bit and decide wich bat file to start.
	bit64 := strings.Contains(runtime.GOARCH, "64")
	path := "ProjectZomboid32.bat"
	if bit64 {
		path = "ProjectZomboid64.bat"
	}

	// Check if a custom starting file is given.
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	// Checks if it exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println(" -", path, "not existing. Can't start Project Zomboid.")
		return
	}

	go web()

	fmt.Println(" + Starting " + path + "...")
	cmd := exec.Command(path)
	cmd.Stdout = trk
	cmd.Run()
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
	fmt.Println("Websocket connected...")
	if ws != nil {
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

	fmt.Println("Websocket disconnected...")
	connection = nil
}

func web() {
	fmt.Println(" + Starting Webservice...")
	http.Handle("/ws", websocket.Handler(wsHandler))
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
