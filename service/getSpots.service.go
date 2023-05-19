package service

import (
	"encoding/json"
	"example/config"
	"example/model"
)

func GetSpotData(lat, long, radius float64, shape string) ([]byte, error) {
	db, err := config.ConnectDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rowData, err := model.ExecuteShapeQuery(shape, db, lat, long, radius)
	if err != nil {
		return nil, err
	}

	newData, err := model.ExecuteThirdQuery(db, lat, long, radius)
	if err != nil {
		return nil, err
	}

	if err = rowData.Err(); err != nil {
		return nil, err
	}

	jsonData, err := json.Marshal(newData)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}
