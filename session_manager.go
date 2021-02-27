package main
import (
	"fmt"
	"github.com/dchest/uniuri"
)


type session struct {
	id             string
	lastUpdateTime time.Time
	conn           net.Conn
	clients        cmap.ConcurrentMap
}

func createSession() session {
	return randomString();
}

func getSession(id string) {

}

func main() { 
	fmt.Println(randomString((4)))
}

func randomString(len int) string {
	return uniuri.NewLenChars(len, []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ"))
}