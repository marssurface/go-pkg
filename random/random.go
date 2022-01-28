package random

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

/**
*  RandRange
*  @Description: 范围内随机数
*  @param min
*  @param max
*  @return int
 */
func RandRange(min, max int) int {
	return rand.Intn(max-min) + min
}
