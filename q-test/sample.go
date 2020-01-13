package main

import "github.com/chunqian/q"

func main() {
	greeting := "hello world"
	pi := 3.14159265359
	cidrs := []string{"10.0.0.0/16", "172.16.0.0/20", "192.168.0.0/24"}

	type Bar struct {
		Baz string
	}

	type Foo struct {
		a string
		b Bar
		c int
	}

	mt := Foo{
		a: "look how pretty this is!",
		b: Bar{
			Baz: "it follow pointers too",
		},
		c: 123,
	}
	fc := func(n float64) bool {
		if n < 1 {
			return true
		}
		return false
	}
	mp := map[string]interface{}{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": mt,
	}
	q.Q(greeting, cidrs, cidrs[1:3], mt, fc(pi), mp)
}
