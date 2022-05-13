package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"time"

	term "github.com/nsf/termbox-go"
	"nspark.com/start/tool"
)

// var (
// 	SlaveDB = database.SlaveDB
// )

var MyMap map[int]string = map[int]string{
	0: "  ",
	1: "🟧",
	2: "😠",
	3: "🤢",
	4: "🏠",
	5: "😆",
}

var Load [15][15]int
var Idx [2]int
var monsterIdx [6][2]int
var stop bool = false
var gameLive = true

var point = 0

var success bool = false

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	Start()
}

func Init() {
	Idx = [2]int{0, 0}
	Load = [15][15]int{
		{2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 1, 0, 1, 1, 1, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 1, 1, 1},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0},
		{0, 0, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0},
		{0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0},
		{1, 1, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 1, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 4},
	}
	stop = false
	gameLive = true
	point = 0
	success = false
	for i := 0; i < len(monsterIdx); i++ {
		MakeMonster(i)
	}
}

func Start() {
	Init()
	Draw()
	go MvMonster()
	err := term.Init()
	tool.CkErr(err)
	defer term.Close()

	go func() {
		for {
			if stop {
				goto EXIST
			}
			getPoint()
		}
	EXIST:
	}()

	for gameLive {
		ev := term.PollEvent()
		Mv(ev.Ch)
		time.Sleep(200)
	}
}

func Mv(a rune) {
	if !stop {
		switch a {
		case 119:
			if Idx[0]-1 > -1 && Load[Idx[0]-1][Idx[1]] != 1 {
				Idx[0] = Idx[0] - 1
				CkDie(&Idx, true)
				Load[Idx[0]+1][Idx[1]] = 0
				Load[Idx[0]][Idx[1]] = 2
			}
		case 115:
			if Idx[0]+1 < len(Load) && Load[Idx[0]+1][Idx[1]] != 1 {
				Idx[0] = Idx[0] + 1
				CkDie(&Idx, true)
				Load[Idx[0]-1][Idx[1]] = 0
				Load[Idx[0]][Idx[1]] = 2
			}
		case 97:
			if Idx[1]-1 > -1 && Load[Idx[0]][Idx[1]-1] != 1 {
				Idx[1] = Idx[1] - 1
				CkDie(&Idx, true)
				Load[Idx[0]][Idx[1]+1] = 0
				Load[Idx[0]][Idx[1]] = 2
			}
		case 100:
			if Idx[1]+1 < len(Load[0]) && Load[Idx[0]][Idx[1]+1] != 1 {
				Idx[1] = Idx[1] + 1
				CkDie(&Idx, true)
				Load[Idx[0]][Idx[1]-1] = 0
				Load[Idx[0]][Idx[1]] = 2
			}
		case 113:
			gameLive = false
		}
		if stop && !success {
			Load[Idx[0]][Idx[1]] = 3
		}
		if success {
			Load[Idx[0]][Idx[1]] = 5
			stop = true
		}
		Draw()
	} else {
		switch a {
		case -1:
			if success {
				fmt.Println("success\n종료 : q , 재시작 : r")
			} else {
				fmt.Println("game over\n종료 : q , 재시작 : r")
			}
		case 113:
			gameLive = false
		case 114:
			Start()
		}

	}
}

func getPoint() {
	point++
	time.Sleep(time.Second)
}

func CkDie(location *[2]int, isUser bool) bool {
	switch isUser {
	case true:
		if Load[location[0]][location[1]] == 3 {
			stop = true
			Mv(-1)
		}
		if Load[location[0]][location[1]] == 4 {
			stop = true
			success = true
			Mv(-1)
		}
	case false:
		if Load[location[0]][location[1]] == 2 {
			stop = true
			Mv(-1)
		}
	}
	return stop
}

func Clear() {
	cmd := exec.Command("clear") //Linux example, its tested
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func Draw() {
	Clear()
	println("time : ", point)
	fmt.Println("🟫🟫🟫🟫🟫🟫🟫🟫🟫🟫🟫🟫🟫🟫🟫🟫🟫")
	for i := 0; i < len(Load); i++ {
		for j := 0; j < len(Load[i]); j++ {
			if j == 0 {
				fmt.Print("🟫")
			}
			fmt.Print(MyMap[Load[i][j]])
			if j == len(Load[i])-1 {
				fmt.Print("🟫")
			}
		}
		fmt.Println()
	}
	fmt.Println("🟫🟫🟫🟫🟫🟫🟫🟫🟫🟫🟫🟫🟫🟫🟫🟫🟫")
}

func MakeMonster(idx int) {
	for {
		location := [2]int{rand.Intn(len(Load)), rand.Intn(len(Load))}
		loadLocation := Load[location[0]][location[1]]
		if loadLocation != 1 &&
			loadLocation != 2 &&
			loadLocation != 3 {
			Load[location[0]][location[1]] = 3
			monsterIdx[idx] = location
			break
		}
	}
}

func MvMonster() {

	for {
		if stop {
			Mv(-1)
			break
		}
		for idx, Idx := range monsterIdx {
			randLocation := rand.Intn(4)
			switch randLocation {
			case 0:
				if Idx[0]-1 > -1 && Load[Idx[0]-1][Idx[1]] != 1 && Load[Idx[0]-1][Idx[1]] != 3 && Load[Idx[0]-1][Idx[1]] != 4 {
					Idx[0] = Idx[0] - 1
					CkDie(&Idx, false)
					monsterIdx[idx] = Idx
					Load[Idx[0]+1][Idx[1]] = 0
					Load[Idx[0]][Idx[1]] = 3
				}
			case 1:
				if Idx[0]+1 < len(Load) && Load[Idx[0]+1][Idx[1]] != 1 && Load[Idx[0]+1][Idx[1]] != 3 && Load[Idx[0]+1][Idx[1]] != 4 {
					Idx[0] = Idx[0] + 1
					CkDie(&Idx, false)
					monsterIdx[idx] = Idx
					Load[Idx[0]-1][Idx[1]] = 0
					Load[Idx[0]][Idx[1]] = 3
				}
			case 2:
				if Idx[1]-1 > -1 && Load[Idx[0]][Idx[1]-1] != 1 && Load[Idx[0]][Idx[1]-1] != 3 && Load[Idx[0]][Idx[1]-1] != 4 {
					Idx[1] = Idx[1] - 1
					CkDie(&Idx, false)
					monsterIdx[idx] = Idx
					Load[Idx[0]][Idx[1]+1] = 0
					Load[Idx[0]][Idx[1]] = 3
				}
			case 3:
				if Idx[1]+1 < len(Load[0]) && Load[Idx[0]][Idx[1]+1] != 1 && Load[Idx[0]][Idx[1]+1] != 3 && Load[Idx[0]][Idx[1]+1] != 4 {
					Idx[1] = Idx[1] + 1
					CkDie(&Idx, false)
					monsterIdx[idx] = Idx
					Load[Idx[0]][Idx[1]-1] = 0
					Load[Idx[0]][Idx[1]] = 3
				}
			}
		}
		Draw()
		time.Sleep(time.Second / 4)
	}
}
