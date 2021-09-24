package sqlutils

import (
	"fmt"
	"strconv"
	"strings"
)

type Param struct {
	Key       string
	Value     interface{}
	Condition string
}

type Filter struct {
	params    []Param
	Limit     int
	Offset    int
	WithIndex bool
}

func (f *Filter) Add (p Param) {
	f.params = append(f.params, p)
}

func (f *Filter) PrepareFilter() (string, string, []interface{}) {
	var whereStatement string
	var filterString string
	var limitOffsetString string
	var queryArgs []interface{}
	if len(f.params) > 0 {
		for i, p := range f.params {
			s := strconv.Itoa(i+1)
			filterString += fmt.Sprintf(" and %s%s", p.Key, f.condition(p.Condition,  s))
			queryArgs = append(queryArgs, p.Value)
		}
		filterString = strings.Trim(filterString, " and")
		whereStatement = "where " + filterString
	}
	if f.Limit > 0 {
		limitOffsetString += fmt.Sprintf(" limit $%d ", len(queryArgs) + 1)
		queryArgs = append(queryArgs, f.Limit)
	}

	if f.Offset > 0 {
		limitOffsetString += fmt.Sprintf(" offset $%d ", len(queryArgs) + 1)
		queryArgs = append(queryArgs, f.Offset)
	}

	return whereStatement, limitOffsetString, queryArgs
}

func (f *Filter) condition(condition string, index string) string {
	specialChar := "$"
	conditions := []string{">", "<", "!=", "=", "<=", ">=", "in"}
	result := "="
	if Contains(conditions, condition) {
		result = condition
	}

	if !f.WithIndex {
		specialChar = "?"
		index=""
	}
	if condition == "in" {
		result = fmt.Sprintf(" %s(%s%s) ", result, specialChar, index)
	} else {
		result = fmt.Sprintf(" %s%s%s ", result, specialChar, index)
	}

	return result
}
