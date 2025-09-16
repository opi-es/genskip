package main

import (
	"fmt"
	"syscall"
	"time"
)

var (
	user32           = syscall.MustLoadDLL("user32.dll")
	keybdEvent       = user32.MustFindProc("keybd_event")
	getAsyncKeyState = user32.MustFindProc("GetAsyncKeyState")
)

// Константы для кодов клавиш
const (
	VK_F9 = 0x78
	VK_F  = 0x46
)

func pressF() {
	// Key down (F = 0x46)
	keybdEvent.Call(uintptr(VK_F), 0, 0, 0)
	time.Sleep(10 * time.Millisecond)
	// Key up
	keybdEvent.Call(uintptr(VK_F), 0, 2, 0)
}

// isKeyPressed проверяет, нажата ли указанная клавиша
func isKeyPressed(keyCode int) bool {
	state, _, _ := getAsyncKeyState.Call(uintptr(keyCode))
	return uint16(state)&0x8000 != 0
}

func main() {
	fmt.Println("Программа запущена. Нажмите F9 для паузы/продолжения.")

	paused := false
	keyPressed := false // для отслеживания состояния клавиши

	// даем время переключиться на окно игры
	time.Sleep(5 * time.Second)
	// таймер для периода между нажатиями Ф
	tickerAction := time.NewTicker(10 * time.Second)
	defer tickerAction.Stop()
	// тикер для опроса клавиш
	tickerPoll := time.NewTicker(50 * time.Millisecond)
	defer tickerPoll.Stop()

	for {
		select {
		case <-tickerPoll.C:
			if isKeyPressed(VK_F9) {
				if !keyPressed {
					paused = !paused
					if paused {
						fmt.Println("Пауза. Нажмите F9 для продолжения.")
					} else {
						fmt.Println("Продолжение. Нажмите F9 для паузы.")
					}
					keyPressed = true
				}
			} else {
				keyPressed = false
			}
		case <-tickerAction.C:
			if !paused {
				pressF()
			}
		}
	}
}
