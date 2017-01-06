package backend

import (
	"io/ioutil"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"

	"github.com/k2wanko/go-vue/server/render"
)

func init() {
	http.Handle("/", http.HandlerFunc(serverRoot))
}

func serverRoot(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	m := make(map[string]interface{})
	m["url"] = r.URL.String()
	ren := &render.Renderer{
		Path:  "./app/server-build.js",
		Data:  m,
		Cache: NewComponentCache(ctx),
	}
	b, err := ioutil.ReadAll(ren)
	if err != nil {
		log.Criticalf(ctx, "render: %v", err)
		panic(err)
	}
	w.Write(b)
}
