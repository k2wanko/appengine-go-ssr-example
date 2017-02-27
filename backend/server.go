package backend

import (
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"

	"github.com/k2wanko/go-vue/server/render"
)

func init() {
	http.HandleFunc("/api/", handleAPI)
	http.HandleFunc("/", serverRoot)
}

func serverRoot(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	m := make(map[string]interface{})
	m["code"] = 200
	m["url"] = r.URL.String()
	ren := &render.Renderer{
		Path:  "./app/server-build.js",
		Data:  m,
		Cache: NewComponentCache(ctx),
	}

	b, err := ren.Render()
	if renerr, ok := err.(*render.RenderError); ok {
		status := 500
		log.Errorf(ctx, "RenderError: %v", err)
		if code, ok := renerr.Get("code").(int64); ok {
			status = int(code)
		}
		w.WriteHeader(status)
		return
	} else if err != nil {
		log.Errorf(ctx, "render: %#v", err)
		panic(err)
	}

	if code, ok := m["code"].(int64); ok {
		w.WriteHeader(int(code))
	} else {
		w.WriteHeader(200)
	}
	w.Write(b)
}
