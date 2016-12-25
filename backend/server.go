package backend

import (
	"errors"
	"net/http"

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

	w.Header().Add("content-type", "text/html")

	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	} else {
		panic(errors.New("no flusher"))
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

	_, err = renderStream(nil, resObj)
	if jserr, ok := err.(*goja.Exception); ok {
		log.Criticalf(ctx, "%#v", jserr.Value().String())
	} else if err != nil {
		log.Criticalf(ctx, "%#v", err)
	}

	tr.Wait()
}
