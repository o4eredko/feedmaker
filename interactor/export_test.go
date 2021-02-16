package interactor

type ExportedFeedInteractor feedInteractor

func (i *feedInteractor) GenerationRepo() FeedRepo {
	return i.feeds
}

func (i *feedInteractor) Presenter() Presenter {
	return i.presenter
}
