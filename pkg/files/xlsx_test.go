package files

import (
	"fmt"
	"strings"
	"testing"
)

func TestSlice(t *testing.T) {
	my := []string{"a", "b", "c", "d", "e"}
	mya := my[1:3]
	fmt.Println(mya[0:2])
	fmt.Println(len(mya))
	fmt.Println(len(my))
}

func TestTable_MapTableHeaders(t *testing.T) {
	filterFile := "/home/jk/CRO/CRO_BASE/openmedia_backup/filters/eurovolby - zadání.xlsx"
	rows, err := ReadExcelFileSheetRows(filterFile, "data")
	if err != nil {
		t.Error(err)
	}
	table := CreateTable(rows, 0, 0)
	fmt.Println(len(table.RowHeaderToColumnMap))
	fmt.Println(table.RowHeaderToColumnMap)
	fmt.Println(len(table.ColumnHeaderMap))
	// match := table.MatchRow("Pochman Stanislav", "navrhující strana", "KAN")
	match := table.MatchRow("Rohel Petr", "navrhující strana", "Levice")
	// match := table.MatchRow("Juřica Vojtěch", "navrhující strana", "Levice")
	// match := table.MatchRow("Široký Jan", "navrhující strana", "Levice")
	fmt.Println(match)
}

func TestXLSXtableBuild(t *testing.T) {
	type args struct {
		fileName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"simple", args{"/tmp/kek.xlsx"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := XLSXtableBuild(tt.args.fileName); (err != nil) != tt.wantErr {
				t.Errorf("XLSXtableBuild() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestXLSXtableStreamSave(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"simple", args{"/tmp/sream.xlsx"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := XLSXtableStreamSave(tt.args.filePath); (err != nil) != tt.wantErr {
				t.Errorf("XLSXtableStreamSave() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateTableTransformColumn(t *testing.T) {
	rows := [][]string{
		{"Ahoj", "Mahoj", "Sumak"}, // Column header
		{"Tahoj", "Bahoj", "Tumak"},
		{"Kahoj", "Čahoj", "Kakak"},
		{"Žake", "Lahoj", "Rakak"},
	}
	table := CreateTableTransformColumn(rows, 0, 0, strings.ToLower)
	fmt.Println(table.RowHeaderToColumnMap)
}
