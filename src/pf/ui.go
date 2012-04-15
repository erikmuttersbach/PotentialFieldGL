package pf

import (
	"github.com/banthar/Go-SDL/sdl"
	"github.com/banthar/gl"
)

type UI struct {
	running bool
}



func InitUI() *UI {
	ui := &UI{
		running: false,
	}

	sdl.Init(sdl.INIT_VIDEO)
	
	

	var screen = sdl.SetVideoMode(300, 300, 32, sdl.OPENGL)
	if screen == nil {
		panic("sdl error")
	}

	if gl.Init() != 0 {
		panic("gl error")	
	}

	gl.MatrixMode(gl.PROJECTION)

	gl.Viewport(0, 0, int(screen.W), int(screen.H))
	gl.LoadIdentity()
	gl.Ortho(0, float64(screen.W), float64(screen.H), 0, -1.0, 1.0)

	gl.ClearColor(1, 1, 1, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	ui.running = true
	
	return ui
}