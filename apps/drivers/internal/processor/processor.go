package processor

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"strconv"
)

type DriversCsv struct {
	csvReader *csv.Reader
}

type DriversCsvRow struct {
	Driver   string
	CustId   int
	Location string
	Class    string
	Irating  int
}

func checkCsvReader(header []string) error {
	if header[0] != "DRIVER" || header[1] != "CUSTID" || header[2] != "LOCATION" || header[13] != "CLASS" || header[14] != "IRATING" {
		return fmt.Errorf("invalid header: %v", header)
	}

	return nil
}

func Process(data io.ReadCloser) (*DriversCsv, error) {
	log.Println("Analyzing data")
	csvReader := csv.NewReader(data)

	// Check the header
	header, err := csvReader.Read()
	if err != nil {
		return nil, err
	}

	err = checkCsvReader(header)
	if err != nil {
		return nil, err
	}

	// Return the DriversCsv struct
	return &DriversCsv{csvReader: csvReader}, nil
}

func (d *DriversCsv) Read() (*DriversCsvRow, error) {
	record, err := d.csvReader.Read()
	if err != nil {
		return nil, err
	}

	custId, err := strconv.Atoi(record[1])
	if err != nil {
		return nil, err
	}

	irating, err := strconv.Atoi(record[14])
	if err != nil {
		return nil, err
	}

	return &DriversCsvRow{
		Driver:   record[0],
		CustId:   custId,
		Location: record[2],
		Class:    record[13],
		Irating:  irating,
	}, nil
}
