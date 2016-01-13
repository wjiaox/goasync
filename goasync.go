// gosync project async.go
package async

import (
	"errors"
	"reflect"
	"strings"
)

type AsyncItem struct {
	Handler reflect.Value
	Params  []reflect.Value
}

type Async struct {
	Count uint
	Items map[string]*AsyncItem
}

func NewAsync() *Async {
	return &Async{Count: 0, Items: make(map[string]*AsyncItem)}
}

func (a *Async) Add(name string, function interface{}, params ...interface{}) error {
	var handler reflect.Value
	var p4func []string

	if _, ok := a.Items[name]; ok {
		return errors.New("Task exist")
	}

	if reflect.TypeOf(function).Kind() == reflect.Func {
		handler = reflect.ValueOf(function)
		name := reflect.TypeOf(function).String()
		p := strings.Split(name, " ")
		if p[0] != "" {
			p4func = strings.Split(p[0][len("func("):len(p[0])-1], ",")
		}
	} else {
		return errors.New("need pass a function method")
	}
	if len(p4func) != len(params) {
		return errors.New("nums of params not match")
	}
	a.Items[name] = &AsyncItem{Params: make([]reflect.Value, len(p4func))}
	for i, v := range params {
		if p4func[i] != reflect.TypeOf(v).String() {
			return errors.New("function's params not match")
		} else {
			a.Items[name].Params[i] = reflect.ValueOf(v)
		}
	}
	a.Items[name].Handler = handler
	a.Count++
	return nil
}

func (a *Async) Go() (map[string][]interface{}, error) {
	cmap := make(chan map[string][]interface{}, a.Count)
	rmap := make(chan map[string][]interface{})
	if a.Count <= 0 {
		return nil, errors.New("no task")
	}

	go func() {
		last := make(map[string][]interface{})
		defer func(last map[string][]interface{}) {
			rmap <- last
		}(last)
		for {
			if a.Count <= 0 {
				break
			}
			select {
			case c := <-cmap:
				a.Count--
				for k, v := range c { //break out off mapping
					last[k] = v
				}
			}

		}
	}()
	for k, v := range a.Items {

		go func(key string, i *AsyncItem) {
			vals := i.Handler.Call(i.Params)
			l := len(vals)
			results := make([]interface{}, l)
			if l > 0 {
				for i, val := range vals {
					results[i] = val.Interface()
				}
				cmap <- map[string][]interface{}{key: results}

			}

		}(k, v)

	}
	result := <-rmap
	return result, nil

}

func init() {

}
