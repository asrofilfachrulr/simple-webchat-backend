package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var activeConnections = make(map[*websocket.Conn]bool)

func checkOrigin(r *http.Request) bool {
	return true
}

func main() {
	runtime.GOMAXPROCS(2)

	r := mux.NewRouter()

	r.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})
	r.HandleFunc("/ws", wsHanlder)

	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("Listening on :8080...")
	log.Fatalln(srv.ListenAndServe())
}

func wsHanlder(w http.ResponseWriter, r *http.Request) {
	// upgrading the protocol
	upgrader.CheckOrigin = checkOrigin
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err.Error())
		return
	}

	activeConnections[conn] = true

	fmt.Printf("Updated Connections: %v\n", activeConnections)

	defer func() {
		conn.Close()
		delete(activeConnections, conn)
		fmt.Printf("Updated Connections: %v\n", activeConnections)
	}()

	// main loop to continously do socket listening
	for {
		mt, p, err := conn.ReadMessage()

		if err != nil {
			log.Println(err.Error())
			return
		}

		// incoming data model
		data := struct {
			Name    string `json:"name"`
			Content string `json:"content"`
		}{}

		// parse incoming message to predefined struct
		if err := json.Unmarshal(p, &data); err != nil {
			log.Fatalln(err.Error())
			return
		}

		log.Printf("Received Message:\n\ttype:\t%v\n\tcontent:\t%v\n\n", mt, data)

		// Parse data to json []byte then send to active sockets/connections
		if message, err := json.Marshal(data); err != nil {
			log.Fatalln(err.Error())
		} else {
			// broadcast message concurently
			go broadcastMessage(message)
		}
	}
}

func broadcastMessage(message []byte) {
	var wg sync.WaitGroup
	logChan := make(chan string)
	successAttemptCounter := 0

	// logging successful sent message attempt, receive log data from channel logChan
	go logSentMessages(logChan, &successAttemptCounter)

	for conn := range activeConnections {
		// using waitgroup so upon broadcasted message to all active connections, is logged correctly
		wg.Add(1)

		// attempt to send message to a connection concurently
		go attemptWriteMessage(message, conn, &wg, logChan)
	}

	// wait all attempt done
	wg.Wait()

	// close logging channel
	close(logChan)

	// log all attempts are done
	log.Printf("Done sending attempt to all active connections\nSent\t: %d\nFail\t: %d", successAttemptCounter, len(activeConnections)-successAttemptCounter)
}

func attemptWriteMessage(message []byte, conn *websocket.Conn, wg *sync.WaitGroup, logChan chan<- string) {
	if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
		log.Println(err.Error())
	} else {
		logChan <- fmt.Sprintf("Message sent to connection: %p", conn)
	}

	// done either success or fail
	wg.Done()
}

func logSentMessages(logChan <-chan string, counter *int) {
	for logCh := range logChan {
		*counter = *counter + 1
		log.Println(logCh)
	}
}
