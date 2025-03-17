/*
 * @Author: hugo
 * @Date: 2024-06-18 15:01
 * @LastEditors: hugo2lee
 * @LastEditTime: 2025-03-17 15:11
 * @FilePath: \gotox\internal\pkg\pkg_test.go
 * @Description:
 *
 * Copyright (c) 2024 by hugo, All Rights Reserved.
 */

package pkg

import (
	"reflect"
	"testing"
)

func TestIF(t *testing.T) {
	type args[R any] struct {
		condition bool
		trueVal   R
		falseVal  R
	}
	tests := []struct {
		name string
		args args[any]
		want any
	}{
		// TODO: Add test cases.
		{
			"if ture",
			args[any]{
				condition: true,
				trueVal:   1,
				falseVal:  2,
			},
			1,
		},
		{
			"if false",
			args[any]{
				condition: false,
				trueVal:   1,
				falseVal:  2,
			},
			2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IF(tt.args.condition, tt.args.trueVal, tt.args.falseVal); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IF() = %v, want %v", got, tt.want)
			}
		})
	}
}
