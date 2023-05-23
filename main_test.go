/*
 * @Author: shgopher shgopher@gmail.com
 * @Date: 2023-05-23 22:16:07
 * @LastEditors: shgopher shgopher@gmail.com
 * @LastEditTime: 2023-05-23 22:18:52
 * @FilePath: /downM3U8/main_test.go
 * @Description:
 *
 * Copyright (c) 2023 by shgopher, All Rights Reserved.
 */
package main

import (
	"math/rand"
	"reflect"
	"strconv"
	"testing"
)

func TestSortString(t *testing.T) {
	strs := []string{"a1", "b12", "c2", "d3", "e10"}
	expected := []string{"a1", "c2", "d3", "e10", "b12"}

	// 调用待测试函数
	sortString(strs)

	// 检查排序后的结果是否符合预期
	if !reflect.DeepEqual(strs, expected) {
		t.Errorf("sortString(%v) = %v, want %v", strs, strs, expected)
	}
}

func BenchmarkSortString(b *testing.B) {
	// 生成一个包含 1000 个字符串的切片，用于测试
	strs := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		strs[i] = "a" + strconv.Itoa(rand.Intn(1000))
	}

	// 重置计时器
	b.ResetTimer()

	// 循环测试排序函数的性能
	for i := 0; i < b.N; i++ {
		sortString(strs)
	}
}
