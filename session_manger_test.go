package main

import "testing"

func TestCreateSession(t *testing.T) {
	sessionID := createSession()
   if ( len(sessionID)  != 4) {
	   t.Error("Expect the session id to be 4 characters long")
   }
}