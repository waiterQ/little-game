package main

import (
	"fmt"
	"math/rand"
	"time"
)

// func init() {
// 	for {
// 		fmt.Printf("%s\r", time.Now().Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05"))
// 	}
// }

func main() {
	var row, col, lei int = 10, 10, 10
	fmt.Println("键入区域大小和雷数([row] [col] [num]):")
	fmt.Scanf("%d%d%d", &row, &col, &lei)
	if row <= 0 {
		row = 10
	}
	if col <= 0 {
		col = 10
	}
	if lei <= 0 {
		lei = 10
	}

	var wei, l2 int = 0, col

	for {
		if l2/10 == 0 {
			break
		}
		l2 /= 10
		wei += 1
	}

	var arr2 [][]Mine
	for i := 0; i < row; i++ {
		var list []Mine
		for i := 0; i < col; i++ {
			list = append(list, Mine{})
		}
		arr2 = append(arr2, list)
	}
	rand.Seed(time.Now().UnixNano())

	print(arr2, row, col, wei, 0)

	x, y := 0, 0
	// fmt.Println("键入坐标开始([row] [col]):")
	fmt.Scanf("%d%d", &x, &y)
	first_click := point{x, y}

	mime_num := 0
	var Mimes []point
	startTime := time.Now()
	defer func() {
		fmt.Printf("\n用时(s): %d", time.Now().Sub(startTime)/time.Second)
	}()

	for { // mime
		if mime_num > lei {
			break
		}
		x := rand.Intn(row)
		y := rand.Intn(col)
		if x >= first_click.X-1 && x <= first_click.X+1 && y >= first_click.Y-1 && y <= first_click.Y+1 { // first click area not allow mime
			continue
		}
		if arr2[x][y].Typ == 0 {
			arr2[x][y].Typ = -1
			mime_num += 1
			Mimes = append(Mimes, point{x, y})
		}
	}
	// fmt.Println("Mimes", Mimes)
	for i := 0; i < len(Mimes); i++ { // 数字
		for mx := Mimes[i].X - 1; mx <= Mimes[i].X+1; mx++ {
			if mx < 0 || mx > row-1 {
				continue
			}
			for my := Mimes[i].Y - 1; my <= Mimes[i].Y+1; my++ {
				if my < 0 || my > col-1 {
					continue
				}
				if arr2[mx][my].Typ != -1 {
					arr2[mx][my].Typ += 1
				}
			}
		}
	}

	click := make(chan point, 1)
	click <- first_click

	go func(click chan point) {
		for {
			x, y := 0, 0
			fmt.Scanf("%d%d", &x, &y)
			click <- point{X: x, Y: y}
		}
	}(click)

	var is_over bool
	var discovered int
	for {
		select {
		case p := <-click:
			m := &arr2[p.X][p.Y]
			if m.Checked {
			} else if m.Typ == -1 {
				m.Checked = true
				print2(arr2, row, col, wei, discovered)
				fmt.Println("\ngame over")
				is_over = true
				break
			} else if m.Typ == 0 {
				recDicover(arr2, p, &discovered, row, col)
			} else {
				m.Checked = true
				discovered += 1
			}
			// fmt.Println("discovered=", discovered)

			if discovered+len(Mimes) == row*col {
				print2(arr2, row, col, wei, discovered)
				fmt.Println("\nyou win")
				is_over = true
				break
			}

			print(arr2, row, col, wei, discovered)
		}
		if is_over {
			break
		}
	}
}

type point struct {
	X int
	Y int
}

type Mine struct {
	Checked bool
	Typ     int // 0 空白 1-8周围 -1boom
}

func print(arr2 [][]Mine, row, col, wei, discovered int) {
	view := ""
	if discovered > 0 {
		view = fmt.Sprintf("\033[%dA", 2*row+4)
	} else {
		view = "\n\n"
	}
	for i := 0; i < row; i++ {
		for j := 0; j < col; j++ {
			for n := 0; n < wei; n++ {
				view += " "
			}
			if !arr2[i][j].Checked {
				view += "?"
			} else if arr2[i][j].Typ == 0 {
				view += " "
			} else if arr2[i][j].Typ == -1 {
				view += "X"
			} else {
				view += fmt.Sprintf("%d", arr2[i][j].Typ)
			}
			view += " "
		}
		view += fmt.Sprintf("  %d\n\n", i)
	}
	view += "\n\n"
	for i := 0; i < col; i++ {
		is := fmt.Sprintf("%d", i)
		for n := 0; n < wei-len(is)+1; n++ {
			view += " "
		}

		view += fmt.Sprintf("%d ", i)
	}
	view += "\n"
	if discovered > 0 {
		view += "                                      \r键入一个坐标([x] [y]):"
	} else {
		view += "键入坐标开始([row] [col]):"
	}
	fmt.Printf("%s", view)
}

func print2(arr2 [][]Mine, row, col, wei, discovered int) {
	view := ""
	view = fmt.Sprintf("\033[%dA", 2*row+4)

	for i := 0; i < row; i++ {
		for j := 0; j < col; j++ {
			for n := 0; n < wei; n++ {
				view += " "
			}
			if arr2[i][j].Typ == 0 {
				view += " "
			} else if arr2[i][j].Typ == -1 {
				view += "X"
			} else {
				view += fmt.Sprintf("%d", arr2[i][j].Typ)
			}
			view += " "
		}
		view += fmt.Sprintf("  %d\n\n", i)
	}
	view += "\n\n"
	for i := 0; i < col; i++ {
		is := fmt.Sprintf("%d", i)
		for n := 0; n < wei-len(is)+1; n++ {
			view += " "
		}

		view += fmt.Sprintf("%d ", i)
	}
	view += "\n"
	fmt.Printf("%s", view)
}

func recDicover(arr2 [][]Mine, p point, discovered *int, row, col int) {
	for px := p.X - 1; px <= p.X+1; px++ {
		if px < 0 || px > row-1 {
			continue
		}
		for py := p.Y - 1; py <= p.Y+1; py++ {
			if py < 0 || py > col-1 {
				continue
			}
			if arr2[px][py].Checked {
			} else if arr2[px][py].Typ == 0 {
				arr2[px][py].Checked = true
				*discovered += 1
				recDicover(arr2, point{px, py}, discovered, row, col)
			} else if arr2[px][py].Typ > 0 {
				arr2[px][py].Checked = true
				*discovered += 1
			}
		}
	}
}
