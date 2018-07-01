package main

import (
	"os"
	"net/http"
	"strings"
	"encoding/json"
	"bytes"
	"time"
	"fmt"
	"io/ioutil"
)

func main() {
	token := readConfig("TOKEN")
	organization := readConfig("ORGANIZATION")
	volumes := strings.Split(readConfig("VOLUMES"), ",")
	timestamp := time.Now().Unix()
	httpClient := &http.Client{}
	for _, volume := range volumes {
		requestJson, err := json.Marshal(map[string]string{"name": fmt.Sprintf("%d_%s", timestamp, volume),
			"organization": organization,
			"volume_id": volume,
		})
		if err != nil {
			panic(err)
		}
		request, err := http.NewRequest("POST", "https://cp-par1.scaleway.com/snapshots", bytes.NewReader(requestJson))
		if err != nil {
			panic(err)
		}
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("X-Auth-Token", token)
		response, err := httpClient.Do(request)
		if err != nil {
			panic(err)
		}

		if !strings.HasPrefix(response.Status, "20") {
			defer response.Body.Close()
			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				panic(err)
			}
			println(string(body))
		} else {
			fmt.Printf("Backup created for volume %s.\n", volume)
		}
	}
}

func readConfig(key string) string {
	value, found := os.LookupEnv(key)
	if !found {
		panic("Configuration not found: " + key)
	}
	return value
}
