package service

import (
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"

	"github.com/sacurio/jb-challenge/internal/app/model"
)

type Stock interface {
	Parse() (model.Stock, error)
	String() string
}

type stock struct {
	data  string
	model model.Stock
}

func NewStock(data string) Stock {
	return &stock{
		data: data,
	}
}

func (s *stock) Parse() (model.Stock, error) {
	reader := csv.NewReader(strings.NewReader(s.data))

	_, _ = reader.Read()
	var st model.Stock

	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		st = model.Stock{
			Symbol: record[0],
			Date:   record[1],
			Time:   record[2],
		}

		st.Open, _ = strconv.ParseFloat(record[3], 64)
		st.High, _ = strconv.ParseFloat(record[4], 64)
		st.Low, _ = strconv.ParseFloat(record[5], 64)
		st.Close, _ = strconv.ParseFloat(record[6], 64)
		st.Volume, _ = strconv.Atoi(record[7])
	}

	s.model = st
	return st, nil
}

func (s *stock) String() string {
	return fmt.Sprintf("%s quote is $%.2f per share.", strings.ToUpper(s.model.Symbol), s.model.High)
}
