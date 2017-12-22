package req_validator

import (
	"errors"
	"strings"
	"time"
	"fmt"
	"strconv"
)

type Rule interface{
	Validate(key, tagValue, value string) (bool, error)
}

var Rules = map[string]Rule{
	"notnull": &NotnullRule{},
	"int": &IntRule{}, // "123", "1"
	"required": &NotnullRule{}, // 因为 r.Form.Get(value) 如果取不到值就会返回 "", 结果和 null 是一样的,
	"bool": &BooleanRule{}, // "1" "0" "true" "False" "TRUE"
	"datefmt": &DatefmtRule{},
	"enum": &EnumRule{}, // enum:2,34,sdf,23
	"between": &BetweenRule{}, // between:1,10, must be int
}

func RegisterRule(name string, r Rule) {
	Rules[name] = r
}

type NotnullRule struct {}

type RequiredRule NotnullRule

// Validate: key 参数名, tagValue tag的取值,如果有  value 参数值
func (r *NotnullRule) Validate(key, tagValue, value string) (bool, error) {
	if value == "" {
		return false, errors.New("parameter '" + key + "' is required")
	}
	return true, nil
}

func (r *NotnullRule) Error() error { return nil }


type IntRule struct {}

func (d *IntRule) Validate(key, tagValue, value string) (bool, error) {
	if value == "" {
		return true, nil
	}
	for i:= 0; i<len(value); i ++ {
		if !isDigit(value[i]) {
			return false, errors.New("parameter '" + key + "' with value '" + value + "' is not an integer")
		}
	}
	return true, nil
}

func (d *IntRule) Error() error { return nil }

func isDigit(b byte) bool {
	return uint8(b) > 47 && uint8(b) < 58
}

type BooleanRule struct {}

var booleans = []string{"1", "0", "true", "false"}

func (b *BooleanRule) Validate(key, tagValue, value string) (bool, error) {
	if value == "" || value == "1" || value == "0" {
		return true, nil
	} else {
		value = strings.ToLower(value)
		if value == "true" || value == "false" {
			return true, nil
		}
	}
	return false, errors.New("parameter '" + key + "' with value '" + value + "' is not an boolean value")
}


type DatefmtRule struct {}

func (d *DatefmtRule) Validate(key, tagValue, value string) (bool, error) {
	_, err := time.Parse(tagValue, value)
	if err != nil {
		return false, fmt.Errorf("parameter %s with value %s is not a '%s' date", key, value, tagValue)
	}
	return true, nil
}


// BetweenRule "between:2,10" : 2 <= n <= 10
type BetweenRule struct {}

func (*BetweenRule) Validate(key, tagValue, value string) (bool, error) {
	v, err := strconv.Atoi(value)
	if err != nil {
		return false, errors.New("parameter '"+ key +"' is not an integer")
	}

	tagValues := strings.Split(tagValue, ",") // tagValue fmt: 2,10
	if len(tagValues) == 2 {
		min, err := strconv.Atoi(tagValues[0])
		if err != nil {
			return false, errors.New("rule with wrong format: "+ tagValue)
		}
		max, err := strconv.Atoi(tagValues[1])
		if err != nil {
			return false, errors.New("rule with wrong format: "+ tagValue)
		}
		if v >= min && v <= max {
			return true, nil
		} else {
			return false, errors.New("parameter '" + key + "' is not between "+ tagValue)
		}
	}
	return false, errors.New("rule with wrong format: "+ tagValue)
}

// EnumRule "enum:2,34,7,d,sfs"
type EnumRule struct {}

func (*EnumRule) Validate(key, tagValue, value string) (bool, error) {
	tagValues := strings.Split(tagValue, ",") // tagValue fmt: 2,10,ds
	for _, sub := range tagValues {
		if value == sub {
			return true, nil
		}
	}
	return false, errors.New("parameter " + key + "'s value "+ value + " is not in "+tagValue)
}






