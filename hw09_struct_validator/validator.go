package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field   string
	Rule    string
	Message string
}

func (err ValidationError) Error() string {
	return fmt.Sprintf("%s: %s (%s)", err.Field, err.Message, err.Rule)
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var sb strings.Builder
	for _, e := range v {
		sb.WriteString(e.Error())
		sb.WriteString("; ")
	}
	return strings.TrimSpace(sb.String())
}

type InternalValidationError struct {
	Message string
}

func (e InternalValidationError) Error() string {
	return fmt.Sprintf("validator internal: %s", e.Message)
}

func Validate(v interface{}) error {
	if err := validateStruct(v); err != nil {
		return InternalValidationError{Message: err.Error()}
	}

	var vErrs ValidationErrors

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

		fvErrs, err := validateField(field, tag, value)
		if err != nil {
			var internalErr InternalValidationError
			if errors.As(err, &internalErr) {
				return internalErr
			}
			return err
		}
		vErrs = append(vErrs, fvErrs...)
	}

	if len(vErrs) > 0 {
		return vErrs
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

func validateField(field reflect.StructField, tag string, value reflect.Value) (ValidationErrors, error) {
	var errs ValidationErrors
	rules := strings.Split(tag, "|")

	switch value.Kind() { //nolint:exhaustive
	case reflect.String:
		str := value.String()
		for _, rule := range rules {
			vErrMsg, err := validateString(rule, str)
			if err != nil {
				return nil, InternalValidationError{Message: fmt.Sprintf("field %s: %v", field.Name, err)}
			}
			if vErrMsg != "" {
				errs = append(errs, ValidationError{
					Field:   field.Name,
					Rule:    rule,
					Message: vErrMsg,
				})
			}
		}
	case reflect.Int:
		num := int(value.Int())
		for _, rule := range rules {
			vErrMsg, err := validateInt(rule, num)
			if err != nil {
				return nil, InternalValidationError{Message: fmt.Sprintf("field %s: %v", field.Name, err)}
			}
			if vErrMsg != "" {
				errs = append(errs, ValidationError{
					Field:   field.Name,
					Rule:    rule,
					Message: vErrMsg,
				})
			}
		}
	case reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			elem := value.Index(i)
			vErrs, err := validateField(field, tag, elem)
			if err != nil {
				var internalErr InternalValidationError
				if errors.As(err, &internalErr) {
					return nil, internalErr
				}
				return nil, err
			}
			for _, e := range vErrs {
				e.Field = fmt.Sprintf("%s[%d]", e.Field, i)
				errs = append(errs, e)
			}
		}
	default:
		return nil, InternalValidationError{
			Message: fmt.Sprintf("unsupported field type %s for field %s", value.Kind(), field.Name),
		}
	}

	return errs, nil
}

func validateString(rule string, value string) (string, error) {
	ruleParts := strings.SplitN(rule, ":", 2)
	if len(ruleParts) != 2 {
		return "", fmt.Errorf("invalid rule format: %q", rule)
	}

	key, param := ruleParts[0], ruleParts[1]
	switch key {
	case "len":
		expected, err := strconv.Atoi(param)
		if err != nil {
			return "", fmt.Errorf("must number: %w", err)
		}
		if len(value) != expected {
			return fmt.Sprintf("length must be %d", expected), nil
		}

	case "regexp":
		re, err := regexp.Compile(param)
		if err != nil {
			return "", fmt.Errorf("invalid regexp: %w", err)
		}
		if !re.MatchString(value) {
			return fmt.Sprintf("must match regexp %q", param), nil
		}

	case "in":
		opts := strings.Split(param, ",")
		for _, opt := range opts {
			if value == opt {
				return "", nil
			}
		}
		return fmt.Sprintf("must be one of [%s]", strings.Join(opts, ", ")), nil

	default:
		return "", fmt.Errorf("unknown rule: %q", key)
	}

	return "", nil
}

func validateInt(rule string, value int) (string, error) {
	ruleParts := strings.SplitN(rule, ":", 2)
	if len(ruleParts) != 2 {
		return "", fmt.Errorf("invalid rule format: %q", rule)
	}

	key, param := ruleParts[0], ruleParts[1]
	switch key {
	case "min":
		minV, err := strconv.Atoi(param)
		if err != nil {
			return "", fmt.Errorf("must number: %w", err)
		}
		if value < minV {
			return fmt.Sprintf("must be >= %d", minV), nil
		}
	case "max":
		maxV, err := strconv.Atoi(param)
		if err != nil {
			return "", fmt.Errorf("must number: %w", err)
		}
		if value > maxV {
			return fmt.Sprintf("must be <= %d", maxV), nil
		}
	case "in":
		opts := strings.Split(param, ",")
		for _, opt := range opts {
			n, err := strconv.Atoi(opt)
			if err != nil {
				return "", fmt.Errorf("must number: %w", err)
			}
			if value == n {
				return "", nil
			}
		}
		return fmt.Sprintf("must be one of [%s]", strings.Join(opts, ", ")), nil
	default:
		return "", fmt.Errorf("unknown rule: %q", key)
	}

	return "", nil
}
