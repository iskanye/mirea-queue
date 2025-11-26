package app

type App struct{}

func New() *App {
	return &App{}
}

func (App) MustRun() {
}

func (App) Stop() {
}
