package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

//  Creates type Person, with a name and an age
type Person struct {
	Name string
	Age  int
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {

	// Opens and reads file index html file
	indexFile, _ := os.Open("html/index.html")
	index, _ := ioutil.ReadAll(indexFile)

	//  HTTP handler for upgrading to websocket
	http.HandleFunc("/websocket", func(res http.ResponseWriter, req *http.Request) {
		conn, err := upgrader.Upgrade(res, req, nil)

		// basic error checking
		if err != nil {
			fmt.Println(err)
			return
		}

		//  Append to client's text area, that they've connected
		fmt.Println("Client Connected!")

		//  Hardcoded variable of type Person for test purposes
		myPerson := Person{
			Name: "Duffy",
			Age:  0,
		}

		//  infinite loop to handle websocket - ends at channel closure, or age-out @ 40 seconds.
		for {
			time.Sleep(2 * time.Second) //wait two seconds
			if myPerson.Age < 40 {      // if not age'd out
				JsonToView, err := json.Marshal(myPerson) //marshal variable into json for sending err
				if err != nil {
					fmt.Println(err)
					return
				}
				err = conn.WriteMessage(websocket.TextMessage, JsonToView) //sends JSON that contains person variable to view
				if err != nil {
					fmt.Println(err)
					break
				}
				myPerson.Age += 2 //increment age by two (because we waited two second)
			} else {
				conn.Close() // if age'd out, close connection.
				break
			}
		}
		fmt.Println("Client unsubscribed")
	})

	//  serves index
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(res, string(index))
	})

	// Start server on port 3000
	http.ListenAndServe(":3000", nil)
}
