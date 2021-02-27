package main

import (
	"github.com/dchest/uniuri"
	// "github.com/patrickmn/go-cache"
)


func createSession() string {
	return randomString(4)
}

func getSession(id string) {

}

func randomString(len int) string {
	return uniuri.NewLenChars(len, []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ"))
}