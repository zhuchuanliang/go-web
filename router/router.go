package router

import (
	"git.woa.com/alanclzhu/go-web/core"
	"git.woa.com/alanclzhu/go-web/middlerware"
)

/**
 * @Description $
 * @Date 2023/2/22 13:03
 **/

func RegisterRouter(core *core.Core) {
	//core.Get("foo", FooControllerHandler)
	//静态路由
	core.Get("/user/login", middlerware.Test1(),UserLoginHandler)

	//批量通用前缀
	subjectApi := core.Group("/subject")
	{

		subjectApi.Get("/:id", SubjectGetHandler)
		subjectApi.Delete("/:id", SubjectDelHandler)
		subjectApi.Put("/:id", SubjectUpdateHandler)
		subjectApi.Get("/list/all", SubjectListHandler)
		subjectInnerApi := subjectApi.Group("/info")
		{
			subjectInnerApi.Get("/name", SubjectNameHandler)
		}
	}

}
