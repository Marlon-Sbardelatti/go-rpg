package main

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Sprite struct {
	Image *ebiten.Image
	X     float64
	Y     float64
}

type Player struct {
	*Sprite
	Health uint
}

type Enemy struct {
	*Sprite
	FollowsPlayer bool
}

type Potion struct {
	*Sprite
	AmtHeal uint
}

type Game struct {
	player  Player
	enemies []*Enemy
	potions []*Potion
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
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{120, 180, 255, 255})

	// draw player
	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Translate(g.player.X, g.player.Y)

	screen.DrawImage(g.player.Image.SubImage(
		image.Rect(0, 0, 16, 16),
	).(*ebiten.Image), &opts)

	for _, enemy := range g.enemies {
		opts.GeoM.Reset()
		opts.GeoM.Translate(enemy.X, enemy.Y)
		screen.DrawImage(enemy.Image.SubImage(image.Rect(0, 0, 16, 16)).(*ebiten.Image), &opts)
	}

	for _, potion := range g.potions {
		opts.GeoM.Reset()
		opts.GeoM.Translate(potion.X, potion.Y)
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

	game := Game{
		player: Player{
			&Sprite{
				Image: playerImage,
				X:     0,
				Y:     100,
			},
			100,
		},
		enemies: []*Enemy{
			{
				&Sprite{
					Image: skeletonImage,
					X:     50,
					Y:     50,
				},
				true,
			},
			{
				&Sprite{
					Image: skeletonImage,
					X:     75,
					Y:     75,
				},
				false,
			},
			{
				&Sprite{
					Image: skeletonImage,
					X:     150,
					Y:     150,
				},
				true,
			},
		},
		potions: []*Potion{
			{
				&Sprite{
					Image: potionImage,
					X:     80,
					Y:     50,
				},
				1.0,
			},
			{
				&Sprite{
					Image: potionImage,
					X:     100,
					Y:     20,
				},
				1.0,
			},
		},
	}

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
