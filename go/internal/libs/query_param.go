package libs

import "strings"

type QueryParam struct {
	PageSize int
	PageNumber int
	SortKey string
	SortOrder string
}

func CheckValidSortOrder(sortOrder string) bool {
	sort := strings.ToUpper(sortOrder)
	return sort == "ASC" || sort == "DESC"
}