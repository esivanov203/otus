### До оптимизации

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

### Изучение причин и решения
- Две функции содержат 2 цикла, которые накапливают данные и передают по значнию
    Решение: оставляем только один цикл - в которой считывание из буфера, сравнение в регулярке, запись в мапу с итерацией значения
             (функции и их бенчмарки становятся не нужны - удаляем)
- При построчном чтении создается новый []byte(line) для json.Unmarshal — это лишние аллокации
    Решение: используем потоковый парсер  json.NewDecoder
- strings.Split, regexp.Match - тоже дают лишние аллокации
    Решение: strings.IndexByte и strings.HasSuffix (дали -500 000 аллокаций)

### Результат
```
evg226@MacBook-Air-Evgeny hw10_program_optimization % go test -bench=. -benchmem
goos: darwin
goarch: arm64
pkg: github.com/esivanov203/otus
cpu: Apple M1
BenchmarkGetDomainStat-8               6         184273403 ns/op         2913386 B/op     112117 allocs/op
PASS
ok      github.com/esivanov203/otus     2.268s
```
