package consul

import (
	"encoding/base64"
	"github.com/hashicorp/consul/api"

	"fmt"
)

func UpdateKv(client *api.Client, key, content string) error {
	valueBytes, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		fmt.Printf("Error decoding base64 value for key %s: %v\n", key, err)
		return err
	}
	_, err = client.KV().Put(&api.KVPair{Key: key, Value: valueBytes}, nil)
	return err
}

func DeleteKv(client *api.Client, key string) error {
	_, err := client.KV().Delete(key, nil)
	return err
}

func GetConsulKeys(client *api.Client, prefix, separator string) ([]string, error) {
	keys, _, err := client.KV().Keys(prefix, separator, nil)
	if err != nil {
		return nil, err
	}
	return keys, nil
}

func GetConsulValue(client *api.Client, key string) (string, error) {
	kvPair, _, err := client.KV().Get(key, nil)
	if err != nil {
		return "", err
	}
	return string(kvPair.Value), nil
}
