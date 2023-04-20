package main

import "github.com/xuri/excelize/v2"

type matchApiPostData struct {
	FieldHeaders []string `json:"fieldHeaders"`
}

func getTitles(excelPath string) ([]string, error) {
	xlData, err := excelize.OpenFile(excelPath)
	if err != nil {
		return nil, err
	}
	rows, err := xlData.GetRows("Sheet1")
	if err != nil {
		return nil, err
	}
	titles := rows[0]
	return titles, nil
}
