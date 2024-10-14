package middlerware

import ct "git.woa.com/alanclzhu/go-web/context"

/**
 * @Description $
 * @Date 2024/10/14 22:12
 **/
func Recovery() ct.ControllerHandler {
	return func(c *ct.Context) error {
		defer func() {
			if err := recover(); err != nil {
				c.Json(500, err)
			}
		}()
		c.Next()
		return nil
	}
}
