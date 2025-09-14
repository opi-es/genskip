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
	VK_F10 = 0x79
	VK_F   = 0x46
)

func pressF() {
	// Key down (пробел = 0x20)
	keybdEvent.Call(uintptr(VK_F), 0, 0, 0)
	time.Sleep(10 * time.Millisecond)
	// Key up
	keybdEvent.Call(uintptr(VK_F), 0, 2, 0)
}

// isKeyPressed проверяет, нажата ли указанная клавиша
func isKeyPressed(keyCode int) bool {
	state, _, _ := getAsyncKeyState.Call(uintptr(keyCode))
	return state&0x8000 != 0
}

func main() {
	fmt.Println("Программа запущена. Нажмите F10 для паузы/продолжения.")

	paused := false
	keyPressed := false // Для отслеживания состояния клавиши

	// Даем время переключиться на окно игры
	time.Sleep(3 * time.Second)

	for {
		// Проверяем, нажата ли F10
		if isKeyPressed(VK_F10) {
			if !keyPressed {
				// Меняем состояние только если клавиша только что нажата
				paused = !paused
				if paused {
					fmt.Println("Программа приостановлена. Нажмите F10 для продолжения.")
				} else {
					fmt.Println("Программа продолжена.")
				}
				keyPressed = true
			}
		} else {
			keyPressed = false
		}

		// Если не на паузе, нажимаем пробел
		if !paused {
			pressF()
		}

		// Ждем немного перед следующей итерацией
		time.Sleep(100 * time.Millisecond)
	}
}
