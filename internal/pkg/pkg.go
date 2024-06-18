/*
 * @Author: hugo
 * @Date: 2024-06-18 15:01
 * @LastEditors: hugo
 * @LastEditTime: 2024-06-18 15:01
 * @FilePath: \gotox\internal\pkg\pkg.go
 * @Description:
 *
 * Copyright (c) 2024 by hugo, All Rights Reserved.
 */
package pkg

import "github.com/google/uuid"

func GenUuid() string {
	return uuid.New().String()
}
