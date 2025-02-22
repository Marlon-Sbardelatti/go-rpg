package main

import (
	"fmt"
	"game/entities"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	player      entities.Player
	enemies     []*entities.Enemy
	potions     []*entities.Potion
	tilemapJSON *TilemapJSON
	tilemapImg  *ebiten.Image
	cam         *Camera
}

func (g *Game) Update() error {
	// react to key presses
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.player.X += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.player.X -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.player.Y -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.player.Y += 2
	}

	for _, enemy := range g.enemies {
		if enemy.FollowsPlayer {
			if enemy.X < g.player.X {
				enemy.X += 1
			} else if enemy.X > g.player.X {
				enemy.X -= 1
			}

			if enemy.Y < g.player.Y {
				enemy.Y += 1
			} else if enemy.Y > g.player.Y {
				enemy.Y -= 1
			}
		}
	}

	for _, potion := range g.potions {
		if g.player.X > potion.X {
			g.player.Health += potion.AmtHeal
			fmt.Printf("Health: %v", g.player.Health)
		}
	}

	g.cam.FollowTarget(g.player.X, g.player.Y, 320, 240)
	g.cam.Constrain(
		float64(g.tilemapJSON.Layers[0].Width)*16,
		float64(g.tilemapJSON.Layers[0].Height)*16,
		320,
		240,
	)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{120, 180, 255, 255})
	opts := ebiten.DrawImageOptions{}

	for _, layer := range g.tilemapJSON.Layers {
		for index, id := range layer.Data {
			x := index % layer.Width
			y := index / layer.Width

			x *= 16
			y *= 16

			srcX := (id - 1) % 22
			srcY := (id - 1) / 22

			srcX *= 16
			srcY *= 16

			opts.GeoM.Translate(float64(x), float64(y))
			opts.GeoM.Translate(g.cam.X, g.cam.Y)

			screen.DrawImage(g.tilemapImg.SubImage(image.Rect(srcX, srcY, srcX+16, srcY+16)).(*ebiten.Image), &opts)
			opts.GeoM.Reset()
		}
	}

	// draw player
	opts.GeoM.Translate(g.player.X, g.player.Y)
	opts.GeoM.Translate(g.cam.X, g.cam.Y)

	screen.DrawImage(g.player.Image.SubImage(
		image.Rect(0, 0, 16, 16),
	).(*ebiten.Image), &opts)

	for _, enemy := range g.enemies {
		opts.GeoM.Reset()
		opts.GeoM.Translate(enemy.X, enemy.Y)
		opts.GeoM.Translate(g.cam.X, g.cam.Y)
		screen.DrawImage(enemy.Image.SubImage(image.Rect(0, 0, 16, 16)).(*ebiten.Image), &opts)
	}

	for _, potion := range g.potions {
		opts.GeoM.Reset()
		opts.GeoM.Translate(potion.X, potion.Y)
		opts.GeoM.Translate(g.cam.X, g.cam.Y)

		screen.DrawImage(potion.Image.SubImage(image.Rect(0, 0, 16, 16)).(*ebiten.Image), &opts)
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	playerImage, _, err := ebitenutil.NewImageFromFile("assets/images/ninja.png")
	if err != nil {
		log.Fatal(err)
	}

	skeletonImage, _, err := ebitenutil.NewImageFromFile("assets/images/skeleton.png")
	if err != nil {
		log.Fatal(err)
	}

	potionImage, _, err := ebitenutil.NewImageFromFile("assets/images/potion.png")
	if err != nil {
		log.Fatal(err)
	}

	tilemapJSON, err := NewTilemapJSON("assets/maps/spawn.json")
	if err != nil {
		log.Fatal(err)
	}

	tilemapImage, _, err := ebitenutil.NewImageFromFile("assets/images/TilesetFloor.png")
	if err != nil {
		log.Fatal(err)
	}

	game := Game{
		player: entities.Player{
			Sprite: &entities.Sprite{
				Image: playerImage,
				X:     0,
				Y:     100,
			},
			Health: 100,
		},
		enemies: []*entities.Enemy{
			{
				Sprite: &entities.Sprite{
					Image: skeletonImage,
					X:     50,
					Y:     50,
				},
				FollowsPlayer: true,
			},
			{
				Sprite: &entities.Sprite{
					Image: skeletonImage,
					X:     75,
					Y:     75,
				},
				FollowsPlayer: false,
			},
			{
				Sprite: &entities.Sprite{
					Image: skeletonImage,
					X:     150,
					Y:     150,
				},
				FollowsPlayer: true,
			},
		},
		potions: []*entities.Potion{
			{
				Sprite: &entities.Sprite{
					Image: potionImage,
					X:     80,
					Y:     50,
				},
				AmtHeal: 1.0,
			},
			{
				Sprite: &entities.Sprite{
					Image: potionImage,
					X:     100,
					Y:     20,
				},
				AmtHeal: 1.0,
			},
		},
		tilemapJSON: tilemapJSON,
		tilemapImg:  tilemapImage,
		cam:         NewCamera(0.0, 0.0),
	}

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
