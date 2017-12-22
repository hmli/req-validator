package req_validator

import (
	"testing"
	"net/http"
)

func TestValidator_Validate(t *testing.T) {
	req, err  := http.NewRequest("GET", "http://localhost?a=1&b=2&c=ddd", nil)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	v := New(map[string]string{ // true
		"a": "required|int",
		"b": "int",
	})
	t.Log(v.Validate(req))
	v2 := New(map[string]string{ // f
		"c": "datefmt:2006-01-02 15:04:05",
	})
	t.Log(v2.Validate(req))
	v3 := New(map[string]string{ // f
		"d": "required",
	})
	t.Log(v3.Validate(req))

	v4 := New(map[string]string{ // t
		"a": "bool",
	})
	t.Log(v4.Validate(req))
	v5 := New(map[string]string{ // f
		"b": "bool",
	})
	t.Log(v5.Validate(req))
	v6 := New(map[string]string{ // t
		"b": "enum:1,2,3",
	})
	t.Log(v6.Validate(req))
	v7 := New(map[string]string{ // f
		"b": "enum:1,3",
	})
	t.Log(v7.Validate(req))
	v8 := New(map[string]string{ // f
		"b": "between:3,6",
	})
	t.Log(v8.Validate(req))
	v9 := New(map[string]string{ // t
		"b": "between:1,3",
	})
	t.Log(v9.Validate(req))
}
