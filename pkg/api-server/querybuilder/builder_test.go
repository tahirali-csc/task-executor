package querybuilder

import(
	"log"
	"testing"
)

func TestBuilder(t *testing.T){
	values := make(map[string][]string)
	values["status"] = []string{"3"}
	values["buildId"] = []string{"4"}
	// values["sortBy"] = []string{"-createdTs,updatedTs"}
	values["sortBy"] = []string{"-createdTs"}

	colInfo := make(map[string]Column)
	colInfo["buildId"] = NewColumn("build_id", NumberType)
	colInfo["status"] = NewColumn("stauts", NumberType)
	colInfo["createdTs"] = NewColumn("created_ts", NumberType)
	colInfo["updatedTs"] = NewColumn("updated_ts", NumberType)

	clause, _ := GetFilterClause(values, colInfo)
	log.Println(clause)
}