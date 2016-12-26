package backend

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"

	"sync"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"github.com/k2wanko/appengine-ssr-example/timer"
)

func init() {
	http.Handle("/", http.HandlerFunc(serverRoot))
}

var requireMu sync.Mutex

func serverRoot(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	vm := goja.New()
	mr := new(require.Registry)
	rm := mr.Enable(vm)
	tr := timer.NewRegistry()
	tr.Enable(vm)

	vm.Set("debug", func(c goja.FunctionCall) goja.Value {
		log.Debugf(ctx, "%v\n", c.Arguments)
		return goja.Undefined()
	})

	process := vm.NewObject()
	env := vm.NewObject()
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		env.Set(pair[0], pair[1])
	}
	process.Set("browser", true)
	process.Set("env", env)
	setTimeout, ok := goja.AssertFunction(vm.Get("setTimeout"))
	if !ok {
		panic(errors.New("setTimeout is not fcuntion"))
	}
	process.Set("nextTick", func(c goja.FunctionCall) goja.Value {
		if len(c.Arguments) == 0 {
			return vm.NewGoError(errors.New("arguments is 0"))
		}
		setTimeout(nil, c.Argument(0))
		return goja.Undefined()
	})
	vm.Set("process", process)

	w.Header().Add("content-type", "text/html")

	// For standarone.
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}

	ComponentCahce(ctx, vm)
	v, err := rm.Require("./app/server-build.js")

	if jserr, ok := err.(*goja.Exception); ok {
		log.Criticalf(ctx, "%#v", jserr.Value().String())
	} else if err != nil {
		log.Criticalf(ctx, "%#v", err)
	}

	renderStream, ok := goja.AssertFunction(v.ToObject(vm).Get("renderStream"))
	if !ok {
		panic(errors.New("not function render"))
	}

	jsCtx := vm.NewObject()
	jsCtx.Set("url", r.URL.String())

	resObj := vm.NewObject()
	resObj.Set("write", func(c goja.FunctionCall) goja.Value {
		data := c.Argument(0).String()
		w.Write([]byte(data))
		return goja.Undefined()
	})
	resObj.Set("end", func(c goja.FunctionCall) goja.Value {
		data := c.Argument(0).String()
		w.Write([]byte(data))
		return goja.Undefined()
	})
	resObj.Set("error", func(c goja.FunctionCall) goja.Value {
		err := c.Argument(0)
		log.Errorf(ctx, "renderStream: %#v", err)
		return goja.Undefined()
	})

	jsCtx.Set("res", resObj)

	_, err = renderStream(nil, jsCtx)
	if jserr, ok := err.(*goja.Exception); ok {
		log.Criticalf(ctx, "%#v", jserr.Value().String())
	} else if err != nil {
		log.Criticalf(ctx, "%#v", err)
	}

	tr.Wait()
}
