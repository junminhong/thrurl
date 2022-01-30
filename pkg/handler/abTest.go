package handler

import (
	"math/rand"
	"time"
)

func ABTest(sourceUrl string, sourceUrlB string, percent int) string {
	mr1 := ModelRes{
		Id:      1,
		Name:    sourceUrl,
		Percent: 100 - percent,
	}
	mr2 := ModelRes{
		Id:      2,
		Name:    sourceUrlB,
		Percent: percent,
	}
	var mrs []ModelRes
	mrs = append(mrs, mr1, mr2)
	prize := make(map[int]string)
	pack := make(map[int]int)
	for _, v := range mrs {
		prize[v.Id] = v.Name
		pack[v.Id] = v.Percent
	}
	return RedPack(prize, pack)
}

type ModelRes struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Percent int    `json:"percent"`
}

func RedPack(prize map[int]string, pack map[int]int) string {
	randArr := make(map[int][2]int)
	var sum int
	for k, v := range pack {
		var randArr1 [2]int
		randArr1[0] = sum
		sum += v
		randArr1[1] = sum
		randArr[k] = randArr1
	}
	rand.Seed(time.Now().Unix())
	s := rand.Intn(sum)
	id := 0
	for m, n := range randArr {
		if s >= n[0] && s < n[1] {
			id = m
			break
		}
	}
	return prize[id]
}
