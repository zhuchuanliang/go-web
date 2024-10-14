package middlerware

import (
	ct "context"
	"fmt"
	"git.woa.com/alanclzhu/go-web/context"
	"log"
	"time"
)

/**
 * @Description $
 * @Date 2024/9/27 10:11
 **/
func TimeoutHandler(d time.Duration) context.ControllerHandler {
	return func(c *context.Context) error {
		finish := make(chan struct{}, 1)
		panicChan := make(chan interface{}, 1)

		//执行业务逻辑前预操作，初始化超时context
		durationCtx, cancel := ct.WithTimeout(c.BaseContext(), d)
		defer cancel()
		c.Request.WithContext(durationCtx)

		go func() {
			defer func() {
				if p:=recover(); p != nil {
                    panicChan <- p
                }
			}()

			c.Next()

			finish <- struct{}{}
		}()
		select {
        case <-finish:
			fmt.Println("finish")
            return nil
        case p := <-panicChan:
            log.Println(p)
			c.ResponseWriter.WriteHeader(500)
        case <-durationCtx.Done():
			c.SetHasTimeout()
            c.ResponseWriter.Write([]byte("time out"))
        }
		return nil
	}

}
