package kurin

import (
	"log"
	"os"
	"os/signal"
)

type (
	App struct {
		name            string
		adapters        []Adapter
		fallibleSystems []Fallible
		closableSystems []Closable
		fail            chan error
	}

	Fallible interface {
		NotifyFail(chan error)
	}

	Closable interface {
		Close()
	}

	Adapter interface {
		Closable
		Open()
	}
)

func NewApp(name string, adapters ...Adapter) *App {
	closableSystems := make([]Closable, len(adapters))
	for i, a := range adapters {
		closableSystems[i] = Closable(a)
	}

	return &App{
		name:            name,
		adapters:        adapters,
		closableSystems: closableSystems,
		fallibleSystems: []Fallible{},
	}
}

func (a *App) RegisterFallibleSystems(systems ...Fallible) {
	a.fallibleSystems = append(a.fallibleSystems, systems...)
}

func (a *App) RegisterClosableSystems(systems ...Closable) {
	a.closableSystems = append(a.closableSystems, systems...)
}

func (a *App) Run() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	log.Printf("Starting %s ...\n", a.name)
	for _, a := range a.adapters {
		go a.Open()
	}

	for _, system := range a.fallibleSystems {
		system.NotifyFail(a.fail)
	}

	<-stop

	log.Println("Shutting down server...")

	for _, c := range a.closableSystems {
		c.Close()
	}
}
