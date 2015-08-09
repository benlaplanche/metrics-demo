package main

import (
	"errors"
	"fmt"
	redis "github.com/garyburd/redigo/redis"
	"strings"
)

type client struct {
	host string
	port int

	connection redis.Conn
}

type metric struct {
	name  string
	value string
	unit  string
}

type metronAgent struct {
}

type Client interface {
	Disconnect() error
	Info() (map[string]string, error)
	InfoField(fieldName string) (string, error)
	EmitMetric(m metric, agent metronAgent) error
	Address() string
}

func Connect() (Client, error) {
	client := &client{
		host: "127.0.0.1",
		port: 6379,
	}

	address := fmt.Sprintf("%s:%d", client.host, client.port)

	var err error
	client.connection, err = redis.Dial("tcp", address)
	if err != nil {
		fmt.Errorf("Error connecting to the Redis Server at: %s", address)
		return nil, err
	}

	return client, nil

}

func (c *client) EmitMetric(m metric, agent metronAgent) error {
	return nil
}

func (c *client) Address() string {
	return fmt.Sprintf("%s:%d", c.host, c.port)
}

func (c *client) Disconnect() error {
	return c.connection.Close()
}

func (client *client) Info() (map[string]string, error) {
	info := map[string]string{}

	response, err := redis.String(client.connection.Do("info"))
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
		return "", errors.New(fmt.Sprintf("Unknown field: %s", fieldName))
	}

	return value, nil
}

func main() {
	redisConn, err := Connect()
	defer redisConn.Disconnect()

	if err != nil {
		fmt.Errorf("Error connecting to Redis at %s", redisConn.Address())
		return
	}

	field := "uptime_in_seconds"
	uptime, err := redisConn.InfoField(field)
	if err != nil {
		fmt.Errorf("Error reading field: %s", field)
		return
	}

	fmt.Println(uptime)
}
