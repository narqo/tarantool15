# tarantool15

Faster queries for [go-tnt][] (Go client for tarantool 1.5).

```
% go test ./... -run XX -bench 'BenchmarkQuery_Pack' -benchmem

  goos: darwin
  goarch: amd64
  pkg: github.com/narqo/tarantool15
  BenchmarkQuery_Pack/Select-4      	10000000	       164 ns/op	     256 B/op	       1 allocs/op
  BenchmarkQuery_Pack/tnt.Select-4  	10000000	       166 ns/op	     100 B/op	       5 allocs/op  ☚ :P
  BenchmarkQuery_Pack/Insert-4      	10000000	       175 ns/op	     256 B/op	       1 allocs/op
  BenchmarkQuery_Pack/tnt.Insert-4  	 2000000	       844 ns/op	    1000 B/op	      17 allocs/op  ☚ :(
  BenchmarkQuery_Pack/Update-4      	10000000	       220 ns/op	     256 B/op	       1 allocs/op
  BenchmarkQuery_Pack/tnt.Update-4  	 1000000	      1153 ns/op	    1196 B/op	      23 allocs/op  ☚ :[[
  BenchmarkQuery_Pack/Delete-4      	10000000	       172 ns/op	     256 B/op	       1 allocs/op
  BenchmarkQuery_Pack/tnt.Delete-4  	 2000000	       867 ns/op	    1000 B/op	      17 allocs/op  ☚ :(
  BenchmarkQuery_Pack/Call-4        	10000000	       134 ns/op	     256 B/op	       1 allocs/op
  BenchmarkQuery_Pack/tnt.Call-4    	 1000000	      1024 ns/op	    1464 B/op	      15 allocs/op  ☚ :(
  PASS
  ok  	github.com/narqo/tarantool15	39.076s
```

[go-tnt]: github.com/lomik/go-tnt
