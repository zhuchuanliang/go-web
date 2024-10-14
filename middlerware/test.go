package middlerware

import (
	"fmt"
	ct "git.woa.com/alanclzhu/go-web/context"
)

/**
 * @Description $
 * @Date 2024/10/11 20:21
 **/
func Test1() ct.ControllerHandler {
    return func(ctx *ct.Context) error{
        fmt.Println("middleware test1 pre")
		ctx.Next()
		fmt.Println("middleware test1 post")
		return nil
    }
}