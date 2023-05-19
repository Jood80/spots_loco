package service

import (
	"encoding/json"
	"example/config"
	"example/database"
)

func GetSpotData(lat, long, radius float64, shape string) ([]byte, error) {
	db, err := config.ConnectDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rowData, err := database.ExecuteShapeQuery(shape, db, lat, long, radius)
	if err != nil {
		return nil, err
	}

	newData, err := database.ExecuteThirdQuery(db, lat, long, radius)
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
