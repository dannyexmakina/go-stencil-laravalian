package middleware

import (
	"log"
	"time"

	"github.com/labstack/echo/v4"
)

// RequestLogger はリクエストごとにメソッド・パス・処理時間を出力する
func RequestLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		err := next(c)
		log.Printf("%s %s (%s)", c.Request().Method, c.Path(), time.Since(start))
		return err
	}
}
