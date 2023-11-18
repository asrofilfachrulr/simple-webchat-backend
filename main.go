package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func checkOrigin(r *http.Request) bool {
	return true
}

func main() {
	http.HandleFunc("/ws", wsHanlder)
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})

	fmt.Println("Start server at :8080")
	log.Fatalln(http.ListenAndServe(":8080", nil))
}

func wsHanlder(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = checkOrigin
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err.Error())
		return
	}

	defer conn.Close()

	for {
		mt, p, err := conn.ReadMessage()

		if err != nil {
			log.Println(err.Error())
			return
		}

		fmt.Printf("Received Message:\n\ttype:\t%v\n\tcontent:\t%s\n\n", mt, p)
	}
}
