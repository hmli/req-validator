# Req-Validator

对 HTTP 请求里传入的参数 (包括 Query param 和 form value) 进行验证, 实现参考了 laravel 框架. [Laravel validation](https://laravel.com/docs/5.5/validation)



## 示例
```golang
// 一个 HTTP 请求, 有 a b c 三个参数
req, err  := http.NewRequest("GET", "http://localhost?a=1&b=2&c=ddd", nil)
if err != nil {
    // error handling
}
// 新建一个 validator, 并验证 a必须有值且为整数, b必须有值也为整数
v := New(map[string]string{ // true
    "a": "required|int",
    "b": "int",
})
validated, err := v.Validate(req)
// validated: 是否通过验证
// err 返回的错误
```


## Rules

在上面的示例中的 `required|int` 被称为 Rule.

目前支持以下 rule:

### notnull
这个值必须是非空

### required
这个值必须存在, 其实和 notnull 一样

### int
这个值是整数, 比如 "1" "123" "33333"

### bool
这个值是 bool, 即 0, 1 或者 true false 的任意大小写. 比如 "1", "true", "False", "trUE", "FALSE"

### datefmt
时间格式. eg: `datefmt:2006-01-02 15:03:05`

这个值必须是一个可被后面的时间格式成功解析的时间字符串.

### enum

枚举. eg: `enum:1,3,ddd,sfsf,235`

这个值必须是后面枚举的几个值其中之一, 各值以 `,`分隔

### between

数字范围 eg: `between:2,10`

这个值一定是一个整数, 且取值范围是 2<=n<=10

## 自定义 Rule

实现 `Rule` interface, 然后调用 `RegisterRule()` 进行注册
```golang

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

RegisterRule("int...", &IntRule{})
```
