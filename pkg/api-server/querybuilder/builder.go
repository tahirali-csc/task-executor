package querybuilder

// package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type ColumnType string

var StringType ColumnType = "string"
var NumberType ColumnType = "number"
var TimestampType ColumnType = "timestamp"

type Column struct {
	Name string
	Type ColumnType
}

func NewColumn(name string, fieldType ColumnType) Column {
	return Column{Name: name, Type: fieldType}
}

func GetFilterClause(values map[string][]string, columnInfo map[string]Column) (string, error) {
	var cols []string
	var sortCols []string
	const SortBy = "sortBy"

	for k, v := range values {
		if k != SortBy {
			if info, ok := columnInfo[k]; ok {
				switch info.Type {
				case StringType:
					cols = append(cols, fmt.Sprintf("%s='%s'", info.Name, v[0]))
				case NumberType:
					v, _ := strconv.ParseInt(v[0], 10, 64)
					cols = append(cols, fmt.Sprintf("%s=%d", info.Name, v))
				}
			} else {
				return "", errors.New(fmt.Sprintf("%s is undefined", k))
			}
		} else {
			sortParts := strings.Split(v[0], ",")
			for _, c := range sortParts {

				sortDir := "ASC"
				col := c

				if strings.HasPrefix(col, "-") {
					col = strings.Replace(col, "-", "", 1)
					sortDir = "DESC"
				}

				if info, ok := columnInfo[col]; ok {
					sortCols = append(sortCols, fmt.Sprintf("ORDER BY %s %s", info.Name, sortDir))
				}
			}
		}
	}

	whereClause := strings.Join(cols, " AND ")
	sortClause := strings.Join(sortCols, ",")

	return fmt.Sprintf("WHERE %s\n%s", whereClause, sortClause), nil
}
