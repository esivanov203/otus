package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		Levels []int    `validate:"in:1,4,5"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	Name struct {
		Name string `validate:"in-Jake"`
	}

	Salary struct {
		Salary int `validate:"max=5000"`
	}

	Unknown struct {
		Name  string `validate:"no-string:ok"`
		Count int    `validate:"no-int:ok"`
	}

	NoType struct {
		Height float64 `validate:"min:10"`
	}

	EmptyTag struct {
		Name string `validate:""`
	}

	InRegexp struct {
		Name  string `validate:"in:admin,stuff"`
		Login string `validate:"regexp:("`
	}

	MinMaxBorder struct {
		Min int `validate:"min:10"`
		Max int `validate:"max:100"`
	}

	InNotANumber struct {
		count int `validate:"in:two,10"`
	}

	EmptySlice struct {
		Jobs []string `validate:"len:11"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name string
		in   interface{}
	}{
		{"empty struct", struct{}{}},
		{"empty tag", EmptyTag{Name: "ok"}},
		{"empty slice element", EmptySlice{[]string{}}},
		{"simple string len", App{"1.0.0"}},
		{"struct by ptr", &App{"1.0.0"}},
		{"min max border", MinMaxBorder{10, 100}},
		{
			"no any tags in struct",
			Token{
				Header:    []byte("header-data"),
				Payload:   []byte("payload-data"),
				Signature: []byte("signature-data"),
			},
		},
		{"struct has other tag", Response{200, "body"}},
		{
			"complex struct",
			User{
				ID:     "123e4567-e89b-12d3-a456-426614174000",
				Name:   "Alice",
				Age:    30,
				Email:  "alice@example.com",
				Role:   "admin",
				Phones: []string{"79991234567", "79997654321"},
				Levels: []int{1, 5},
				meta:   json.RawMessage(`{"active":true}`),
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("positive case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			result := Validate(tt.in)
			require.Equal(t, nil, result)
			_ = tt
		})
	}
}

func TestValidateNegative(t *testing.T) {
	tests := []struct {
		name         string
		in           interface{}
		expectedErrs []string
	}{
		{"nil", nil, []string{"nil is not allowed"}},
		{"not struct", "", []string{"is not a struct"}},
		{"simple string len", App{"1.0"}, []string{"len"}},
		{"string in regexp", InRegexp{"ok", "Joe"}, []string{"in", "regexp"}},
		{"int is not in", Response{700, "body"}, []string{"in"}},
		{"int in not int", InNotANumber{10}, []string{"in", "must contain numbers"}},
		{
			"complex struct",
			User{
				ID:     "123e4567-e89b-12d3-a456-426614174000",
				Name:   "Alice",
				Age:    17,      // min int
				Email:  "alice", // regexp string
				Role:   "adm",   // in string
				Phones: []string{"79991234567", "79997654321"},
				meta:   json.RawMessage(`{"active":true}`),
			},
			[]string{"min", "regexp", "in"},
		},
		{"string rule not valid", Name{"Jake"}, []string{"invalid rule"}},
		{"between int rule not valid", Salary{4000}, []string{"invalid rule"}},
		{"unknown rules", Unknown{"Jake", 20}, []string{"unknown int rule", "unknown string rule"}},
		{"unsupported type", NoType{15.2}, []string{"unsupported type"}},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("negative case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			result := Validate(tt.in)
			require.Error(t, result)
			for _, expectedErr := range tt.expectedErrs {
				require.Contains(t, result.Error(), expectedErr)
			}
		})
	}
}
