package main

import (
	"fmt"
)

func Delete[T any](src []T, index int) ([]T, T, error) {
	// 源数组长度
	length := len(src)
	// 校验所要删除的下表是否合法
	if index < 0 || index >= length {
		// 定义零值
		var zero T
		// 错误
		err := fmt.Errorf("index out of range, length : %d, index : %d", length, index)
		return nil, zero, err
	}
	delElem := src[index]
	// i = length - 1 时，是最后1个位置的索引，索引从0开始
	for i := index; i < length-1; i++ {
		// 后面元素往前面移动1个位置
		src[i] = src[i+1]
	}
	// 子切片，共享底层数组，没有开辟新内存，切片是左闭右开，[:length - 1]为从第0个位置到倒数第2个位置
	src = src[:length-1]
	// 缩容
	src = shrink(src)
	return src, delElem, nil
}

func shrink[T any](src []T) []T {
	capacity, length := cap(src), len(src)
	// 重新计算切片容量
	newCapacity := recalculateCapacity(capacity, length)
	if capacity == newCapacity {
		// 数组容量没有发生变化，直接返回原切片
		return src
	}
	newSrc := make([]T, 0, newCapacity)
	newSrc = append(newSrc, src...)
	return newSrc
}

func recalculateCapacity(capacity int, length int) int {
	if capacity <= 32 {
		// 如果容量小于等于32，则不进行缩容
		return capacity
	}
	if capacity > 1024 && (capacity/length) >= 2 {
		// 如果容量大于1024，并且容量是元素数量的2倍以上，进行缩容，缩容因子为0.625
		factor := 0.625
		return int(float32(capacity) * float32(factor))
	}
	if capacity <= 1024 && (capacity/length >= 4) {
		// 如果容量小于等于1024 并且 容量是元素数量的4倍及以上，则缩容一半
		return capacity / 2
	}
	return capacity
}

func main() {
	// 测试缩容
	length := 1026
	delCount := length/2 + 1

	src := make([]int, 0, length)
	for i := 0; i < length; i++ {
		src = append(src, i)
	}
	fmt.Println(src)
	/**
	为何这里要先声明 delElem 和 err？
	如果此处不声明这两个变量，for循环体中就需要使用声明变量的写法，此时：
	src, delElem, err := Delete(src, delIndex)
	:= 左边的src是 for 循环体的里的局部变量，发生了变量覆盖，和上面的 src 不是同一个东西，而Delete方法参数的 src 始终是上面的src！
	声明了这两个变量后，for 循环体中就变成了对变量的重新赋值！
	src, delElem, err = Delete(src, delIndex)
	*/
	var (
		delElem int
		err     error
		// 始终删除第1个元素
		delIndex = 0
	)
	for i := 0; i < delCount; i++ {
		fmt.Printf("del index : %d\n", delIndex)

		src, delElem, err = Delete(src, delIndex)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("del element : %d\n", delElem)
		fmt.Printf("after delete, src slice length : %d, capacity : %d\n", len(src), cap(src))
		fmt.Printf("after delete, src slice : %v\n", src)
	}
}
