/*
 * @Author: hugo2lee
 * @Date: 2024-10-30 18:04
 * @LastEditors: hugo2lee
 * @LastEditTime: 2025-01-15 15:03
 * @FilePath: \gotox\requestx\request.go
 * @Description:
 *
 * Copyright (c) 2025 by hugo, All Rights Reserved.
 */
package requestx

import (
	"github.com/go-resty/resty/v2"
	jsoniter "github.com/json-iterator/go"
)

type Requestx struct {
	*resty.Client
}

func NewRequestx() *Requestx {
	r := resty.New()
	r.SetJSONMarshaler(jsoniter.Marshal)
	r.SetJSONUnmarshaler(jsoniter.Unmarshal)
	return &Requestx{
		Client: r,
	}
}

func (r *Requestx) Get(url string) (*resty.Response, error) {
	return r.Client.R().Get(url)
}

func (r *Requestx) Post(url string, body any) (*resty.Response, error) {
	return r.Client.R().
		SetBody(body).
		Post(url)
}
