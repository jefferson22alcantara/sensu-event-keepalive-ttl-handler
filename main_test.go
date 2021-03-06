package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	corev2 "github.com/sensu/sensu-go/api/core/v2"
	"github.com/sensu/sensu-go/types"
)

var server = httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		event := corev2.FixtureEvent("localhost", "keepalive")
		_ = json.NewEncoder(w).Encode(event)
	}
}))

// Testannotations teste
func Testannotations(t *testing.T) (bool, error) {
	anotation := make(map[string]string)
	anotation["sensu.io/plugins/sensu-event-keepalive-ttl-handler/conf/enable"] = "true"
	event := &types.Event{
		ObjectMeta: corev2.ObjectMeta{Annotations: anotation},
	}
	if event.Check.Annotations != nil {
		for key, value := range event.Check.Annotations {
			if key == "sensu.io/plugins/sensu-event-keepalive-ttl-handler/conf/enable" {
				if value == "true" {
					return true, nil
				}

			}

		}

	}
	return false, nil
}

//Testparameter test
func Testparameter(t *testing.T) bool {
	event := &types.Event{
		Timestamp:            0,
		Entity:               &types.Entity{ObjectMeta: types.ObjectMeta{Name: "EntityName"}},
		Check:                &types.Check{ObjectMeta: types.ObjectMeta{Name: "CheckName"}, Output: "Last check execution was 30 seconds ago"},
		Metrics:              &corev2.Metrics{},
		ObjectMeta:           corev2.ObjectMeta{},
		ID:                   []byte{},
		Sequence:             0,
		XXX_NoUnkeyedLiteral: struct{}{},
		XXX_unrecognized:     []byte{},
		XXX_sizecache:        0,
	}
	var MaxOccurrence int64
	MaxOccurrence = 60

	re, err := regexp.Compile(`\w+\s\w+\s\w+\s\w+\s\d+\s\w+\s\w+`)
	if err != nil {
		fmt.Println("Regex Compile Error ")
		return false
	}

	if re.MatchString(event.Check.Output) {
		fmt.Println("Event Check Output [PASS]", event.Check.Output)
	} else {
		fmt.Println("Event Check Output [FAIL] Event not Match with Regex ttl [nothing to do]", event.Check.Output)
		return false
	}

	if event.Check.Status == 1 {
		fmt.Println("Event Check Status [PASS]", event.Check.Status)
	} else {
		fmt.Println("Event Check Status [FAIL], Event not Match with Check Status 1 [nothing to do]", event.Check.Status)
		return false
	}
	if event.Check.Occurrences > MaxOccurrence {
		fmt.Println("Event Check Ocurrences [PASS]", event.Check.Occurrences)
	} else {
		fmt.Printf("Event Check Ocurrences [FAIL] event.Occurrences:%v is less Than MaxOccurrence Configured:%v [nothing to do]", event.Check.Occurrences, MaxOccurrence)
		return false
	}

	return true
}
