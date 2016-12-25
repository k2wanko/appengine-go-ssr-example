package backend

import (
	"golang.org/x/net/context"

	"github.com/dop251/goja"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/memcache"
)

type componentCahce struct {
	ctx context.Context
	vm  *goja.Runtime
}

func (c *componentCahce) get(fc goja.FunctionCall) (v goja.Value) {
	key := fc.Argument(0).String()
	item, err := memcache.Get(c.ctx, key)
	if err != nil {
		log.Warningf(c.ctx, "memcache error %v", err)
		return
	}

	return c.vm.ToValue(string(item.Value))
}

func (c *componentCahce) set(fc goja.FunctionCall) (v goja.Value) {
	memcache.Set(c.ctx, &memcache.Item{
		Key:   fc.Argument(0).String(),
		Value: []byte(fc.Argument(1).String()),
	})
	return
}

func ComponentCahce(ctx context.Context, vm *goja.Runtime) {
	c := &componentCahce{ctx: ctx, vm: vm}
	o := vm.NewObject()
	o.Set("get", c.get)
	o.Set("set", c.set)
	vm.Set("ComponentCache", o)
}
