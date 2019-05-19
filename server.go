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
}

// TODO Some sort of reaper process that cleans up old sessions
var sessionMap = cmap.New()

// TODO Some better id generation
var id = 1

func createSession(w http.ResponseWriter, r *http.Request) {
	sessionID := fmt.Sprintf("%s%d", "id-", id)
	newSession := session{id: sessionID, lastUpdateTime: time.Now()}
	fmt.Println(newSession)
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
	} else {
		// Some error handling?
	}
}

func connectToSession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["sessionID"]
	if tmp, ok := sessionMap.Get(sessionID); ok {
		newSession := tmp.(*session)
		fmt.Println(newSession)
		readConn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			// handle error
		}
		writeConn := newSession.conn
		go func() {
			defer readConn.Close()

			for {
				msg, op, err := wsutil.ReadClientData(readConn)
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Println("Writing : "+string(msg), len(msg))
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
	// TODO Just make a single file?
	r.Handle("/", http.FileServer(http.Dir("./static/home")))
	r.PathPrefix("/controller/").Handler(http.FileServer(http.Dir("./static")))
	r.PathPrefix("/library/").Handler(http.FileServer(http.Dir("./static")))
	r.PathPrefix("/examples/").Handler(http.FileServer(http.Dir("./static")))
	// TODO Fix the semantics of this at some stage(make post?)
	r.HandleFunc("/session/", createSession)
	r.HandleFunc("/session/{sessionID}/start", startSession)
	r.HandleFunc("/session/{sessionID}/connect", connectToSession)
	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "80"
	}
	if err := http.ListenAndServe(":"+port, r); err != nil {
		panic(err)
	}
}
