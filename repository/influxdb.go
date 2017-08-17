package repository

import (
	"time"

	client "github.com/influxdata/influxdb/client/v2"
)

type InfluxDB struct {
	client client.Client
}

func NewInfluxDBClient(host string) (*InfluxDB, error) {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: host,
	})
	if err != nil {
		return nil, err
	}

	idb := &InfluxDB{
		client: c,
	}
	return idb, nil
}

func (idb *InfluxDB) Save(table string, data map[string]interface{}) error {

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "tests",
		Precision: "s",
	})

	if err != nil {
		return err
	}

	pt, err := client.NewPoint(table, nil, data, time.Now())
	if err != nil {
		return err
	}

	bp.AddPoint(pt)

	if err := idb.client.Write(bp); err != nil {
		return err
	}

	return nil
}
