package game

import (
	"math"
	"math/rand"

	"github.com/bunyk/fasolasi/src/ui"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type ParticleSystem struct {
	Sprites   []*pixel.Sprite
	Particles []*Particle
	Batch     *pixel.Batch
}

type Particle struct {
	Sprite  int
	Time    float64
	MaxTime float64

	Source      pixel.Vec
	Destination pixel.Vec
	Frequency   float64
	Amplitude   float64
	Phase       float64
}

func (p Particle) CurrentLocation() pixel.Vec {
	t := p.Time / p.MaxTime
	ampEnvelope := math.Sin(t*math.Pi) * p.Amplitude
	return pixel.Lerp(p.Source, p.Destination, t).Add(pixel.V(
		ampEnvelope*math.Cos(p.Phase+t*p.Frequency),
		ampEnvelope*math.Sin(p.Phase+t*p.Frequency),
	))
}

func NewParticleSystem(filename string, width, height float64) *ParticleSystem {
	spritesheet, err := ui.LoadPicture(filename)
	if err != nil {
		panic(err)
	}

	var sprites []*pixel.Sprite
	for x := spritesheet.Bounds().Min.X; x < spritesheet.Bounds().Max.X; x += width {
		for y := spritesheet.Bounds().Min.Y; y < spritesheet.Bounds().Max.Y; y += height {
			sprites = append(sprites,
				pixel.NewSprite(spritesheet, pixel.R(x, y, x+width, y+height)),
			)
		}
	}

	batch := pixel.NewBatch(&pixel.TrianglesData{}, spritesheet)
	return &ParticleSystem{
		Sprites: sprites,
		Batch:   batch,
	}
}

func (ps *ParticleSystem) Spawn(src, dst pixel.Vec) {
	ps.Particles = append(ps.Particles, &Particle{
		Sprite: rand.Intn(len(ps.Sprites)),

		MaxTime:     1.5,
		Source:      src,
		Destination: dst,

		Frequency: math.Pi * (rand.Float64()*2.0 - 1.0),
		Amplitude: rand.Float64() * src.To(dst).Len() * 0.25,
		Phase:     2 * math.Pi * rand.Float64(),
	})
}

func (ps *ParticleSystem) UpdateAndRender(win *pixelgl.Window, dt float64) {
	haveParticles := len(ps.Particles)
	for i := 0; i < haveParticles; i++ {
		ps.Particles[i].Time += dt
		if ps.Particles[i].Time > ps.Particles[i].MaxTime {
			ps.Particles[i] = ps.Particles[haveParticles-1]
			haveParticles--
			i--
		}
	}
	ps.Particles = ps.Particles[:haveParticles]

	ps.Batch.Clear()
	for _, p := range ps.Particles {
		ps.Sprites[p.Sprite].Draw(ps.Batch, pixel.IM.Moved(p.CurrentLocation()))
	}
	ps.Batch.Draw(win)
}
