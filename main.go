package main

import "cart/w4"

var gameOver = false
var awaitReset = false
var frameCounter = 0

type Player struct {
	x, y int
}

type Obstacle struct {
	x, y, speed int
	width       uint
}

var player = Player{80, 130} 
var obstacles = []Obstacle{
	{0, 100, 2, 30},     // Carro na faixa 1
	{160, 80, -3, 40},   // Carro na faixa 2 (vai para a esquerda)
	{0, 60, 1, 20},      // Carro na faixa 3
}

var score = 0
var lastPlayerY = player.y

// Atualiza a pontuação 
func updateScore() {
	// Se jogador subiu (y diminuiu), incrementa a pontuação
	if player.y < lastPlayerY {
		score++
		lastPlayerY = player.y
	}

	// Evita pontuação negativa se descer
	if player.y > lastPlayerY {
		lastPlayerY = player.y
	}
}

// Converte inteiro positico em string
func intToStr(n int) string {
	if n == 0 {
		return "0"
	}
	digits := []byte{}
	for n > 0 {
		d := byte(n % 10)
		digits = append([]byte{ '0' + d }, digits...)
		n /= 10
	}
	return string(digits)
}

// Cria a iteração com os buttons de movimentação 
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

// Move os obstaculos na tela horizontalmente 
func moveObstacles() {
	for i := range obstacles {
		obstacles[i].x += obstacles[i].speed

		// Se sair da tela, volta para o lado oposto
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
}

// Verifica se houve colisão entre obstaculo e player 
func checkCollision() {
	for _, o := range obstacles {
		if rectsOverlap(player.x, player.y, 8, 8, o.x, o.y, int(o.width), 8) {
			gameOver = true

			// Som de colisão
			w4.Tone(220, 20, 100, w4.TONE_PULSE1|w4.TONE_MODE1)
			w4.Tone(150, 20, 150, w4.TONE_PULSE1|w4.TONE_MODE1)
		}
	}
}

// Cria o ambiente do jogo
func draw() {
	// Paleta: céu azul, jogador branco
	w4.PALETTE[0] = 0x29ADFF // Céu azul claro
	w4.PALETTE[1] = 0xFFFFFF // Jogador branco
	w4.PALETTE[2] = 0x808080 // Estrada cinza
	w4.PALETTE[3] = 0xFF0000 // Carro vermelho

	// Define quais cores usar: jogador = branco (1), fundo = azul (0)
	*w4.DRAW_COLORS = 0x21

	// Fundo azul
	w4.Rect(0, 0, 160, 160)

	// Estrada cinza 
	*w4.DRAW_COLORS = 3
	for _, o := range obstacles {
		w4.Rect(0, o.y, 160, 10)
	}

	// Carro vermelho
	*w4.DRAW_COLORS = 4
	for _, o := range obstacles {
		w4.Rect(o.x, o.y, o.width, 8)
	}

	// Jogador branco
	*w4.DRAW_COLORS = 0x21
	w4.Rect(player.x, player.y, 8, 8)

	w4.Text("Pontos: ", 5, 5)
	w4.Text(intToStr(score), 60, 5)
}

// Chama todas as funções
func update() {
	handleInput()
	draw()

	if gameOver {
		gameOverScreen()

		var gamepad = *w4.GAMEPAD1
		if gamepad&w4.BUTTON_1 != 0 {
			resetGame()
		}
	} else {
		moveObstacles()
		checkCollision()
		updateScore()
	}
}

func gameOverScreen() {
	*w4.DRAW_COLORS = 2 // cinza ou preto
	w4.Rect(20, 60, 120, 40)

	// Texto branco sobre a caixa
	*w4.DRAW_COLORS = 1
	w4.Text("COLIDIU", 40, 65)
	w4.Text("Pressione X", 44, 80)
	w4.Text("para reiniciar", 36, 90)
}
