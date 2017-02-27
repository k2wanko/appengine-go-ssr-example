package backend

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/net/context"

	"github.com/mjibson/goon"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

type APIError struct {
	Code    int
	Message string
	Err     error
}

func (e *APIError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s\n%v", e.Message, e.Err)
	}

	if e.Code != 0 && e.Message == "" {
		e.Message = http.StatusText(e.Code)
	}

	return e.Message
}

var (
	NotFound            = &APIError{Code: http.StatusNotFound}
	MethodNotAllowed    = &APIError{Code: http.StatusMethodNotAllowed}
	InternalServerError = &APIError{Code: http.StatusInternalServerError}
)

type Todo struct {
	ID      int64  `json:"id" datastore:"-" goon:"id"`
	Title   string `json:"title"`
	Updated int64  `json:"updated"`
}

func handleAPI(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	var res interface{}
	var err error

	path := r.URL.Path[len("/api"):]
	if r.Method == "POST" {
		switch path {
		case "/addTodoItem":
			res, err = addTodoItem(ctx, r)
		case "/getTodoItems":
			res, err = getTodoItems(ctx, r)
		case "/removeTodoItem":
			res, err = removeTodoItem(ctx, r)
		default:
			err = NotFound
		}
	} else {
		err = MethodNotAllowed
	}

	w.Header().Add("content-type", "application/json")

	if err != nil {
		code := 500
		if apiErr, ok := err.(*APIError); ok {
			code = apiErr.Code
		}
		log.Errorf(ctx, "API Error: %v", err)
		w.WriteHeader(code)
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err.Error())))
		return
	}

	if res == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	json, err := json.Marshal(res)
	if err != nil {
		log.Errorf(ctx, "json marshal error: %v", err)
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err.Error())))
		return
	}

	w.WriteHeader(200)
	log.Infof(ctx, "PATH %s\ncode = %d\nbody =\n%s", path, 200, json)
	w.Write(json)
}

func addTodoItem(ctx context.Context, r *http.Request) (interface{}, error) {
	title := r.FormValue("title")
	if title == "" || len(title) > 32 {
		return nil, &APIError{Code: http.StatusBadRequest, Message: "invalid title"}
	}

	todo := &Todo{
		Title:   title,
		Updated: time.Now().Unix(),
	}

	key, err := goon.FromContext(ctx).Put(todo)
	if err != nil {
		return nil, err
	}

	todo.ID = key.IntID()

	return todo, nil
}

func getTodoItems(ctx context.Context, r *http.Request) (res interface{}, err error) {
	q := datastore.NewQuery("Todo").Order("Updated")

	g := goon.FromContext(ctx)
	var todos []*Todo
	_, err = g.GetAll(q, &todos)
	if err != nil {
		return
	}

	return todos, nil
}

func removeTodoItem(ctx context.Context, r *http.Request) (res interface{}, err error) {
	id, err := strconv.ParseInt(r.FormValue("id"), 10, 64)
	if err != nil {
		msg := fmt.Sprintf("Invalid id: %s", r.FormValue("id"))
		log.Warningf(ctx, msg)
		return nil, &APIError{
			Code:    http.StatusBadRequest,
			Message: msg,
		}
	}

	log.Infof(ctx, "remove request: %d", id)

	g := goon.FromContext(ctx)
	todo := &Todo{ID: id}
	err = g.Get(todo)
	if err == datastore.ErrNoSuchEntity {
		return nil, NotFound
	} else if err != nil {
		return
	}

	err = g.Delete(datastore.NewKey(ctx, "Todo", "", todo.ID, nil))
	if err != nil {
		return
	}

	return todo, nil
}
