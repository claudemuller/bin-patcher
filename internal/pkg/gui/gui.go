package gui

import (
	"errors"
	"fmt"
	"os"

	g "github.com/AllenDang/giu"
	log "github.com/claudemuller/bin-patcher/internal/pkg/logger"
	"github.com/claudemuller/bin-patcher/internal/pkg/patcher"
	"github.com/sqweek/dialog"
)

type App struct {
	Logger  *log.Log
	inFile  string
	outFile string
	sig     string
	patch   string
}

func NewApp() App {
	return App{
		Logger: log.NewLogger(),
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

	a.Logger.Log("Loaded file: " + a.inFile)
}

func (a *App) onLoadOutFile() {
	var err error

	a.outFile, err = dialog.File().Save()
	if err != nil && !errors.Is(err, dialog.ErrCancelled) {
		g.Msgbox("Error opening file for saving", err.Error())
	}

	a.Logger.Log("Loaded file: " + a.outFile)
}

func (a *App) onPatch() {
	if err := patcher.Patch(a.inFile, a.outFile, a.sig, a.patch, a.Logger); err != nil {
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
		content := a.Logger.GetLogs()
		winSize, _ := w.GetSize()

		var infile string

		if len(a.inFile) > 0 {
			infile = fmt.Sprintf("%s...%s", a.inFile[:15], a.inFile[len(a.inFile)-15:])
		}

		var outFile string

		if len(a.outFile) > 0 {
			outFile = fmt.Sprintf("%s...%s", a.outFile[:15], a.outFile[len(a.outFile)-15:])
		}

		widgets := []g.Widget{
			g.Row(
				g.Label("Select input file:"),
				g.Label(infile),
				g.Button("Select").OnClick(a.onLoadInFile),
			),
			g.Row(
				g.Label("Select output file:"),
				g.Label(outFile),
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
