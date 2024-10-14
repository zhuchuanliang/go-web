package core

import (
	ct "git.woa.com/alanclzhu/go-web/context"
)

/**
 * @Description $
 * @Date 2024/9/13 09:59
 **/
type IGroup interface {
	// 实现HttpMethod方法
	Get(string, ...ct.ControllerHandler)
	Post(string, ...ct.ControllerHandler)
	Put(string, ...ct.ControllerHandler)
	Delete(string, ...ct.ControllerHandler)

	//实现嵌套的group
	Group(string) IGroup
	// 嵌套中间件
	Use(middlewares ...ct.ControllerHandler)
	//Group(string)
}

type Group struct {
	core   *Core
	parent *Group //上一级前缀
	prefix string //这个group的通用前缀
	middlewares []ct.ControllerHandler
}

func NewGroup(core *Core, prefix string) *Group {

	return &Group{
		core:   core,
		parent: nil,
		prefix: prefix,
	}
}
func (g *Group) Use(middleware ...ct.ControllerHandler) {
	g.middlewares = append(g.middlewares, middleware...)
}

//	func (g *Group) Group(prefix string) {
//		g.core.Group(prefix)
//	}
func (g *Group) Get(uri string, handler ...ct.ControllerHandler) {
	uri = g.getAbsolutePrefix() + uri
	allHandlers := append(g.getMiddlewares(), handler...)
	g.core.Get(uri, allHandlers...)
}

// 实现Post方法
func (g *Group) Post(uri string, handler ...ct.ControllerHandler) {
	uri = g.getAbsolutePrefix() + uri
	allHandlers := append(g.getMiddlewares(), handler...)
	g.core.Post(uri, allHandlers...)
}

func (g *Group) Put(uri string, handler ...ct.ControllerHandler) {
	uri = g.getAbsolutePrefix() + uri
	allHandlers := append(g.getMiddlewares(), handler...)
	g.core.Put(uri, allHandlers...)
}

func (g *Group) Delete(uri string, handler ...ct.ControllerHandler) {
	uri = g.getAbsolutePrefix() + uri
	allHandlers := append(g.getMiddlewares(), handler...)
	g.core.Delete(uri, allHandlers...)
}

// 获取当前group的绝对路径
func (g *Group) getAbsolutePrefix() string {
	if g.parent == nil {
		return g.prefix
	}
	return g.parent.getAbsolutePrefix() + g.prefix
}

// 实现Group方法
func (g *Group) Group(uri string) IGroup {
	cgroup := NewGroup(g.core, uri)
	cgroup.parent = g
	return cgroup

}
func (g *Group) getMiddlewares() []ct.ControllerHandler {
    if g.parent == nil {
        return g.middlewares
    }
	return append(g.parent.getMiddlewares(), g.middlewares...)
}