package main

import (
	"fmt"
	"math"
)

type monkey struct {
	items []int64
	op    func(int64) int64
	tdiv  int64
	t1    int64
	t2    int64
	insp  int64
}

func main() {
	// f, _ := os.ReadFile("input.txt")
	// lines := strings.Split(string(f), "\n")
	fmt.Printf("Day11: \n")

	monkeys := []monkey{
		monkey{
			items: []int64{79, 98},
			op: func(old int64) int64 {
				res := old * 19
				return res

				if res/19 != old {
					return math.MaxInt64
				}
				return res
			},
			tdiv: 23,
			t1:   2,
			t2:   3,
		},
		monkey{
			items: []int64{54, 65, 75, 74},
			op: func(old int64) int64 {
				res := old + 6
				return res
				if res-6 != old {
					return math.MaxInt64
				}
				return res
			},
			tdiv: 19,
			t1:   2,
			t2:   0,
		},
		monkey{
			items: []int64{79, 60, 97},
			op: func(old int64) int64 {
				res := old * old
				return res
				if res/old != old {
					return math.MaxInt64
				}
				return res
			},
			tdiv: 13,
			t1:   1,
			t2:   3,
		},
		monkey{
			items: []int64{74},
			op: func(old int64) int64 {
				res := old + 3
				return res
				if res-3 != old {
					return math.MaxInt64
				}
				return res

			},
			tdiv: 17,
			t1:   0,
			t2:   1,
		},
	}
	/*monkeys := []monkey{
		monkey{
			items: []int64{91, 66},
			op:    func(old int64) int64 { return old * 13 },
			tdiv:  19,
			t1:    6,
			t2:    2,
		},
		monkey{
			items: []int64{78, 97, 59},
			op:    func(old int64) int64 { return old + 7 },
			tdiv:  5,
			t1:    0,
			t2:    3,
		},
		monkey{
			items: []int64{57, 59, 97, 84, 72, 83, 56, 76},
			op:    func(old int64) int64 { return old + 6 },
			tdiv:  11,
			t1:    5,
			t2:    7,
		},
		monkey{
			items: []int64{81, 78, 70, 58, 84},
			op:    func(old int64) int64 { return old + 5 },
			tdiv:  17,
			t1:    6,
			t2:    0,
		},
		monkey{
			items: []int64{60},
			op:    func(old int64) int64 { return old + 8 },
			tdiv:  7,
			t1:    1,
			t2:    3,
		},
		monkey{
			items: []int64{57, 69, 63, 75, 62, 77, 72},
			op:    func(old int64) int64 { return old * 5 },
			tdiv:  13,
			t1:    7,
			t2:    4,
		},
		monkey{
			items: []int64{73, 66, 86, 79, 98, 87},
			op:    func(old int64) int64 { return old * old },
			tdiv:  3,
			t1:    5,
			t2:    2,
		},
		monkey{
			items: []int64{95, 89, 63, 67},
			op:    func(old int64) int64 { return old + 2 },
			tdiv:  2,
			t1:    1,
			t2:    4,
		},
	}*/
	/*
		 2  2 275 99
		 3  4 210 97
		10  3 288 8
		 7  6 277 103
	*/

	for i := 0; i < 1000; i++ {

		for j := 0; j < len(monkeys); j++ {
			m := &monkeys[j]
			for _, item := range m.items {
				m.insp += 1
				newLvl := m.op(item)
				if newLvl < 0 {
					fmt.Println("overflow")
					//newLvl = math.MaxInt64
					// newLvl = m.items[k]
					// fmt.Printf("%d %d %+v\n", newLvl, m.items[k], m)
					// return
					//newLvl = m.items[k]
					//fmt.Println("overflow")
					//newLvl = -newLvl
				}
				dest := m.t2
				if newLvl%m.tdiv == 0 {
					dest = m.t1
				}
				monkeys[dest].items = append(monkeys[dest].items, newLvl)
			}
			m.items = []int64{}
		}
	}

	for i := 0; i < len(monkeys); i++ {
		fmt.Printf("%+v\n", monkeys[i].insp)
	}
}
