package hw02unpackstring

import (
	"errors"
	"strconv"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	sr := []rune(s)
	ret := make([]rune, 0, len(sr))
	isPrevNum := false

	for i, c := range sr {
		isCurrentNum := unicode.IsNumber(c)
		if isCurrentNum {
			if i == 0 || isPrevNum {
				return "", ErrInvalidString
			}
			num, err := strconv.Atoi(string(c))
			if err != nil {
				return "", ErrInvalidString
			}
			ret = ret[:len(ret)-1]
			for j := 0; j < num; j++ {
				ret = append(ret, sr[i-1])
			}
		} else {
			ret = append(ret, sr[i])
		}
		isPrevNum = isCurrentNum
	}

	return string(ret), nil
}
