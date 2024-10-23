package handler

import (
	"encoding/json"
	"fmt"
)

type Query struct {
	Method    string        `json:"method"`
	Attribute string        `json:"attribute,omitempty"`
	Values    []interface{} `json:"values,omitempty"`
}

func (q Query) String() string {
	jsonData, _ := json.Marshal(q)
	return string(jsonData)
}

func Equal(attribute string, value interface{}) string {
	return newQuery("equal", attribute, parseValue(value)).String()
}

func NotEqual(attribute string, value interface{}) string {
	return newQuery("notEqual", attribute, parseValue(value)).String()
}

func LessThan(attribute string, value interface{}) string {
	return newQuery("lessThan", attribute, parseValue(value)).String()
}

func LessThanEqual(attribute string, value interface{}) string {
	return newQuery("lessThanEqual", attribute, parseValue(value)).String()
}

func GreaterThan(attribute string, value interface{}) string {
	return newQuery("greaterThan", attribute, parseValue(value)).String()
}

func GreaterThanEqual(attribute string, value interface{}) string {
	return newQuery("greaterThanEqual", attribute, parseValue(value)).String()
}

func Search(attribute string, value string) string {
	return newQuery("search", attribute, []interface{}{value}).String()
}

func IsNull(attribute string) string {
	return newQuery("isNull", attribute, nil).String()
}

func IsNotNull(attribute string) string {
	return newQuery("isNotNull", attribute, nil).String()
}

func Between(attribute string, start interface{}, end interface{}) string {
	return newQuery("between", attribute, []interface{}{start, end}).String()
}

func StartsWith(attribute string, value string) string {
	return newQuery("startsWith", attribute, []interface{}{value}).String()
}

func EndsWith(attribute string, value string) string {
	return newQuery("endsWith", attribute, []interface{}{value}).String()
}

func Select(attributes []string) string {
	return newQuery("select", "", toInterfaceSlice(attributes)).String()
}

func OrderAsc(attribute string) string {
	return newQuery("orderAsc", attribute, nil).String()
}

func OrderDesc(attribute string) string {
	return newQuery("orderDesc", attribute, nil).String()
}

func CursorBefore(documentID string) string {
	return newQuery("cursorBefore", "", []interface{}{documentID}).String()
}

func CursorAfter(documentID string) string {
	return newQuery("cursorAfter", "", []interface{}{documentID}).String()
}

func Limit(limit int) string {
	return newQuery("limit", "", []interface{}{limit}).String()
}

func Offset(offset int) string {
	return newQuery("offset", "", []interface{}{offset}).String()
}

func Contains(attribute string, value interface{}) string {
	return newQuery("contains", attribute, parseValue(value)).String()
}

func Or(queries []string) string {
	return newQuery("or", "", parseQueries(queries)).String()
}

func And(queries []string) string {
	return newQuery("and", "", parseQueries(queries)).String()
}

// Helper functions
func newQuery(method, attribute string, values []interface{}) Query {
	return Query{
		Method:    method,
		Attribute: attribute,
		Values:    values,
	}
}

func parseValue(value interface{}) []interface{} {
	switch v := value.(type) {
	case []interface{}:
		return v
	default:
		return []interface{}{value}
	}
}

func toInterfaceSlice(strings []string) []interface{} {
	interfaces := make([]interface{}, len(strings))
	for i, s := range strings {
		interfaces[i] = s
	}
	return interfaces
}

func parseQueries(queries []string) []interface{} {
	parsed := make([]interface{}, len(queries))
	for i, q := range queries {
		var query Query
		err := json.Unmarshal([]byte(q), &query)
		if err != nil {
			fmt.Println("Error parsing query:", err)
			return nil
		}
		parsed[i] = query
	}
	return parsed
}
