package main

import (
	"log"
	"time"
)

type JSONTime time.Time

func (jt *JSONTime) UnmarshalJSON(b []byte) error {
	t, err := time.Parse("2006-01-02 15:04:05", string(b))
	if err != nil {
		log.Fatal(err)
	}

	*jt = JSONTime(t)

	return err
}

type Parameter struct {
	Title           string  `json:"parametername"`
	ChemicalFormula string  `json:"chemicalFormula"`
	Name            string  `json:"name"`
	Class           string  `json:"class"`
	Concentration   float32 `json:"modifyav"`
	Norma1          float32 `json:"norma"`
	Norma2          float32 `json:"norma_2"`
	TimePoint       string  `json:"dateTime"`
}

type AirData struct {
	StationName string      `json:"stationName"`
	Parameters  []Parameter `json:"parameters"`
}
