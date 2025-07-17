package main

import "cart/w4"

type Player struct {
	x, y int
}

type Obstacle struct {
	x, y, speed int
	width       uint
}

type Coin struct {
	x, y      int
	collected bool
}

var player = Player{80, 130}
var obstacles = []Obstacle{
	{0, 100, 2, 30},
	{160, 80, -3, 40},
	{0, 60, 1, 20},
}

var lastPlayerY = player.y
var score = 0
var coins = [3]Coin{}
var roadY = [3]int{100, 80, 60}
var gameOver = false
var awaitReset = false
var showStartMenu = true
var frameCounter = 0
var rngSeed uint32 = 1
var blinkCounter = 0
var gameStarted = false
var menuFrameCounter = 0
var coinBlink = false
var playerBlink = false
var musicCounter = 0

func start() {}

func randInt(min, max int) int {
	rngSeed ^= rngSeed << 13
	rngSeed ^= rngSeed >> 17
	rngSeed ^= rngSeed << 5
	return min + int(rngSeed%(uint32(max-min+1)))
}

func drawChicken(x, y int) {
	// Corpo branco da galinha
	*w4.DRAW_COLORS = 2
	w4.Rect(x+1, y+2, 6, 4)

	// Cabeça branca
	*w4.DRAW_COLORS = 2
	w4.Rect(x+2, y, 4, 3)

	// Bico amarelo
	*w4.DRAW_COLORS = 0x40
	w4.Rect(x+1, y+1, 3, 2)

	// Crista vermelha
	*w4.DRAW_COLORS = 0x30
	w4.Rect(x+3, y-1, 2, 1)

	// Olho
	*w4.DRAW_COLORS = 0x04
	w4.Rect(x+2, y+1, 1, 1)

	// Pernas amarelas
	*w4.DRAW_COLORS = 0x40
	w4.Rect(x+2, y+6, 1, 2)
	w4.Rect(x+4, y+6, 1, 2)

	// Pés
	*w4.DRAW_COLORS = 0x40
	w4.Rect(x+1, y+7, 2, 1)
	w4.Rect(x+4, y+7, 2, 1)
}

func drawDetailedCar(x, y int, width uint, facingLeft bool) {
	w := int(width)

	if facingLeft {
		// Borda branca da cabine
		*w4.DRAW_COLORS = 2
		w4.Rect(x-1, y-1, uint(w/3)+2, 8)

		// Cabine verde
		*w4.DRAW_COLORS = 0x01
		w4.Rect(x, y, uint(w/3), 6)

		// Borda branca da carroceria principal
		*w4.DRAW_COLORS = 2
		w4.Rect(x+w/3-1, y+1, uint((w*2)/3)+2, 6)

		// Carroceria principal verde
		*w4.DRAW_COLORS = 0x01
		w4.Rect(x+w/3, y+2, uint((w*2)/3), 4)

		// Rodas
		*w4.DRAW_COLORS = 0x40
		w4.Oval(x+w-8, y+6, 4, 4)
		w4.Oval(x+4, y+6, 4, 4)
	} else {
		// Borda branca da carroceria principal
		*w4.DRAW_COLORS = 2
		w4.Rect(x-1, y+1, uint((w*2)/3)+2, 6)

		// Carroceria principal verde
		*w4.DRAW_COLORS = 0x01
		w4.Rect(x, y+2, uint((w*2)/3), 4)

		// Borda branca da cabine
		*w4.DRAW_COLORS = 2
		w4.Rect(x+(w*2)/3-1, y-1, uint(w/3)+2, 8)

		// Cabine verde
		*w4.DRAW_COLORS = 0x01
		w4.Rect(x+(w*2)/3, y, uint(w/3), 6)

		// Rodas
		*w4.DRAW_COLORS = 0x40
		w4.Oval(x+4, y+6, 4, 4)
		w4.Oval(x+w-8, y+6, 4, 4)
	}
}

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

	// Limites da tela
	if player.x < 0 {
		player.x = 0
	}
	if player.x > 160-8 {
		player.x = 160 - 8
	}
	if player.y < 0 {
		player.y = 0
	}
	if player.y > 160-8 {
		player.y = 160 - 8
	}
}

func moveObstacles() {
	for i := range obstacles {
		obstacles[i].x += obstacles[i].speed

		// Reposiciona quando sai da tela
		if obstacles[i].speed > 0 && obstacles[i].x > 160 {
			obstacles[i].x = -int(obstacles[i].width)
		} else if obstacles[i].speed < 0 && obstacles[i].x+int(obstacles[i].width) < 0 {
			obstacles[i].x = 160
		}
	}
}

func rectsOverlap(x1, y1, w1, h1, x2, y2, w2, h2 int) bool {
	return x1 < x2+w2 &&
		x1+w1 > x2 &&
		y1 < y2+h2 &&
		y1+h1 > y2
}

func resetGame() {
	player = Player{80, 130}
	obstacles[0].x = 0
	obstacles[1].x = 160
	obstacles[2].x = 0
	gameOver = false
	score = 0
	lastPlayerY = player.y
	musicCounter = 0

	generateCoins()
}

func playBackgroundMusic() {
	musicCounter++

	if musicCounter%30 == 0 {
		noteIndex := (musicCounter / 30) % 16

		// Melodia simples em loop
		melody := []uint{
			523, 659, 784, 659,
			523, 659, 784, 659,
			698, 784, 880, 784,
			659, 523, 440, 523,
		}

		w4.Tone(melody[noteIndex], 20, 20, w4.TONE_PULSE2|w4.TONE_MODE3)
	}
}

func playGameOverSound() {
	musicCounter = 0

	w4.Tone(523, 12, 100, w4.TONE_PULSE1|w4.TONE_MODE1)
	w4.Tone(415, 12, 100, w4.TONE_PULSE1|w4.TONE_MODE1)
	w4.Tone(330, 12, 100, w4.TONE_PULSE1|w4.TONE_MODE1)
	w4.Tone(220, 16, 100, w4.TONE_PULSE1|w4.TONE_MODE1)
	w4.Tone(247, 12, 90, w4.TONE_PULSE1|w4.TONE_MODE1)
	w4.Tone(175, 20, 100, w4.TONE_PULSE1|w4.TONE_MODE1)
	w4.Tone(131, 30, 100, w4.TONE_PULSE1|w4.TONE_MODE2)
}

func checkCollision() {
	for _, o := range obstacles {
		if rectsOverlap(player.x, player.y, 8, 8, o.x, o.y, int(o.width), 10) {
			gameOver = true
			playGameOverSound()
		}
	}
}

func playCoinSound() {
	w4.Tone(987, 5, 90, w4.TONE_PULSE1|w4.TONE_MODE1)
	w4.Tone(1319, 12, 100, w4.TONE_PULSE1|w4.TONE_MODE1)
	w4.Tone(1047, 8, 80, w4.TONE_PULSE1|w4.TONE_MODE1)
	w4.Tone(784, 4, 60, w4.TONE_PULSE1|w4.TONE_MODE2)
}

func checkCoinCollection() {
	for i := range coins {
		if !coins[i].collected && rectsOverlap(player.x, player.y, 8, 8, coins[i].x, coins[i].y, 6, 6) {
			coins[i].collected = true
			score++
			playCoinSound()
		}
	}

	if allCoinsCollected() {
		generateCoins()
	}
}

func generateCoins() {
	for i := range coins {
		coins[i].x = randInt(5, 155)
		coins[i].y = roadY[i%len(roadY)]
		coins[i].collected = false
	}
}

func allCoinsCollected() bool {
	for _, c := range coins {
		if !c.collected {
			return false
		}
	}
	return true
}

func drawDashedLine(y int) {
	const dashWidth = 4
	const dashSpacing = 4
	for x := 0; x < 160; x += dashWidth + dashSpacing {
		w4.Rect(x, y, dashWidth, 2)
	}
}

func drawRoad() {
	for _, o := range obstacles {
		*w4.DRAW_COLORS = 3
		w4.Rect(0, o.y, 160, 16)

		*w4.DRAW_COLORS = 2
		drawDashedLine(o.y + 8)
	}
}

func drawRoundedCoin(x, y int) {
	*w4.DRAW_COLORS = 0x40
	w4.Rect(x+1, y, 4, 6)
	w4.Rect(x, y+1, 6, 4)
	w4.Rect(x+1, y+1, 4, 4)
}

func drawScore() {
	*w4.DRAW_COLORS = 0x04
	w4.Text("Score: ", 5, 5)

	// Versão mais simples que funciona com números pequenos
	if score < 10 {
		switch score {
		case 0:
			w4.Text("0", 55, 5)
		case 1:
			w4.Text("1", 55, 5)
		case 2:
			w4.Text("2", 55, 5)
		case 3:
			w4.Text("3", 55, 5)
		case 4:
			w4.Text("4", 55, 5)
		case 5:
			w4.Text("5", 55, 5)
		case 6:
			w4.Text("6", 55, 5)
		case 7:
			w4.Text("7", 55, 5)
		case 8:
			w4.Text("8", 55, 5)
		case 9:
			w4.Text("9", 55, 5)
		}
	} else if score < 100 {
		tens := score / 10
		ones := score % 10

		// Desenha dezenas
		switch tens {
		case 1:
			w4.Text("1", 55, 5)
		case 2:
			w4.Text("2", 55, 5)
		case 3:
			w4.Text("3", 55, 5)
		case 4:
			w4.Text("4", 55, 5)
		case 5:
			w4.Text("5", 55, 5)
		case 6:
			w4.Text("6", 55, 5)
		case 7:
			w4.Text("7", 55, 5)
		case 8:
			w4.Text("8", 55, 5)
		case 9:
			w4.Text("9", 55, 5)
		}

		// Desenha unidades
		switch ones {
		case 0:
			w4.Text("0", 63, 5)
		case 1:
			w4.Text("1", 63, 5)
		case 2:
			w4.Text("2", 63, 5)
		case 3:
			w4.Text("3", 63, 5)
		case 4:
			w4.Text("4", 63, 5)
		case 5:
			w4.Text("5", 63, 5)
		case 6:
			w4.Text("6", 63, 5)
		case 7:
			w4.Text("7", 63, 5)
		case 8:
			w4.Text("8", 63, 5)
		case 9:
			w4.Text("9", 63, 5)
		}
	} else {
		w4.Text("99+", 55, 5)
	}
}

func drawMenu() {
	w4.PALETTE[0] = 0x00AA00
	w4.PALETTE[1] = 0xFFFFFF
	w4.PALETTE[2] = 0x808080
	w4.PALETTE[3] = 0xFFFF00

	// Animações do menu
	menuFrameCounter++
	if menuFrameCounter%15 == 0 {
		coinBlink = !coinBlink
	}
	if menuFrameCounter%20 == 0 {
		playerBlink = !playerBlink
	}
	moveObstacles()

	// Fundo verde
	*w4.DRAW_COLORS = 0x21
	w4.Rect(0, 0, 160, 160)

	drawRoad()

	for _, o := range obstacles {
		drawDetailedCar(o.x, o.y, o.width, o.speed < 0)
	}

	if coinBlink {
		for _, c := range coins {
			drawRoundedCoin(c.x, c.y)
		}
	}

	if playerBlink {
		drawChicken(player.x, player.y)
	}

	*w4.DRAW_COLORS = 0x04
	w4.Text("CROSSY GO", 45, 20)
	w4.Text("Colete moedas!", 28, 40)
	w4.Text("Pressione X", 42, 100)
	w4.Text("para comecar", 36, 110)
}

func draw() {
	w4.PALETTE[0] = 0x00AA00
	w4.PALETTE[1] = 0xFFFFFF
	w4.PALETTE[2] = 0x808080
	w4.PALETTE[3] = 0xFFFF00

	*w4.DRAW_COLORS = 0x21

	// Fundo verde
	w4.Rect(0, 0, 160, 160)

	drawRoad()

	for _, o := range obstacles {
		drawDetailedCar(o.x, o.y, o.width, o.speed < 0)
	}

	for _, c := range coins {
		if !c.collected {
			drawRoundedCoin(c.x, c.y)
		}
	}

	drawChicken(player.x, player.y)
	drawScore()
}

//go:export update
func update() {
	if !gameStarted {
		drawMenu()

		gamepad := *w4.GAMEPAD1
		if gamepad&w4.BUTTON_1 != 0 {
			gameStarted = true
			resetGame()
		}
		return
	}

	if !gameOver {
		handleInput()
		moveObstacles()
		checkCollision()
		checkCoinCollection()
		playBackgroundMusic()
	}

	draw()

	if gameOver {
		gameOverScreen()

		var gamepad = *w4.GAMEPAD1
		if gamepad&w4.BUTTON_1 != 0 {
			resetGame()
		}
	}
}

func gameOverScreen() {
	boxWidth := 120
	boxHeight := 70
	boxX := (160 - boxWidth) / 2
	boxY := (160 - boxHeight) / 2

	*w4.DRAW_COLORS = 2
	w4.Rect(boxX, boxY, uint(boxWidth), uint(boxHeight))

	textWidth := 7 * 8
	textX := (160 - textWidth) / 2
	*w4.DRAW_COLORS = 1
	w4.Text("COLIDIU", textX, boxY+8)

	*w4.DRAW_COLORS = 1
	w4.Text("Score: ", 50, boxY+22)

	// Versão mais simples que funciona com números pequenos
	if score < 10 {
		switch score {
		case 0:
			w4.Text("0", 106, boxY+22)
		case 1:
			w4.Text("1", 106, boxY+22)
		case 2:
			w4.Text("2", 106, boxY+22)
		case 3:
			w4.Text("3", 106, boxY+22)
		case 4:
			w4.Text("4", 106, boxY+22)
		case 5:
			w4.Text("5", 106, boxY+22)
		case 6:
			w4.Text("6", 106, boxY+22)
		case 7:
			w4.Text("7", 106, boxY+22)
		case 8:
			w4.Text("8", 106, boxY+22)
		case 9:
			w4.Text("9", 106, boxY+22)
		}
	} else if score < 100 {
		tens := score / 10
		ones := score % 10

		// Desenha dezenas
		switch tens {
		case 1:
			w4.Text("1", 106, boxY+22)
		case 2:
			w4.Text("2", 106, boxY+22)
		case 3:
			w4.Text("3", 106, boxY+22)
		case 4:
			w4.Text("4", 106, boxY+22)
		case 5:
			w4.Text("5", 106, boxY+22)
		case 6:
			w4.Text("6", 106, boxY+22)
		case 7:
			w4.Text("7", 106, boxY+22)
		case 8:
			w4.Text("8", 106, boxY+22)
		case 9:
			w4.Text("9", 106, boxY+22)
		}

		// Desenha unidades
		switch ones {
		case 0:
			w4.Text("0", 114, boxY+22)
		case 1:
			w4.Text("1", 114, boxY+22)
		case 2:
			w4.Text("2", 114, boxY+22)
		case 3:
			w4.Text("3", 114, boxY+22)
		case 4:
			w4.Text("4", 114, boxY+22)
		case 5:
			w4.Text("5", 114, boxY+22)
		case 6:
			w4.Text("6", 114, boxY+22)
		case 7:
			w4.Text("7", 114, boxY+22)
		case 8:
			w4.Text("8", 114, boxY+22)
		case 9:
			w4.Text("9", 114, boxY+22)
		}
	} else {
		w4.Text("99+", 106, boxY+22)
	}

	textWidth2 := 11 * 8
	textX2 := (160 - textWidth2) / 2
	w4.Text("Pressione X", textX2, boxY+36)

	textWidth3 := 14 * 8
	textX3 := (160 - textWidth3) / 2
	w4.Text("para reiniciar", textX3, boxY+50)
}
