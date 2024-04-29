/*
 * @Author: hugo
 * @Date: 2024-04-23 19:56
 * @LastEditors: hugo
 * @LastEditTime: 2024-04-29 14:44
 * @FilePath: \gotox\webx\response.go
 * @Description:
 *
 * Copyright (c) 2024 by hugo, All Rights Reserved.
 */

package webx

// Response 统一的响应格式
type Response struct {
	// Code 响应的业务错误码。0表示业务执行成功，非0表示业务执行失败。
	Code int `json:"code"`
	// Message 响应的参考消息。前端可使用msg来做提示
	Message string `json:"msg"`
	// Data 响应的具体数据
	Data any `json:"data"`
}

func ResponseSuccess(data any) Response {
	return Response{
		Code: 200,
		Data: data,
	}
}

type ErrMsg struct {
	Code    int
	Message string
}

func ResponseErr(err ErrMsg) Response {
	return Response{
		Code:    err.Code,
		Message: err.Message,
	}
}
