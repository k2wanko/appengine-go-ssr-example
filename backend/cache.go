package backend

import (
	"github.com/k2wanko/go-vue/server/render"
	"golang.org/x/net/context"

	"google.golang.org/appengine/log"
	"google.golang.org/appengine/memcache"
)

type componentCahce struct {
	ctx context.Context
}

func NewComponentCache(c context.Context) render.ComponentCacher {
	return &componentCahce{ctx: c}
}

func (c *componentCahce) Get(key string) ([]byte, error) {
	item, err := memcache.Get(c.ctx, key)
	if err != nil {
		log.Warningf(c.ctx, "memcache error %v", err)
		return nil, nil
	}
	return item.Value, nil
}

func (c *componentCahce) Set(key string, v []byte) error {
	memcache.Set(c.ctx, &memcache.Item{
		Key:   key,
		Value: v,
	})
	return nil
}
