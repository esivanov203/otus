package hw10programoptimization

import (
	"archive/zip"
	"testing"
)

func BenchmarkGetDomainStat(b *testing.B) {
	b.StopTimer()
	r, err := zip.OpenReader("testdata/users.dat.zip")
	if err != nil {
		b.Fatal(err)
	}
	defer func() { _ = r.Close() }()
	for i := 0; i < b.N; i++ {
		data, err := r.File[0].Open()
		if err != nil {
			b.Fatal(err)
		}

		b.StartTimer()
		_, err = GetDomainStat(data, "biz")
		b.StopTimer()

		if err != nil {
			b.Fatal(err)
		}
		err = data.Close()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGetUsers(b *testing.B) {
	b.StopTimer()
	r, err := zip.OpenReader("testdata/users.dat.zip")
	if err != nil {
		b.Fatal(err)
	}
	defer func() { _ = r.Close() }()
	for i := 0; i < b.N; i++ {
		data, err := r.File[0].Open()
		if err != nil {
			b.Fatal(err)
		}

		b.StartTimer()
		_, err = getUsers(data)
		b.StopTimer()

		if err != nil {
			b.Fatal(err)
		}
		err = data.Close()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCountDomains(b *testing.B) {
	b.StopTimer()
	r, err := zip.OpenReader("testdata/users.dat.zip")
	if err != nil {
		b.Fatal(err)
	}
	defer func() { _ = r.Close() }()
	data, err := r.File[0].Open()
	defer func() { _ = data.Close() }()
	u, err := getUsers(data)
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {

		b.StartTimer()
		_, err = countDomains(u, "biz")
		b.StopTimer()

		if err != nil {
			b.Fatal(err)
		}
	}
}
