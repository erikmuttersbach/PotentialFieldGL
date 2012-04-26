package pf

import (
	"github.com/banthar/Go-SDL/sdl"
	"sdl/ttf"
	"github.com/banthar/gl"
)

type UI struct {
	running bool
	font *ttf.Font
	screen *sdl.Surface
}



func InitUI() *UI {
	ui := &UI{
		running: false,
	}

	if sdl.Init(sdl.INIT_EVERYTHING) != 0 {
		panic(sdl.GetError())
	}
	
	if ttf.Init() != 0 {
		panic(sdl.GetError())
	}

	ui.screen = sdl.SetVideoMode(300, 300, 32, sdl.OPENGL)
	
	if ui.screen == nil {
		panic("sdl error")
	}

	if gl.Init() != 0 {
		panic("gl error")	
	}

	gl.MatrixMode(gl.PROJECTION)

	gl.Viewport(0, 0, int(ui.screen.W), int(ui.screen.H))
	gl.LoadIdentity()
	gl.Ortho(0, float64(ui.screen.W), float64(ui.screen.H), 0, -1.0, 1.0)

	gl.ClearColor(1, 1, 1, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	
	// TTF
	sdl.EnableUNICODE(1)	
	ui.font = ttf.OpenFont("FontinSans.otf", 20)
	ui.font.SetStyle(ttf.STYLE_UNDERLINE)
	if ui.font == nil {
		panic(sdl.GetError())
	}

	ui.running = true
	
	return ui
}