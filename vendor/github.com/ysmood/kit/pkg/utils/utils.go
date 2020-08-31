package utils

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"reflect"
	"sync"
	"text/template"

	"github.com/tidwall/gjson"
)

// JSONResult shortcut for gjson.Result
type JSONResult = *gjson.Result

// Nil used to create empty channel
type Nil struct{}

// Noop swallow all args and do nothing
func Noop(_ ...interface{}) {}

// ErrArg get the last arg as error
func ErrArg(args ...interface{}) error {
	return args[len(args)-1].(error)
}

// E if the last arg is error, panic it
func E(args ...interface{}) []interface{} {
	err, ok := args[len(args)-1].(error)
	if ok {
		panic(err)
	}
	return args
}

// E1 if the second arg is error panic it, or return the first arg
func E1(arg interface{}, err error) interface{} {
	if err != nil {
		panic(err)
	}
	return arg
}

// ErrInjector let you easily mock error for testing
type ErrInjector struct {
	fn func(error) error
}

// CountInject inject err after E is called with specified times
func (e *ErrInjector) CountInject(times int, err error) {
	count := 1
	e.Inject(func(origin error) error {
		if count == times {
			e.Inject(nil)
			return err
		}
		count++
		return origin
	})
}

// Inject the fn and enable the enjection, call it with nil to disable injection
func (e *ErrInjector) Inject(fn func(error) error) {
	e.fn = fn
}

// E inject error
func (e *ErrInjector) E(err error) error {
	if e.fn == nil {
		return err
	}

	return e.fn(err)
}

// MustToJSONBytes encode data to json bytes
func MustToJSONBytes(data interface{}) []byte {
	bytes, err := json.Marshal(data)
	E(err)
	return bytes
}

// MustToJSON encode data to json string
func MustToJSON(data interface{}) string {
	return string(MustToJSONBytes(data))
}

// JSON parse json for easily access the value from json path
func JSON(data interface{}) JSONResult {
	var res gjson.Result
	switch v := data.(type) {
	case string:
		res = gjson.Parse(v)
	case []byte:
		res = gjson.ParseBytes(v)
	}

	return &res
}

// All run all actions concurrently, returns the wait function for all actions.
func All(actions ...func()) func() {
	wg := &sync.WaitGroup{}

	wg.Add(len(actions))

	runner := func(action func()) {
		defer wg.Done()
		action()
	}

	for _, action := range actions {
		go runner(action)
	}

	return wg.Wait
}

// RandBytes generate random bytes with specified byte length
func RandBytes(len int) []byte {
	b := make([]byte, len)
	_, _ = rand.Read(b)
	return b
}

// RandString generate random string with specified string length
func RandString(len int) string {
	b := RandBytes(len)
	return hex.EncodeToString(b)
}

// Try try fn with recover, return the panic as value
func Try(fn func()) (err interface{}) {
	defer func() {
		err = recover()
	}()

	fn()

	return err
}

// S Template render, the params is key-value pairs
func S(tpl string, params ...interface{}) string {
	var out bytes.Buffer

	dict := map[string]interface{}{}
	fnDict := template.FuncMap{}

	l := len(params)
	for i := 0; i < l-1; i += 2 {
		k := params[i].(string)
		v := params[i+1]
		if reflect.TypeOf(v).Kind() == reflect.Func {
			fnDict[k] = v
		} else {
			dict[k] = v
		}
	}

	t := template.Must(template.New("").Funcs(fnDict).Parse(tpl))
	E(t.Execute(&out, dict))

	return out.String()
}
