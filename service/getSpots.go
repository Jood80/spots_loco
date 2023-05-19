package service

import (
	"encoding/json"
	"example/config"
	"example/database"
)

func GetSpotData(latFlo, longFlo, radiusFlo float64, shape string) ([]byte, error) {
	db, err := config.ConnectDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rowData, err := database.ExecuteShapeQuery(shape, db, latFlo, longFlo, radiusFlo)
	if err != nil {
		return nil, err
	}

	newData, err := database.ExecuteThirdQuery(db, latFlo, longFlo, radiusFlo)
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
