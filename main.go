package main

import (
	"context"
	"fmt"
	"regexp"

	"github.com/sensu-community/sensu-plugin-sdk/httpclient"
	"github.com/sensu-community/sensu-plugin-sdk/sensu"
	corev2 "github.com/sensu/sensu-go/api/core/v2"
	"github.com/sensu/sensu-go/types"
)

// Config represents the handler plugin config.
type Config struct {
	sensu.PluginConfig
	APIKey        string
	APIURL        string
	MaxOccurrence int
}

var (
	plugin = Config{
		PluginConfig: sensu.PluginConfig{
			Name:     "sensu-event-keepalive-ttl-handler",
			Short:    "The Sensu Go Remove events as keepalive",
			Keyspace: "sensu.io/plugins/sensu-event-keepalive-ttl-handler/config",
		},
	}

	options = []*sensu.PluginConfigOption{
		{
			Path:      "apikey",
			Env:       "API_JEY",
			Argument:  "apikey",
			Shorthand: "k",
			Default:   "",
			Usage:     "The Sensu Api Key , use default from API_JEY env var",
			Value:     &plugin.APIKey,
		},
		{
			Path:      "url",
			Env:       "SENSU_URL",
			Argument:  "sensu_url",
			Shorthand: "u",
			Default:   "",
			Usage:     "The Sensu Api URL  , use default from SENSU_URL env var",
			Value:     &plugin.APIURL,
		},
		{
			Path:      "max_occurrence",
			Env:       "SENSU_MAX_OCCURRENCE",
			Argument:  "max_occurrence",
			Shorthand: "m",
			Default:   0,
			Usage:     "The Max event Occurence after ttl to remove event   , use default from SENSU_MAX_OCCURRENCE env var",
			Value:     &plugin.MaxOccurrence,
		},
	}
)

func main() {
	handler := sensu.NewGoHandler(&plugin.PluginConfig, options, checkArgs, executeHandler)
	handler.Execute()
}

func checkArgs(_ *types.Event) error {
	if len(plugin.APIKey) == 0 {
		return fmt.Errorf("apikey sensu  is empty")
	}
	if len(plugin.APIURL) == 0 {
		return fmt.Errorf("url sensu  is empty")
	}
	return nil
	if plugin.MaxOccurrence == 0 {
		return fmt.Errorf("max occurrence sensu  is 0")
	}
	return nil
}
func executeHandler(event *types.Event) error {

	res, err := CheckAnnotations(event)
	if err != nil {
		fmt.Println("Check Annotations error ", err)
	}
	if res {
		for {
			fmt.Println("Event as correct Annotations")
			fmt.Println("Check Name :", event.Check.Name)
			fmt.Println("Namespace: ", event.Entity.Namespace)
			fmt.Println("Entity: ", event.Entity.Name)
			dl := ClientDeleteResourceEvent(event)
			if dl {
				fmt.Println("Removed Sucess")

			} else {
				fmt.Println("Fail to Remove Event by ttl")

			}
			if CheckEventAlreadexist(event) {
				break
			}
		}
		// return nil

	}

	return nil
}

// CheckEventAlreadexist  teste
func CheckEventAlreadexist(event *types.Event) bool {

	config := httpclient.CoreClientConfig{
		URL:    plugin.APIURL,
		APIKey: plugin.APIKey,
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

// ClientDeleteResourceEvent  teste
func ClientDeleteResourceEvent(event *types.Event) bool {
	config := httpclient.CoreClientConfig{
		URL:    plugin.APIURL,
		APIKey: plugin.APIKey,
	}
	if CheckEventParameter(event) {

		cl := httpclient.NewCoreClient(config)
		req := httpclient.NewEventRequest(event.Check.Namespace, event.Entity.Name, event.Check.Name)
		resp, err := cl.DeleteResource(context.Background(), req)
		if err != nil {
			fmt.Println("Remove Event Check is not possible ")
			return false
		}
		switch {

		case resp.StatusCode == 204:
			fmt.Printf("HTTP Response %v\n", resp.Status)
			fmt.Printf("Event Check %v Removed Success from Entity %v \n", event.Check.Name, event.Entity.Name)
			return true

		case resp.StatusCode == 404:

			fmt.Printf("Event Check %v has already been removed from Entity %v \n", event.Check.Name, event.Entity.Name)
			return false
		}
	} else {
		fmt.Println("\nEvent not inclued on ttl parameters ")
		return false
	}
	return false
}

// CheckAnnotations teste
func CheckAnnotations(event *types.Event) (bool, error) {

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

//CheckEventParameter Check event
func CheckEventParameter(event *types.Event) bool {
	var MaxOccurrence int64
	MaxOccurrence = int64(plugin.MaxOccurrence)

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
