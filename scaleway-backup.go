package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func main() {
	token := readConfig("TOKEN")
	organization := readConfig("ORGANIZATION")
	volumes := strings.Split(readConfig("VOLUMES"), ",")

	backupsNumberConfig := readConfigOrDefault("BACKUPS_NUMBER", "2")
	backupsNumber, err := strconv.Atoi(backupsNumberConfig)
	if err != nil {
		panic(err)
	}

	httpClient := &http.Client{}

	makeBackup(httpClient, token, organization, volumes)
	purgeOldBackups(httpClient, token, organization, volumes, backupsNumber)
}

func makeBackup(httpClient *http.Client, token string, organization string, volumes []string) {
	timestamp := time.Now().Unix()
	for _, volume := range volumes {
		requestJson, err := json.Marshal(map[string]string{"name": fmt.Sprintf("%d_%s", timestamp, volume),
			"organization": organization,
			"volume_id":    volume,
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

func purgeOldBackups(httpClient *http.Client, token string, organization string, volumes []string, backupsNumber int) {
	request, err := http.NewRequest("GET", "https://cp-par1.scaleway.com/snapshots", nil)
	if err != nil {
		panic(err)
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-Auth-Token", token)
	response, err := httpClient.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	snapshotsFromResponse := snapshots{}
	err = json.Unmarshal(body, &snapshotsFromResponse)
	if err != nil {
		panic(err)
	}

	for _, volume := range volumes {
		var volumeBackups []snapshot
		for _, snapshot := range snapshotsFromResponse.Snapshots {
			if snapshot.Organization == organization && snapshot.BaseVolume.Id == volume {
				volumeBackups = append(volumeBackups, snapshot)
			}
		}
		sort.Slice(volumeBackups, func(i, j int) bool {
			return volumeBackups[i].CreationDate > volumeBackups[j].CreationDate
		})
		for i, snapshot := range volumeBackups {
			if i >= backupsNumber {
				request, err := http.NewRequest("DELETE", "https://cp-par1.scaleway.com/snapshots/"+snapshot.Id, nil)
				if err != nil {
					panic(err)
				}
				request.Header.Add("X-Auth-Token", token)
				fmt.Printf("About to delete old (%s) snapshot of volume %s with id %s.\n", snapshot.CreationDate, snapshot.BaseVolume.Id, snapshot.Id)
				_, err = httpClient.Do(request)
				if err != nil {
					panic(err)
				}
			}
		}
	}
}

// Configuration helpers

func readConfig(key string) string {
	value, found := os.LookupEnv(key)
	if !found {
		panic("Configuration not found: " + key)
	}
	return value
}

func readConfigOrDefault(key string, defaultValue string) string {
	value, found := os.LookupEnv(key)
	if found {
		return value
	}
	return defaultValue
}

// Json types

type snapshots struct {
	Snapshots []snapshot `json:"snapshots"`
}

type snapshot struct {
	BaseVolume   baseVolume `json:"base_volume"`
	Id           string     `json:"id"`
	CreationDate string     `json:"creation_date"`
	Organization string     `json:"organization"`
}

type baseVolume struct {
	Id string `json:"id"`
}
