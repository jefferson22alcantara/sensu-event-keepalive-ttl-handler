package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/sensu-community/sensu-plugin-sdk/httpclient"
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
	plugin.MAXOCCURENCE = 60

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
	if event.Check.Occurrences > plugin.MAXOCCURENCE {
		fmt.Println("Event Check Ocurrences [PASS]", event.Check.Occurrences)
	} else {
		fmt.Printf("Event Check Ocurrences [FAIL] event.Occurrences:%v is less Than MaxOccurrence Configured:%v [nothing to do]", event.Check.Occurrences, plugin.MAXOCCURENCE)
		return false
	}

	return true
}

// Testeventalreadexist  teste
func Testeventalreadexist(event *types.Event) bool {

	config := httpclient.CoreClientConfig{
		URL:    server.URL,
		APIKey: "use transport layer security",
		CACert: server.Certificate(),
	}
	cl := httpclient.NewCoreClient(config)
	req := httpclient.NewEventRequest(event.Check.Namespace, event.Entity.Name, event.Check.Name)
	ev := new(corev2.Event)
	resp, err := cl.GetResource(context.Background(), req, ev)
	if err != nil {
		fmt.Println("Check Event not Found ", err)
	}
	var result bool
	switch {
	case resp.StatusCode == 404:
		result = true
	case resp.StatusCode == 200:
		result = false
	}
	return result
}
