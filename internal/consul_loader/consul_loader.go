package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	consulapi "github.com/hashicorp/consul/api"
)

func main() {
	var (
		m      = map[string]interface{}{}
		data   []byte
		err    error
		prefix = "config/"
	)

	if len(os.Args) < 2 {
		log.Fatal("Error. Give me path to .json file")
		return
	}

	path := os.Args[1]

	if data, err = ioutil.ReadFile(path); err != nil {
		log.Fatal("Read file error:", err.Error())
	}

	if err = json.Unmarshal([]byte(data), &m); err != nil {
		log.Fatal("Unmarhal error:", err.Error())
	}

	consul, err := consulapi.NewClient(consulapi.DefaultConfig())
	if err != nil {
		log.Fatal("Consul new client error:", err.Error())
	}

	sessionID, _, err := consul.Session().Create(nil, nil)
	if err != nil {
		log.Fatal("Consul session error:", err.Error())
	}

	err = addDataToConsulKV(m, prefix, consul, sessionID)
	if err != nil {
		log.Fatal("Some shit happens:", err.Error())
	} else {
		log.Println("Success")
	}
}

func putToConsulKV(key string, value interface{},
	consul *consulapi.Client, sessionID string) error {
	bytes := []byte(fmt.Sprintf("%v", value))

	_, _, err := consul.KV().Acquire(
		&consulapi.KVPair{
			Key:     key,
			Value:   bytes,
			Session: sessionID},
		nil)

	return err
}

func addDataToConsulKV(data map[string]interface{}, prefix string,
	consul *consulapi.Client, sessionID string) error {
	var err error
	for key, value := range data {
		if deeper, yes := value.(map[string]interface{}); yes {
			err = addDataToConsulKV(deeper, prefix+key+"/", consul, sessionID)
		} else {
			err = putToConsulKV(prefix+key, value, consul, sessionID)
		}
		if err != nil {
			break
		}
	}
	return err
}
