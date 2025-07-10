package main

import "cart/w4"

var smiley = [8]byte{
	0b11000011,
	0b10000001,
	0b00100100,
	0b00100100,
	0b00000000,
	0b00100100,
	0b10011001,
	0b11000011,
}

type Player struct {
	x, y int
}

var player = Player{80, 130} 

func handleInput() {
	gamepad := *w4.GAMEPAD1

	if gamepad&w4.BUTTON_UP != 0 {
		player.y -= 2
	}
	if gamepad&w4.BUTTON_DOWN != 0 {
		player.y += 2
	}
	if gamepad&w4.BUTTON_LEFT != 0 {
		player.x -= 2
	}
	if gamepad&w4.BUTTON_RIGHT != 0 {
		player.x += 2
	}
}

func draw() {
	// Paleta: céu azul, jogador branco
	w4.PALETTE[0] = 0x89CFF0 // Azul claro (fundo)
	w4.PALETTE[1] = 0xFFFFFF // Branco (jogador)

	// Céu azul (fundo)
	*w4.DRAW_COLORS = 1
	w4.Rect(0, 0, 160, 160)

	// // Jogador
	*w4.DRAW_COLORS = 2
	w4.Rect(player.x, player.y, 8, 8)
}

//go:export update
func update() {
	handleInput()
	draw()
}
