### Before optimization

```
evg226@MacBook-Air-Evgeny hw10_program_optimization % go test -bench=. -benchmem
goos: darwin
goarch: arm64
pkg: github.com/esivanov203/otus
cpu: Apple M1
BenchmarkGetDomainStat-8               3         334154028 ns/op        343992981 B/op   3045400 allocs/op
BenchmarkGetUsers-8                    5         256179016 ns/op        193769355 B/op   1201015 allocs/op
BenchmarkCountDomains-8               14          86007366 ns/op        139859952 B/op   1844384 allocs/op
PASS
ok      github.com/esivanov203/otus     7.797s
```