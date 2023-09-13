package mapper

import "fmt"

type DestinationNotPointerError struct {
	error
}

func NewDestinationNotPointerError(message string, args ...any) *DestinationNotPointerError {
	return &DestinationNotPointerError{fmt.Errorf(message, args...)}
}

type SourceNotStructError struct {
	error
}

func NewSourceNotStructError(message string, args ...any) *SourceNotStructError {
	return &SourceNotStructError{fmt.Errorf(message, args...)}
}

type DestinationNotStructError struct {
	error
}

func NewDestinationNotStructError(message string, args ...any) *DestinationNotStructError {
	return &DestinationNotStructError{fmt.Errorf(message, args...)}
}

type DestinationFieldNotFoundError struct {
	error
}

func NewDestinationFieldNotFoundError(message string, args ...any) *DestinationFieldNotFoundError {
	return &DestinationFieldNotFoundError{fmt.Errorf(message, args...)}
}

type DestinationFieldNotSettableError struct {
	error
}

func NewDestinationFieldNotSettableError(message string, args ...any) *DestinationFieldNotSettableError {
	return &DestinationFieldNotSettableError{error: fmt.Errorf(message, args...)}
}

type FieldsTypesMismatchError struct {
	error
}

func NewFieldsTypesMismatchError(message string, args ...any) *FieldsTypesMismatchError {
	return &FieldsTypesMismatchError{error: fmt.Errorf(message, args...)}
}

type SourceFieldIsNotValidError struct {
	error
}

func NewSourceFieldIsNotValidError(message string, args ...any) *SourceFieldIsNotValidError {
	return &SourceFieldIsNotValidError{error: fmt.Errorf(message, args...)}
}

type DestFieldIsNotValidError struct {
	error
}

func NewDestFieldIsNotValidError(message string, args ...any) *DestFieldIsNotValidError {
	return &DestFieldIsNotValidError{error: fmt.Errorf(message, args...)}
}

type InvalidMaptoTagError struct {
	error
}

func NewInvalidMaptoTagError(message string, args ...any) *InvalidMaptoTagError {
	return &InvalidMaptoTagError{error: fmt.Errorf(message, args...)}
}
