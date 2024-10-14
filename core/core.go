package core

import (
	"git.woa.com/alanclzhu/go-web/context"
	"git.woa.com/alanclzhu/go-web/tree"
	"log"
	"net/http"
	"strings"
)

/**
 * @Description $
 * @Date 2023/2/21 22:21
 **/

type Core struct {
	router map[string]*tree.Tree
	moddlewares []context.ControllerHandler
}

func NewCore() *Core {
	//getRouter := map[string]context.ControllerHandler{}
	//postRouter := map[string]context.ControllerHandler{}
	//putRouter := map[string]context.ControllerHandler{}
	//deleteRouter := map[string]context.ControllerHandler{}
	router := map[string]*tree.Tree{}
	router["GET"] = tree.NewTree()
	router["POST"] = tree.NewTree()
	router["PUT"] = tree.NewTree()
	router["DELETE"] = tree.NewTree()

	return &Core{router: router}
}
func (core *Core) Group(prefix string) IGroup {
	return NewGroup(core, prefix)
}

func (core *Core) Get(url string, handler ...context.ControllerHandler) {
	//core.router[url] = handler
	//upperUrl := strings.ToUpper(url)
	//core.router["GET"][upperUrl] = handler
	allHandlers:=append(core.moddlewares, handler...)
	if err := core.router["GET"].AddRoute(url, allHandlers); err != nil {
		log.Println("add router error:", err)
	}
}
func (core *Core) Post(url string, handler ...context.ControllerHandler) {
	//upperUrl := strings.ToUpper(url)
	//core.router["POST"][upperUrl] = handler
	allHandlers:=append(core.moddlewares, handler...)
	if err := core.router["POST"].AddRoute(url, allHandlers); err != nil {
		log.Println("add router error:", err)
	}
}
func (core *Core) Put(url string, handler ...context.ControllerHandler) {
	//upperUrl := strings.ToUpper(url)
	//core.router["PUT"][upperUrl] = handler
	allHandlers:=append(core.moddlewares, handler...)
	if err := core.router["PUT"].AddRoute(url, allHandlers); err != nil {
		log.Println("add router error:", err)
	}
}

func (core *Core) Delete(url string, handler ...context.ControllerHandler) {
	//upperUrl := strings.ToUpper(url)
	//core.router["DELETE"][upperUrl] = handler
	allHandlers:=append(core.moddlewares, handler...)
	if err := core.router["DELETE"].AddRoute(url, allHandlers); err != nil {
		log.Println("add router error:", err)
	}
}
func(core *Core) Use(middleware ...context.ControllerHandler) {
	core.moddlewares = append(core.moddlewares, middleware...)
}
    //core.middleware = append(core.middleware, middleware...)
func (core *Core) FindRouterByRequest(request *http.Request) []context.ControllerHandler {
	upperMethod := strings.ToUpper(request.Method)
	upperUrl := strings.ToUpper(request.URL.Path)
	if router, ok := core.router[upperMethod]; ok {
		return router.FindHandler(upperUrl)
	}
	return nil
}

func (core *Core) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	log.Println("core ServeHTTP")
	ctx := context.NewContext(request, response)

	//router := core.router["foo"]
	//if router == nil {
	//	return
	//}
	//
	//log.Println("core router")
	//
	//router(ctx)
	//寻找路由
	router := core.FindRouterByRequest(request)
	if router == nil {
		// 如果没有找到，这里打印日志
		ctx.Json(404, "not found")
		return
	}
	ctx.SetHandlers(router)
	//if err := router(ctx); err != nil {
	//	ctx.Json(500, "inner error")
	//	return
	//}
	if err:=ctx.Next();err != nil {
		ctx.Json(500, "inner error")
		return
	}

}
