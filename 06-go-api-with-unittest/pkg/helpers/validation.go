package helpers

import (
	"06-go-api-with-unittest/pkg/errs"
	"reflect"

	"github.com/asaskevich/govalidator"
)

func ValidateStruct(data interface{}) errs.ErrorMessage {
	_, err := govalidator.ValidateStruct(data)
	if err != nil {
		return errs.NewBadRequestError(err.Error())
	}

	err = validateStructValue(data)
	if err != nil {
		return err.(errs.ErrorMessage)
	}

	return nil
}

func validateStructValue(data interface{}) errs.ErrorMessage {
	reflectValue := reflect.ValueOf(data)
	reflectType := reflect.TypeOf(data)

	for i := 0; i < reflectValue.NumField(); i++ {
		fieldValue := reflectValue.Field(i)
		fieldType := reflectType.Field(i)

		switch fieldValue.Kind() {
		case reflect.Float64:
			if fieldValue.Float() <= 0 {
				errMsg := fieldType.Name + " must be greater than 0"
				return errs.NewBadRequestError(errMsg)
			}
		}
	}

	return nil
}
