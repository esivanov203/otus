package hw10programoptimization

import (
	"bufio"
	"encoding/json"
	"io"
	"strings"
)

type User struct {
	Email string `json:"email"`
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	ds := make(DomainStat)
	dotDomain := "." + domain

	br := bufio.NewReader(r)
	dec := json.NewDecoder(br)

	var user User
	for {
		err := dec.Decode(&user)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		email := user.Email
		if strings.HasSuffix(user.Email, dotDomain) {
			at := strings.IndexByte(user.Email, '@')
			if at < 0 {
				continue
			}
			dom := email[at+1:]
			key := strings.ToLower(dom)
			ds[key]++
		}
	}

	return ds, nil
}
