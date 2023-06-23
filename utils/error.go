package utils

import (
	"strconv"
	"strings"
)

type multiError struct {
	title string
	errs  []error
}

func (m multiError) Error() string {
	var b strings.Builder
	b.WriteString(m.title)
	b.WriteString(":\n")
	for i, err := range m.errs {
		b.WriteString(strconv.Itoa(i))
		b.WriteString(". ")
		b.WriteString(err.Error())
		b.WriteString("\n")
	}
	return b.String()
}

func NewMultiError(title string, errs ...error) error {
	return &multiError{
		title: title,
		errs:  errs,
	}
}
