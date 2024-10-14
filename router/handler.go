package router

import (
	"context"
	"fmt"
	ct "git.woa.com/alanclzhu/go-web/context"
	"log"
	"time"
)

/**
 * @Description $
 * @Date 2024/9/12 20:22
 **/
func FooControllerHandler(c *ct.Context) error {
	finish := make(chan struct{}, 1)
	panicChan := make(chan interface{}, 1)

	durationCtx, cancel := context.WithTimeout(c.BaseContext(), time.Duration(30*time.Second))

	defer cancel()

	go func() {
		defer func() {
			if p := recover(); p != nil {
				panicChan <- p
			}
		}()
		time.Sleep(10 * time.Second)
		c.Json(200, "ok")
		finish <- struct{}{}
	}()
	select {
	case p := <-panicChan:
		c.WriterMux().Lock()
		defer c.WriterMux().Unlock()
		log.Println(p)
		c.Json(500, "panic")
	case <-finish:
		fmt.Println("finish")
	case <-durationCtx.Done():
		c.WriterMux().Lock()
		defer c.WriterMux().Unlock()
		c.Json(500, "time out")
		c.SetHasTimeout()
	}
	return nil
}

func UserLoginHandler(c *ct.Context) error {
	c.Json(200, "ok, UserLoginController")
	return nil
}

func SubjectListHandler(c *ct.Context) error {
	c.Json(200, "ok, SubjectListController")
	return nil
}
func SubjectGetHandler(c *ct.Context) error {
	c.Json(200, "ok, SubjectGetController")
	return nil
}
func SubjectDelHandler(c *ct.Context) error {
	c.Json(200, "ok, SubjectDelController")
	return nil
}
func SubjectUpdateHandler(c *ct.Context) error {
	c.Json(200, "ok, SubjectUpdateController")
	return nil
}
func SubjectNameHandler(c *ct.Context) error {
	c.Json(200, "ok, SubjectNameHandler")
	return nil
}
