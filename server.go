package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	cmap "github.com/orcaman/concurrent-map" // Also has some simple tests
)

type session struct {
	id             string
	lastUpdateTime time.Time
	conn           net.Conn
	clients        cmap.ConcurrentMap
}

// TODO Some sort of reaper process that cleans up old sessions
var sessionMap = cmap.New()

// TODO Some better id generation
var id = 1

func createSession(w http.ResponseWriter, r *http.Request) {
	sessionID := fmt.Sprintf("%s%d", "id-", id)
	newSession := session{
		id:             sessionID,
		lastUpdateTime: time.Now(),
		clients:        cmap.New()}
	sessionMap.Set(sessionID, &newSession)
	id = id + 1
	w.Write([]byte(sessionID))
}

func startSession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["sessionID"]
	if tmp, ok := sessionMap.Get(sessionID); ok {
		newSession := tmp.(*session)
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			// handle error
		}
		newSession.conn = conn
		newSession.lastUpdateTime = time.Now()
		fmt.Println(newSession)
		// Write to clients
		/*
		 * FIXME This whole goroutine is a bit flaky.
		 * It will sometimes error out when reading the client data.
		 * with an EOF
		 */
		go func() {
			// defer conn.Close()

			for {
				msg, op, err := wsutil.ReadClientData(conn)
				if err != nil {
					fmt.Println("Error reading connection from host")
					fmt.Println(err)
					return
				}
				// For now the msg is  clientID
				clientID := string(msg[:])
				if tmp2, ok2 := newSession.clients.Get(clientID); ok2 {
					clientConn := tmp2.(net.Conn)
					err = wsutil.WriteServerMessage(clientConn, op, msg)
					if err != nil {
						fmt.Println(err)
						return
					}
				} else {
					fmt.Println("Error: Could not find client with id: " + clientID)
					// Error handling. Could no find client id
				}
			}
		}()
	} else {
		// Some error handling?
	}
}

func connectToSession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["sessionID"]
	clientID := vars["clientID"]
	if tmp, ok := sessionMap.Get(sessionID); ok {
		newSession := tmp.(*session)
		readConn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			// handle error
			fmt.Println("Error connecting to session")
			fmt.Println(err)
		}
		writeConn := newSession.conn
		newSession.clients.Set(clientID, readConn)
		fmt.Println("Client connected: " + clientID)
		go func() {
			defer readConn.Close()

			for {
				msg, op, err := wsutil.ReadClientData(readConn)
				if err != nil {
					fmt.Println(err)
					return
				}
				err = wsutil.WriteServerMessage(writeConn, op, msg)
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}()
	} else {
		// Some error handling?
	}
}

func main() {
	r := mux.NewRouter()
	r.PathPrefix("/controller/").Handler(http.FileServer(http.Dir("./static")))
	r.PathPrefix("/library/").Handler(http.FileServer(http.Dir("./static")))
	r.PathPrefix("/examples/").Handler(http.FileServer(http.Dir("./static")))
	// TODO Fix the semantics of this at some stage(make post?)
	r.HandleFunc("/session/", createSession)
	r.HandleFunc("/session/{sessionID}/start", startSession)
	r.HandleFunc("/session/{sessionID}/connect", connectToSession).Queries("clientID", "{clientID}")
	// default
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/home")))
	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "80"
	}
	if err := http.ListenAndServe(":"+port, r); err != nil {
		panic(err)
	}
}
