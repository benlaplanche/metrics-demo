package main

import (
	redis "github.com/garyburd/redigo/redis"
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
