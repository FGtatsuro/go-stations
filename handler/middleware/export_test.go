package middleware

import "io"

type AccessLog = accessLog

func NewAccessLogMiddlewareWithWriter(w io.Writer) *accessLogMiddleware {
	return &accessLogMiddleware{
		w: w,
	}
}
