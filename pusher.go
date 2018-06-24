package main

import (
	"log"

	"github.com/influxdata/influxdb/client/v2"
)

type Pusher struct {
	c      client.Client
	config PusherConfig
}

type PusherConfig struct {
	DB       string
	Addr     string
	Username string
	Password string
}

func NewPusher(config PusherConfig) (pusher Pusher, err error) {
	pusher.config = config

	pusher.c, err = client.NewHTTPClient(client.HTTPConfig{
		Addr:     config.Addr,
		Username: config.Username,
		Password: config.Password,
	})

	if err != nil {
		return Pusher{}, err
	}

	return
}

func (p *Pusher) Close() error {
	return p.Close()
}

func (p *Pusher) Push(points []*client.Point) error {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  p.config.DB,
		Precision: "m",
	})
	if err != nil {
		log.Fatal(err)
	}

	bp.AddPoints(points)

	if err := p.c.Write(bp); err != nil {
		return err
	}

	if err := p.c.Close(); err != nil {
		return err
	}

	return err
}
