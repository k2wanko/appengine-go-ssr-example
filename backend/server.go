package backend

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/robertkrimen/otto"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

var renderJS *otto.Script

func initRender() (err error) {
	b, err := ioutil.ReadFile("render.js")
	if err != nil {
		return
	}
	renderJS, err = otto.New().Compile("render.js", string(b))
	return
}

func init() {
	if err := initRender(); err != nil {
		panic(err)
	}
	http.HandleFunc("/", serverTop)
}

func resError(ctx context.Context, w http.ResponseWriter, code int, err error) {
	w.WriteHeader(code)
	log.Errorf(ctx, "%s", err)
	fmt.Fprintf(w, "Error: %v", err)
}

func newVM() (vm *otto.Otto, m map[string]otto.Value) {
	m = map[string]otto.Value{}

	vm = otto.New()
	vm.Set("module", m)

	return
}

func mustToValue(v interface{}) otto.Value {
	res, _ := otto.ToValue(v)
	return res
}

func serverTop(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	vm, m := newVM()
	if _, err := vm.Run(renderJS); err != nil {
		resError(ctx, w, 500, err)
		return
	}

	if e, ok := m["exports"]; ok {
		req, res := map[string]interface{}{
			"url": r.URL.String(),
		}, map[string]interface{}{
			"setHeader": func(call otto.FunctionCall) otto.Value {
				if len(call.ArgumentList) < 2 {
					return mustToValue(errors.New("Invalid arguments"))
				}
				key, val := call.Argument(0), call.Argument(1)
				w.Header().Add(key.String(), val.String())
				return otto.UndefinedValue()
			},
			"writeHead": func(call otto.FunctionCall) otto.Value {
				if len(call.ArgumentList) == 0 {
					return mustToValue(errors.New("Invalid arguments"))
				}
				code, err := call.Argument(0).ToInteger()
				if err != nil {
					return mustToValue(err)
				}
				w.WriteHeader(int(code))
				return otto.UndefinedValue()
			},
			"write": func(call otto.FunctionCall) otto.Value {
				if len(call.ArgumentList) == 0 {
					return mustToValue(errors.New("Invalid arguments"))
				}
				d := call.Argument(0)
				fmt.Fprintf(w, "%s", d.String())
				return otto.UndefinedValue()
			},
		}
		if _, err := e.Call(otto.Value{}, req, res); err != nil {
			resError(ctx, w, 500, err)
			return
		}
	}
}
