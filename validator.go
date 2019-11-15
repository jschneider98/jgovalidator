package jgovalidator


import (
	"time"
	"regexp"
	"database/sql"
	"database/sql/driver"
	"reflect"
	"gopkg.in/go-playground/validator.v9"
)

// 
type Validate struct {
	*validator.Validate
}

const (
	Int               string = "^(?:[-+]?(?:0|[1-9][0-9]*))$"
	Float             string = "^(?:[-+]?(?:[0-9]+))?(?:\\.[0-9]*)?(?:[eE][\\+\\-]?(?:[0-9]+))?$"
	DateExp           string = `\d{4}-\d{2}-\d{2}`
	RF3339WithoutZone string = "2006-01-02T15:04:05"
)

// Singleton. This object caches struct defs, so use singleton.
var (
	validate *validator.Validate
	rxInt = regexp.MustCompile(Int)
	rxFloat = regexp.MustCompile(Float)
	rxDate = regexp.MustCompile(DateExp)
)

//
func GetValidator() *validator.Validate {

	if validate != nil {
		return validate
	}

	validate = validator.New()
	validate.RegisterValidation("notNull", NotNull)
	validate.RegisterValidation("int", IsInt)
	validate.RegisterValidation("float", IsFloat)
	validate.RegisterValidation("date", IsDate)
	validate.RegisterValidation("rfc3339", IsRFC3339)
	validate.RegisterValidation("rfc3339WithoutZone", IsRFC3339WithoutZone)
	validate.RegisterValidation("datetime", IsDatetime)

	// register all sql.Null* types to use the ValidateValuer CustomTypeFunc
	validate.RegisterCustomTypeFunc(ValidateValuer, sql.NullString{}, sql.NullInt64{}, sql.NullBool{}, sql.NullFloat64{})

	return validate
}

// ValidateValuer implements validator.CustomTypeFunc
func ValidateValuer(field reflect.Value) interface{} {

	if valuer, ok := field.Interface().(driver.Valuer); ok {

		val, err := valuer.Value()
		if err == nil {
			return val
		}
	}

	return nil
}

// Used with sql.Null* types. If we get here, the value is not null, so just return true
func NotNull(fl validator.FieldLevel) bool {
	return true
}

// IsInt check if the string is an integer. Empty string is valid.
func IsInt(fl validator.FieldLevel) bool {
	if IsNull(fl.Field().String()) {
		return false
	}
	return rxInt.MatchString(fl.Field().String())
}

// IsFloat check if the string is a float.
func IsFloat(fl validator.FieldLevel) bool {
	return fl.Field().String() != "" && rxFloat.MatchString(fl.Field().String())
}

// IsRFC3339 check if string is valid timestamp value according to RFC3339
func IsDate(fl validator.FieldLevel) bool {
	if IsNull(fl.Field().String()) {
		return false
	}

	return rxDate.MatchString(fl.Field().String())
}

// IsRFC3339 check if string is valid timestamp value according to RFC3339
func IsRFC3339(fl validator.FieldLevel) bool {
	return IsTime(fl.Field().String(), time.RFC3339)
}

// IsRFC3339WithoutZone check if string is valid timestamp value according to RFC3339 which excludes the timezone.
func IsRFC3339WithoutZone(fl validator.FieldLevel) bool {
	return IsTime(fl.Field().String(), RF3339WithoutZone)
}

// datetime with or without timezone
func IsDatetime(fl validator.FieldLevel) bool {
	return ( IsTime(fl.Field().String(), time.RFC3339) || IsTime(fl.Field().String(), RF3339WithoutZone) )
}

// **********************************

// IsNull check if the string is null.
func IsNull(str string) bool {
	return len(str) == 0
}

// IsTime check if string is valid according to given format
func IsTime(str string, format string) bool {
	_, err := time.Parse(format, str)
	return err == nil
}
