package main

import (
	"fmt"
	"math/big"
)

type monkey struct {
	items []*big.Int
	op    func(*big.Int) *big.Int
	tdiv  int64
	t1    int
	t2    int
	insp  int
}

func main() {
	// f, _ := os.ReadFile("input.txt")
	// lines := strings.Split(string(f), "\n")
	fmt.Printf("Day11: \n")

	monkeys := []monkey{
		monkey{
			items: []*big.Int{big.NewInt(79), big.NewInt(98)},
			op: func(old *big.Int) *big.Int {
				return old.Mul(old, big.NewInt(19))
			},
			tdiv: 23,
			t1:   2,
			t2:   3,
		},
		monkey{
			items: []*big.Int{big.NewInt(54), big.NewInt(65), big.NewInt(75), big.NewInt(74)},
			op: func(old *big.Int) *big.Int {
				res := old.Add(old, big.NewInt(6))
				return res
			},
			tdiv: 19,
			t1:   2,
			t2:   0,
		},
		monkey{
			items: []*big.Int{big.NewInt(79), big.NewInt(60), big.NewInt(97)},
			op: func(old *big.Int) *big.Int {
				res := old.Mul(old, old)
				return res
			},
			tdiv: 13,
			t1:   1,
			t2:   3,
		},
		monkey{
			items: []*big.Int{big.NewInt(74)},
			op: func(old *big.Int) *big.Int {
				res := old.Add(old, big.NewInt(3))
				return res
			},
			tdiv: 17,
			t1:   0,
			t2:   1,
		},
	}
	/*monkeys := []monkey{
		monkey{
			items: []big.Int{91, 66},
			op:    func(old big.Int) big.Int { return old * 13 },
			tdiv:  19,
			t1:    6,
			t2:    2,
		},
		monkey{
			items: []big.Int{78, 97, 59},
			op:    func(old big.Int) big.Int { return old + 7 },
			tdiv:  5,
			t1:    0,
			t2:    3,
		},
		monkey{
			items: []big.Int{57, 59, 97, 84, 72, 83, 56, 76},
			op:    func(old big.Int) big.Int { return old + 6 },
			tdiv:  11,
			t1:    5,
			t2:    7,
		},
		monkey{
			items: []big.Int{81, 78, 70, 58, 84},
			op:    func(old big.Int) big.Int { return old + 5 },
			tdiv:  17,
			t1:    6,
			t2:    0,
		},
		monkey{
			items: []big.Int{60},
			op:    func(old big.Int) big.Int { return old + 8 },
			tdiv:  7,
			t1:    1,
			t2:    3,
		},
		monkey{
			items: []big.Int{57, 69, 63, 75, 62, 77, 72},
			op:    func(old big.Int) big.Int { return old * 5 },
			tdiv:  13,
			t1:    7,
			t2:    4,
		},
		monkey{
			items: []big.Int{73, 66, 86, 79, 98, 87},
			op:    func(old big.Int) big.Int { return old * old },
			tdiv:  3,
			t1:    5,
			t2:    2,
		},
		monkey{
			items: []big.Int{95, 89, 63, 67},
			op:    func(old big.Int) big.Int { return old + 2 },
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
	zero := big.NewInt(0)
	mod := big.NewInt(0)

	for i := 0; i < 20; i++ {

		for j := 0; j < len(monkeys); j++ {
			m := &monkeys[j]
			for _, item := range m.items {
				m.insp += 1
				op := m.op(item)
				newLvl := op.Div(op, big.NewInt(3))
				dest := m.t2
				mod.Mod(newLvl, big.NewInt(m.tdiv))
				if mod.Cmp(zero) == 0 {
					dest = m.t1
				}
				monkeys[dest].items = append(monkeys[dest].items, newLvl)
			}
			m.items = []*big.Int{}
		}
	}

	for i := 0; i < len(monkeys); i++ {
		fmt.Printf("%+v\n", monkeys[i].insp)
	}
}
