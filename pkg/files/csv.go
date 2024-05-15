package files

import (
	"encoding/csv"
	"fmt"
	"os"
	"reflect"
)

func CSVcompareRows(fileName1, fileName2 string) {
	file, err := os.Open(fileName1)
	if err != nil {
		fmt.Println(err)
	}
	reader := csv.NewReader(file)
	records, _ := reader.ReadAll()
	fmt.Println(records)
	file2, err := os.Open(fileName2)
	if err != nil {
		fmt.Println(err)
	}
	reader2 := csv.NewReader(file2)
	records2, _ := reader2.ReadAll()
	fmt.Println(records2)
	allrs := reflect.DeepEqual(records, records2)
	fmt.Println(allrs)
}
