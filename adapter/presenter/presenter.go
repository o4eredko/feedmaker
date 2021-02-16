package presenter

import "go-feedmaker/interactor"

type Presenter struct {
}

func (p *Presenter) PresentGenerationTypes(out []string) interface{} {
	return out
}

func (p *Presenter) PresentListGenerations(out *interactor.ListGenerationsOut) interface{} {
	return out
}

func (p *Presenter) PresentErr(err error) error {
	return err
}
