package files

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/triopium/go_utils/pkg/helper"
	"github.com/xuri/excelize/v2"
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

func CSVreadRows(csvFileName string, csvDelim rune) ([][]string, error) {
	file, err := os.Open(csvFileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.Comma = csvDelim
	reader.LazyQuotes = true
	return reader.ReadAll()
}

func CSVtoXLSX(csvFile string, csvDelim rune) error {
	// read csv
	records, err := CSVreadRows(csvFile, csvDelim)
	if err != nil {
		return err
	}
	fmt.Println(len(records))
	// Create a new Excel file
	xlsxFile := excelize.NewFile()
	sheetName := "Sheet1"
	// Create a new sheet
	index, err := xlsxFile.NewSheet(sheetName)
	if err != nil {
		return err
	}
	xlsxFile.SetActiveSheet(index)
	for i, row := range records {
		for j, cell := range row {
			cellRef, _ := excelize.CoordinatesToCellName(j+1, i+1)
			err := xlsxFile.SetCellValue(sheetName, cellRef, cell)
			if err != nil {
				return err
			}
		}
	}
	name := helper.FilenameWithoutExtension(csvFile)
	dir := filepath.Dir(csvFile)
	xlsxFilePath := filepath.Join(dir, name+".xlsx")
	return xlsxFile.SaveAs(xlsxFilePath)
}

func CSVdirToXLSX(csvFolder string, csvDelim rune) error {
	files, err := helper.ListDirFiles(csvFolder)
	if err != nil {
		return err
	}
	for _, f := range files {
		ext := filepath.Ext(f)
		if ext == ".csv" {
			// fmt.Println(filepath.Base(f))
			fmt.Println(f)
			rows, err := CSVreadRows(f, csvDelim)
			if err != nil {
				return fmt.Errorf("%w file: %s", err, f)
			}
			fmt.Println(len(rows))
			err = CSVtoXLSX(f, csvDelim)
			if err != nil {
				return fmt.Errorf("%w file: %s", err, f)
			}
		}
	}
	// fmt.Println(files)
	return nil
}
