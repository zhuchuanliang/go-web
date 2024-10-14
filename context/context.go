package context

/**
 * @Description $
 * @Date 2023/2/21 14:20
 **/

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type ControllerHandler func(c *Context) error

type Context struct {
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	Ctx            context.Context
	Handler        []ControllerHandler
	index int // 当前请求调用到调用链的哪个节点
	// 是否超时标记位
	HasTimeOut bool
	// 写保护机制
	Writermux *sync.Mutex
}

func NewContext(r *http.Request, w http.ResponseWriter) *Context {
	return &Context{
		Request:        r,
		ResponseWriter: w,
		Ctx:            r.Context(),
		Writermux:      &sync.Mutex{},
		index:          -1,
	}
}

// #region base function

func (ctx *Context) WriterMux() *sync.Mutex {
	return ctx.Writermux
}

func (ctx *Context) GetRequest() *http.Request {
	return ctx.Request
}

func (ctx *Context) GetResponse() http.ResponseWriter {
	return ctx.ResponseWriter
}

func (ctx *Context) SetHasTimeout() {
	ctx.HasTimeOut = true
}

func (ctx *Context) HasTimeout() bool {
	return ctx.HasTimeOut
}

// #endregion

func (ctx *Context) BaseContext() context.Context {
	return ctx.Request.Context()
}
// 为context设置handlers
func (ctx *Context) SetHandlers(handlers []ControllerHandler) {
	ctx.Handler = handlers
}
// #region implement context.Context
func (ctx *Context) Deadline() (deadline time.Time, ok bool) {
	return ctx.BaseContext().Deadline()
}

func (ctx *Context) Done() <-chan struct{} {
	return ctx.BaseContext().Done()
}

func (ctx *Context) Err() error {
	return ctx.BaseContext().Err()
}

func (ctx *Context) Value(key interface{}) interface{} {
	return ctx.BaseContext().Value(key)
}

// #endregion

// #region query url
func (ctx *Context) QueryInt(key string, def int) int {
	params := ctx.QueryAll()
	if vals, ok := params[key]; ok {
		len := len(vals)
		if len > 0 {
			intval, err := strconv.Atoi(vals[len-1])
			if err != nil {
				return def
			}
			return intval
		}
	}
	return def
}

func (ctx *Context) QueryString(key string, def string) string {
	params := ctx.QueryAll()
	if vals, ok := params[key]; ok {
		len := len(vals)
		if len > 0 {
			return vals[len-1]
		}
	}
	return def
}

func (ctx *Context) QueryArray(key string, def []string) []string {
	params := ctx.QueryAll()
	if vals, ok := params[key]; ok {
		return vals
	}
	return def
}

func (ctx *Context) QueryAll() map[string][]string {
	if ctx.Request != nil {
		return map[string][]string(ctx.Request.URL.Query())
	}
	return map[string][]string{}
}

// #endregion

// #region form post
func (ctx *Context) FormInt(key string, def int) int {
	params := ctx.FormAll()
	if vals, ok := params[key]; ok {
		len := len(vals)
		if len > 0 {
			intval, err := strconv.Atoi(vals[len-1])
			if err != nil {
				return def
			}
			return intval
		}
	}
	return def
}

func (ctx *Context) FormString(key string, def string) string {
	params := ctx.FormAll()
	if vals, ok := params[key]; ok {
		len := len(vals)
		if len > 0 {
			return vals[len-1]
		}
	}
	return def
}

func (ctx *Context) FormArray(key string, def []string) []string {
	params := ctx.FormAll()
	if vals, ok := params[key]; ok {
		return vals
	}
	return def
}

func (ctx *Context) FormAll() map[string][]string {
	if ctx.Request != nil {
		return map[string][]string(ctx.Request.PostForm)
	}
	return map[string][]string{}
}

// #endregion

// #region application/json post

func (ctx *Context) BindJson(obj interface{}) error {
	if ctx.Request != nil {
		body, err := ioutil.ReadAll(ctx.Request.Body)
		if err != nil {
			return err
		}
		ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		err = json.Unmarshal(body, obj)
		if err != nil {
			return err
		}
	} else {
		return errors.New("ctx.request empty")
	}
	return nil
}

// #endregion

// #region response

func (ctx *Context) Json(status int, obj interface{}) error {
	if ctx.HasTimeout() {
		return nil
	}
	ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
	ctx.ResponseWriter.WriteHeader(status)
	byt, err := json.Marshal(obj)
	if err != nil {
		ctx.ResponseWriter.WriteHeader(500)
		return err
	}
	ctx.ResponseWriter.Write(byt)
	return nil
}

func (ctx *Context) HTML(status int, obj interface{}, template string) error {
	return nil
}

func (ctx *Context) Text(status int, obj string) error {
	return nil
}

func (ctx *Context) Next() error {
    ctx.index ++
	if ctx.index < len(ctx.Handler) {
		if err:=ctx.Handler[ctx.index](ctx);err!=nil{
            return err
        }

    }
	return nil
}