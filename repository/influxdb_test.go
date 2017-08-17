package repository

import (
	"log"
	"testing"

	client "github.com/influxdata/influxdb/client/v2"
)

var (
	idb *InfluxDB
)

func init() {
	idb, _ = NewInfluxDBClient("http://localhost:8086")
}

func TestNewInfluxDBClient(t *testing.T) {
	influxdb, err := NewInfluxDBClient("http://localhost:8086")
	if err != nil {
		t.Errorf("Got error %s", err)
	}

	if influxdb == nil {
		t.Error("InfluxDB shoud not be nil")
	}

	_, err = NewInfluxDBClient("thisisnotanurl")
	if err == nil {
		t.Error("Got no error")
	}
}

func TestSaveInfluxDB(t *testing.T) {
	defer teardown()

	data := map[string]interface{}{
		"stuff":     10.00,
		"something": 20.00,
		"total":     30.00,
	}

	err := idb.Save("test_data", data)
	if err != nil {
		t.Errorf("Save Got error %s", err)
	}

}

func teardown() {
	q := client.NewQuery("DROP SERIES FROM test_data", "tests", "s")
	_, err := idb.client.Query(q)
	if err != nil {
		log.Fatal(err)
	}

}
