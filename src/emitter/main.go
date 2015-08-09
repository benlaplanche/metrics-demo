package main

import (
	"fmt"
	"github.com/cloudfoundry/dropsonde"
	redis "github.com/garyburd/redigo/redis"
	"strconv"
)

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
