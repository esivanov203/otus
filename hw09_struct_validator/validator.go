package hw09structvalidator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
	Rule  string
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	msg := ""
	for _, err := range v {
		msg += fmt.Sprintf("%s: %s %s", err.Field, err.Rule, err.Err)
	}

	return msg
}

func Validate(v interface{}) error {
	if err := validateStruct(v); err != nil {
		return fmt.Errorf("validate struct: %w", err)
	}

	var errs ValidationErrors

	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)
		tag := field.Tag.Get("validate")
		if tag == "" {
			continue
		}
		errs = append(errs, validateField(field, tag, value)...)
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func validateStruct(v interface{}) error {
	t := reflect.TypeOf(v)
	if t == nil {
		return fmt.Errorf("nil is not allowed")
	}
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return fmt.Errorf("%T is not a struct", v)
	}

	return nil
}

func validateField(field reflect.StructField, tag string, value reflect.Value) ValidationErrors {
	var errs ValidationErrors
	rules := strings.Split(tag, "|")

	switch value.Kind() { //nolint:exhaustive // линтер требует реализации всех типов
	case reflect.String:
		str := value.String()
		for _, rule := range rules {
			if err := validateString(rule, str); err != nil {
				errs = append(errs, ValidationError{
					Field: field.Name,
					Err:   err,
					Rule:  rule,
				})
			}
		}
	case reflect.Int:
		num := int(value.Int())
		for _, rule := range rules {
			if err := validateInt(rule, num); err != nil {
				errs = append(errs, ValidationError{
					Field: field.Name,
					Err:   err,
					Rule:  rule,
				})
			}
		}
	case reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			element := value.Index(i)
			subErrs := validateField(field, tag, element)
			for _, e := range subErrs {
				e.Field = fmt.Sprintf("%s[%d]", field.Name, i)
				errs = append(errs, e)
			}
		}
	default:
		errs = append(errs, ValidationError{
			Field: field.Name,
			Err:   fmt.Errorf("unsupported type: %s", value.Kind()),
		})
	}

	return errs
}

func validateString(rule string, value string) error {
	ruleParts := strings.SplitN(rule, ":", 2)
	if len(ruleParts) != 2 {
		return fmt.Errorf("invalid rule")
	}

	key, param := ruleParts[0], ruleParts[1]
	switch key {
	case "len":
		expected, err := strconv.Atoi(param)
		if err != nil {
			return fmt.Errorf("%w", err)
		}
		if len(value) != expected {
			return fmt.Errorf("must be %d", expected)
		}
	case "regexp":
		re, err := regexp.Compile(param)
		if err != nil {
			return fmt.Errorf("must compile regexp: %w", err)
		}
		if !re.MatchString(value) {
			return fmt.Errorf("must match regexp: %s", re.String())
		}
	case "in":
		opts := strings.Split(param, ",")
		for _, opt := range opts {
			if value == opt {
				return nil
			}
		}
		return fmt.Errorf("must be one of [%s]", strings.Join(opts, ", "))
	default:
		return fmt.Errorf("unknown string rule")
	}

	return nil
}

func validateInt(rule string, value int) error {
	ruleParts := strings.SplitN(rule, ":", 2)
	if len(ruleParts) != 2 {
		return fmt.Errorf("invalid rule")
	}
	key, param := ruleParts[0], ruleParts[1]
	switch key {
	case "min":
		minV, err := strconv.Atoi(param)
		if err != nil {
			return fmt.Errorf("min validate: %w", err)
		}
		if value < minV {
			return fmt.Errorf("min must be >= %d", minV)
		}
	case "max":
		maxV, err := strconv.Atoi(param)
		if err != nil {
			return fmt.Errorf("max validate: %w", err)
		}
		if value > maxV {
			return fmt.Errorf("max must be <= %d", maxV)
		}
	case "in":
		opts := strings.Split(param, ",")
		for _, opt := range opts {
			n, err := strconv.Atoi(opt)
			if err != nil {
				return fmt.Errorf("must contain numbers")
			}
			if value == n {
				return nil
			}
		}
		return fmt.Errorf("must in [%s]", strings.Join(opts, ", "))
	default:
		return fmt.Errorf("unknown int rule")
	}

	return nil
}
