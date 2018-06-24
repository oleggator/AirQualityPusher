package main

import (
	"encoding/json"
	"fmt"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/jasonlvhit/gocron"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

func main() {
	addr := os.Getenv("INFLUXDB_ADDRESS")
	if addr == "" {
		log.Fatal("Variable 'INFLUXDB_ADDRESS' is empty")
	}

	db := os.Getenv("DB_NAME")
	if addr == "" {
		log.Fatal("Variable 'DB_NAME' is empty")
	}

	username := os.Getenv("DB_LOGIN")
	password := os.Getenv("DB_PASSWORD")

	pusher, err := NewPusher(PusherConfig{
		DB:       db,
		Addr:     addr,
		Username: username,
		Password: password,
	})
	defer pusher.Close()
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 24; i++ {
		timeString := fmt.Sprintf("%d:11", i)
		gocron.Every(1).Day().At(timeString).Do(refresh, &pusher)
	}

	log.Println("Starting...")
	err = refresh(&pusher)
	if err != nil {
		log.Fatal(err)
	}

	<-gocron.Start()
}

func refresh(pusher *Pusher) (err error) {
	log.Println("Refreshing...")

	airData, err := getAirData()
	if err != nil {
		return err
	}

	points := make([]*client.Point, len(airData.Parameters))
	for i, parameter := range airData.Parameters {
		loc, err := time.LoadLocation("Europe/Moscow")
		if err != nil {
			return err
		}

		t, err := time.ParseInLocation("2006-01-02 15:04:05", parameter.TimePoint, loc)
		if err != nil {
			return err
		}

		tags := map[string]string{
			"name":             parameter.Name,
			"chemical_formula": parameter.ChemicalFormula,
			"class":            parameter.Class,
		}

		fields := map[string]interface{}{
			"concentration": parameter.Concentration,
		}

		points[i], err = client.NewPoint(airData.StationName, tags, fields, t)
		if err != nil {
			return err
		}
	}

	return pusher.Push(points)
}

func getAirData() (*AirData, error) {
	stationName := os.Getenv("STATION_NAME")
	if stationName == "" {
		stationName = "Гурьянова"
	}

	v := url.Values{}
	v.Set("locale", "ru_RU")
	v.Add("station_name", stationName)
	v.Add("mapType", "air")

	res, err := http.PostForm("http://178.208.145.33/wp-content/themes/moseco/map/station-popup.php", v)
	if err != nil {
		return nil, err
	}

	blob, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}

	var airData AirData
	err = json.Unmarshal(blob, &airData)
	if err != nil {
		return nil, err
	}

	return &airData, nil
}
