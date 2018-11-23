package main

import "fmt"

func main() {
	rankInfos := []int{1, 2, 3, 4, 5}
	//在排行榜中排除自己
	for i, v := range rankInfos {
		if v == 5 {
			//排行榜中只有自己一个人
			if i == 0 {
				//自己刚好是排行榜第一个人，这边三个条件主要考虑的是slice index 越界的问题
				rankInfos = rankInfos[1:]
			} else if i == len(rankInfos)-1 {
				//自己刚好是排行榜最后一个人
				rankInfos = rankInfos[:i]
			} else {
				//自己在排行榜中间
				rankInfos = append(rankInfos[:i], rankInfos[i+1:]...)
			}
		}
	}
	fmt.Printf("%v", rankInfos)
}
