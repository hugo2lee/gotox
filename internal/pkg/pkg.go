/*
 * @Author: hugo
 * @Date: 2024-06-18 15:01
 * @LastEditors: hugo2lee
 * @LastEditTime: 2025-03-17 14:59
 * @FilePath: \gotox\internal\pkg\pkg.go
 * @Description:
 *
 * Copyright (c) 2024 by hugo, All Rights Reserved.
 */
package pkg

import (
	"github.com/google/uuid"
)

func GenUuid() string {
	return uuid.New().String()
}

func IF[R any](condition bool, trueVal, falseVal R) R {
	if condition {
		return trueVal
	}
	return falseVal
}
