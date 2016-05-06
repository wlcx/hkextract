package main

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"os"
	"time"
)

type xmlDateTime struct {
	time.Time
}

func (t *xmlDateTime) UnmarshalText(text []byte) error {
	parsed, err := time.Parse("2006-01-02 15:04:05 -0700", string(text))
	if err != nil {
		return err
	}
	*t = xmlDateTime{parsed}
	return nil
}

type HealthData struct {
	XMLName    xml.Name `xml:"HealthData"`
	ExportDate string
	Records    []Record `xml:"Record"`
}

type Record struct {
	Type      string      `xml:"type,attr"`
	Unit      string      `xml:"unit,attr"`
	Value     string      `xml:"value,attr"`
	StartDate xmlDateTime `xml:"startDate,attr"`
}

const (
	TypeIdentifierHeight   = "HKQuantityTypeIdentifierHeight"
	TypeIdentifierBodyMass = "HKQuantityTypeIdentifierBodyMass"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s path/to/data.xml\n", os.Args[0])
	}
	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer f.Close()
	decoder := xml.NewDecoder(f)
	export := HealthData{}
	if err := decoder.Decode(&export); err != nil {
		panic(err)
	}
	w := csv.NewWriter(os.Stdout)
	w.Write([]string{"date", "weight"})
	for _, r := range export.Records {
		if r.Type == TypeIdentifierBodyMass {
			w.Write([]string{r.StartDate.Format("2006-01-02"), r.Value})
		}
	}
	w.Flush()
}
