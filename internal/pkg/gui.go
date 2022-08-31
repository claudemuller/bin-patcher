package pkg

import (
	"errors"
	"fmt"
	"os"

	g "github.com/AllenDang/giu"
	"github.com/sqweek/dialog"
)

type App struct {
	logger  *Log
	inFile  string
	outFile string
	sig     string
	patch   string
}

func NewApp() App {
	return App{
		logger: newLogger(),
	}
}

func (a *App) Run() {
	wnd := g.NewMasterWindow("Bin Patcher", 500, 250, 0)
	wnd.Run(a.loop(wnd))
}

func (a *App) onLoadInFile() {
	var err error

	a.inFile, err = dialog.File().Load()
	if err != nil && !errors.Is(err, dialog.ErrCancelled) {
		g.Msgbox("Error opening file", err.Error())
	}

	a.logger.log("Loaded file: " + a.inFile)
}

func (a *App) onLoadOutFile() {
	var err error

	a.outFile, err = dialog.File().Save()
	if err != nil && !errors.Is(err, dialog.ErrCancelled) {
		g.Msgbox("Error opening file for saving", err.Error())
	}

	a.logger.log("Loaded file: " + a.outFile)
}

func (a *App) onPatch() {
	if err := Patch(a.inFile, a.outFile, a.sig, a.patch, a.logger); err != nil {
		g.Msgbox("Error patching file", err.Error())
	}
}

func (a *App) onQuit() {
	os.Exit(0)
}

func (a *App) canPatch() bool {
	return a.inFile == "" && a.outFile == "" && a.sig == "" && a.patch == ""
}

func (a *App) loop(w *g.MasterWindow) func() {
	return func() {
		content := a.logger.getLogs()
		winSize, _ := w.GetSize()

		widgets := []g.Widget{
			g.Row(
				g.Label("Select input file:"),
				g.Label(fmt.Sprintf("%s...%s", a.inFile[:15], a.inFile[len(a.inFile)-15:])),
				g.Button("Select").OnClick(a.onLoadInFile),
			),
			g.Row(
				g.Label("Select output file:"),
				g.Label(fmt.Sprintf("%s...%s", a.outFile[:15], a.outFile[len(a.outFile)-15:])),
				g.Button("Select").OnClick(a.onLoadOutFile),
			),
			g.Row(
				g.Label("Signature to locate:"),
				g.InputText(&a.sig).Size(g.Auto),
			),
			g.Row(
				g.Label("Patch to apply:"),
				g.InputText(&a.patch).Size(g.Auto),
			),
			g.InputTextMultiline(&content).Flags(g.InputTextFlagsReadOnly).Size(g.Auto, float32(winSize)-385),
			g.Align(g.AlignRight).To(
				g.Row(
					g.Button("Patch Bin!").OnClick(a.onPatch).Disabled(a.canPatch()),
					g.Button("Quit").OnClick(a.onQuit),
					g.Dummy(5.0, 0.0),
				),
			),
			g.PrepareMsgbox(),
		}

		g.SingleWindow().Layout(widgets...)
	}
}
