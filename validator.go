package req_validator

import (
	"net/http"
	"strings"
	"errors"
)

type Validator struct {
	rules map[string][]string // eg: "aid": ["required", "max:100", "integer", "max_len:20"]
}

func New(rawTags map[string]string) (v *Validator) {
	v = &Validator{
		rules: make(map[string][]string),
	}

	for tName, tValue := range rawTags {
		values := strings.Split(tValue, "|")
		v.rules[tName] = values
	}
	return v
}

// Validate 使用现有的 rules 对传入的 req 进行验证
func (v *Validator) Validate(req *http.Request) (validated bool, err error) {
	if req.Form == nil {
		req.ParseForm()
	}
	for ruleName, ruleTags := range v.rules { // 遍历每一条规则
		formValue := req.Form.Get(ruleName)
		for _, tag := range ruleTags {
			validated, err := MatchRule(ruleName, tag, formValue)
			if !validated || err != nil {
				return validated, err
			}
		}
	}
	return true, nil

}

// MatchRule  对传入的 tag 匹配一条rule, 并返回rule的校验结果
// 如果匹配不到rule, 认为校验通过
func MatchRule(key, tag, value string) (validated bool, err error) {
	fullTag := strings.SplitN(tag, ":", 2)
	switch len(fullTag) {
	case 0:
		return true, nil
	case 1:
		rule, exists := Rules[fullTag[0]]
		if !exists {
			return true, nil
		}
		return rule.Validate( key, "", value)
		// only match tag name
	case 2:
		rule, exists := Rules[fullTag[0]]
		if !exists {
			return true, nil
		}
		return rule.Validate(key, fullTag[1], value)
	default:
		return false, errors.New("wrong format tag: " + tag)
	}

}

