package main

import (
	"testing"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestSplitNatsSubject(t *testing.T) {
	subj := "thingsMsg.1a202f88-9d55-11ec-b909-0242ac120002"
	want := "1a202f88-9d55-11ec-b909-0242ac120002"
	thingId, err := getThingIdFromNatsSubjet(subj)

	if thingId != want || err != nil {
		t.Fatal("want", want, "got", thingId)
	}
}
