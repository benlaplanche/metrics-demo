package main

import (
	"fmt"
	dmetrics "github.com/cloudfoundry/dropsonde/metrics"
	redis "github.com/garyburd/redigo/redis"
	"strings"
)

type Client interface {
	Disconnect() error
	Info() (map[string]string, error)
	InfoField(fieldName string) (string, error)
	EmitMetric(m *metric) error
	Address() string
}

func (c *client) EmitMetric(m *metric) error {

	err := dmetrics.SendValue(m.Name, m.Value, m.Unit)

	if err != nil {
		return fmt.Errorf("Error emitting metric %v", m)
	}
	return nil
}

func (c *client) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c *client) Disconnect() error {
	return c.Connection.Close()
}

func (client *client) Info() (map[string]string, error) {
	info := map[string]string{}

	response, err := redis.String(client.Connection.Do("info"))
	if err != nil {
		return nil, err
	}

	for _, entry := range strings.Split(response, "\n") {
		trimmedEntry := strings.TrimSpace(entry)
		if trimmedEntry == "" || trimmedEntry[0] == '#' {
			continue
		}

		pair := strings.Split(trimmedEntry, ":")
		info[pair[0]] = pair[1]
	}

	return info, nil
}

func (client *client) InfoField(fieldName string) (string, error) {
	info, err := client.Info()
	if err != nil {
		return "", fmt.Errorf("Error during redis info: %s" + err.Error())
	}

	value, ok := info[fieldName]
	if !ok {
		return "", fmt.Errorf("Unknown field: %s", fieldName)
	}

	return value, nil
}
