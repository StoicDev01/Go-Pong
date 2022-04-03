package main

import (
	fmt "fmt"
	"time"

	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const WINDOW_WIDTH = 1280
const WINDOW_HEIGHT = 720
const MIN_BALL_RADIUS = 5
const MAX_BALL_RADIUS = 15

type Ball struct {
	position rl.Vector2
	radius   float32
	speedX   float32
	speedY   float32
	color    rl.Color
}

type Player struct {
	position rl.Vector2
	size     rl.Vector2
	speed    float32
	color    rl.Color
	score    int32
}

// Implement draw for ball
func (ball Ball) draw() {
	rl.DrawCircle(int32(ball.position.X), int32(ball.position.Y), ball.radius, ball.color)
}

func (ball *Ball) move(delta float32) {
	ball.position.X += ball.speedX * delta
	ball.position.Y += ball.speedY * delta

	var screen_height = rl.GetScreenHeight()

	if ball.position.Y+ball.radius > float32(screen_height) {
		ball.speedY *= -1
	}
	if ball.position.Y-ball.radius < 0 {
		ball.speedY *= -1
	}
}

func (ball *Ball) reset() {
	var screen_width = rl.GetScreenWidth()
	var screen_height = rl.GetScreenHeight()

	// Center ball
	ball.position.X = float32(screen_width) / 2
	ball.position.Y = float32(screen_height) / 2

	// Apply radius
	ball.radius = MAX_BALL_RADIUS

	// Apply random Speed
	rand.Seed(time.Now().UnixNano())
	ball.speedX = randFloat(-400.0, 400.0)
	ball.speedY = randFloat(-400.0, 400.0)

	var minSpeed float32 = 300.0
	var maxSpeed float32 = 400.0

	// Clamp
	if ball.speedX > 0 {
		ball.speedX = rl.Clamp(ball.speedX, minSpeed, maxSpeed)
	} else if ball.speedX < 0 {
		ball.speedX = rl.Clamp(ball.speedX, -minSpeed, -maxSpeed)
	}

	if ball.speedY > 0 {
		ball.speedY = rl.Clamp(ball.speedY, minSpeed, maxSpeed)
	} else if ball.speedY < 0 {
		ball.speedY = rl.Clamp(ball.speedY, -minSpeed, -maxSpeed)
	}

}

func randFloat(min, max float32) float32 {
	return min + rand.Float32()*(max-min)
}

func (player Player) getRect() rl.Rectangle {
	return rl.Rectangle{
		X: player.position.X, Y: player.position.Y,
		Width: player.size.X, Height: player.size.Y,
	}
}

func (player *Player) move_up(delta float32) {
	player.position.Y -= player.speed * delta

	if player.position.Y < 0 {
		player.position.Y = 0
	}
}
func (player *Player) move_down(delta float32) {
	player.position.Y += player.speed * delta

	var screen_height = rl.GetScreenHeight()
	if player.position.Y+player.size.Y > float32(screen_height) {
		player.position.Y = float32(screen_height - int(player.size.Y))
	}
}

// implement draw for player
func (player Player) draw() {
	rl.DrawRectangleRec(player.getRect(), player.color)
}

func main() {
	rl.InitWindow(WINDOW_WIDTH, WINDOW_HEIGHT, "Go Pong!")
	rl.SetTargetFPS(60)

	var screen_width = rl.GetScreenWidth()
	var screen_height = rl.GetScreenHeight()

	// Create objects
	var ball Ball
	ball.position = rl.Vector2{X: float32(screen_width) / 2, Y: float32(screen_height) / 2}
	ball.color = rl.White
	ball.radius = 15

	ball.reset()

	var player1 Player
	player1.size = rl.Vector2{X: 25, Y: 100}
	player1.position = rl.Vector2{X: 50, Y: float32(screen_height/2) - player1.size.Y/2}
	player1.speed = 800
	player1.color = rl.White

	var player2 Player
	player2.size = rl.Vector2{X: 25, Y: 100}
	player2.position = rl.Vector2{
		X: float32(screen_width) - player2.size.X - 50,
		Y: float32(screen_height/2) - player1.size.Y/2,
	}
	player2.speed = 800
	player2.color = rl.White

	for !rl.WindowShouldClose() {
		var delta = rl.GetFrameTime()

		// Debug
		if rl.IsKeyPressed(rl.KeySpace) {
			ball.reset()
		}

		// Movement
		ball.move(delta)

		// Score count
		if ball.position.X > float32(screen_width) {
			player1.score += 1
			ball.reset()
		}
		if ball.position.X < 0 {
			player2.score += 1
			ball.reset()
		}

		if rl.IsKeyDown(rl.KeyW) {
			player1.move_up(delta)
		}
		if rl.IsKeyDown(rl.KeyS) {
			player1.move_down(delta)
		}

		if rl.IsKeyDown(rl.KeyUp) {
			player2.move_up(delta)
		}
		if rl.IsKeyDown(rl.KeyDown) {
			player2.move_down(delta)
		}

		// Collision
		if rl.CheckCollisionCircleRec(ball.position, ball.radius, player2.getRect()) {
			if ball.speedX > 0 {
				ball.speedX *= -1.1
				ball.radius -= 1
				ball.speedY = ((ball.position.Y - (player2.position.Y + player2.size.Y/2)) / (player2.size.Y / 2)) * -ball.speedX

				if ball.radius < MIN_BALL_RADIUS {
					ball.radius = MIN_BALL_RADIUS
				}
			}
		}
		if rl.CheckCollisionCircleRec(ball.position, ball.radius, player1.getRect()) {
			if ball.speedX < 0 {
				ball.speedX *= -1.1
				ball.radius -= 1
				ball.speedY = ((ball.position.Y - (player1.position.Y + player1.size.Y/2)) / (player1.size.Y / 2)) * +ball.speedX

				if ball.radius < MIN_BALL_RADIUS {
					ball.radius = MIN_BALL_RADIUS
				}

			}
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		ball.draw()
		player1.draw()
		player2.draw()

		// Draw Score
		var score_text = fmt.Sprintf("%d	%d", player1.score, player2.score)
		var score_font_size int32 = 30
		var score_text_size = rl.MeasureText(score_text, score_font_size)
		rl.DrawText(score_text, int32(screen_width)/2-score_text_size/2, 50, score_font_size, rl.White)

		// Draw Center line
		rl.DrawLine(int32(screen_width)/2, 0, int32(screen_width)/2, int32(screen_height), rl.White)

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
