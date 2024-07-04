package utils

import "gorm.io/gorm/clause"

func Returning(field ...string) clause.Returning {
	columns := make([]clause.Column, len(field))
	for i, f := range field {
		columns[i] = clause.Column{Name: f}
	}
	return clause.Returning{Columns: columns}
}
