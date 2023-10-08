package generated
import (
)
func (io *PhoneNumberSortInput) SortSqlStrings() []string {
	if io == nil {
		return []string{}
	}
	var sortSqlStrings []string
	for _, si := range io.By {
		str := si.Field.String() + " " + si.Direction.String()
		sortSqlStrings = append(sortSqlStrings, str)
	}
	return sortSqlStrings
}
func (io *UserSortInput) SortSqlStrings() []string {
	if io == nil {
		return []string{}
	}
	var sortSqlStrings []string
	for _, si := range io.By {
		str := si.Field.String() + " " + si.Direction.String()
		sortSqlStrings = append(sortSqlStrings, str)
	}
	return sortSqlStrings
}