package main

import (
	"fmt"
	"github.com/cloudfoundry/dropsonde"
	dmetrics "github.com/cloudfoundry/dropsonde/metrics"
	redis "github.com/garyburd/redigo/redis"
	"strconv"
	"strings"
)

type client struct {
	Host string
	Port int

	Connection redis.Conn
}

type metric struct {
	Name  string
	Value float64
	Unit  string
}

type emitterDetails struct {
	Destination string
	Origin      string
	Zone        string
	Index       string
}

type Client interface {
	Disconnect() error
	Info() (map[string]string, error)
	InfoField(fieldName string) (string, error)
	EmitMetric(m *metric) error
	Address() string
}

func Connect() (Client, error) {
	client := &client{
		Host: "127.0.0.1",
		Port: 6379,
	}

	address := fmt.Sprintf("%s:%d", client.Host, client.Port)

	var err error
	client.Connection, err = redis.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	return client, nil

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

func main() {
	redisConn, err := Connect()

	if err != nil {
		fmt.Println(err)
		return
	}

	defer redisConn.Disconnect()

	d := &emitterDetails{
		Destination: "localhost:3457",
		Origin:      "metrics-demo",
		Zone:        "z1",
		Index:       "0",
	}

	err = dropsonde.Initialize(d.Destination, d.Origin, d.Zone, d.Index)
	if err != nil {
		fmt.Println(err)
		return
	}

	fieldName := "uptime_in_seconds"
	uptime, err := redisConn.InfoField(fieldName)
	if err != nil {
		fmt.Println(err)
		return
	}

	uptimeConverted, err := strconv.ParseFloat(uptime, 64)
	if err != nil {
		fmt.Println(err)
	}

	metricData := &metric{
		Name:  fieldName,
		Value: uptimeConverted,
		Unit:  "",
	}

	err = redisConn.EmitMetric(metricData)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Successfully emitted %+v", metricData)
}
