package interactor

type ExportedFeedInteractor feedInteractor

func (f *feedInteractor) CsvFetcher() CsvRepo {
	return f.csvFetcher
}

func (f *feedInteractor) FileRepo() FileRepo {
	return f.files
}

func (f *feedInteractor) GenerationRepo() GenerationRepo {
	return f.generations
}

func (f *feedInteractor) Presenter() Presenter {
	return f.presenter
}
