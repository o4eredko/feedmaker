package repository

import (
	"context"
	"strings"

	"go-feedmaker/entity"
	"go-feedmaker/interactor"
)

type feedRepo struct {
}

func (r *feedRepo) GetFactoryByGenerationType(generationType string) (interactor.FeedFactory, error) {
	switch {
	case strings.HasPrefix(generationType, "yandex-"):
		return NewYandexFactory(), nil
	default:
		return NewDefaultFactory(), nil
	}
}

func (r *feedRepo) StoreGeneration(ctx context.Context, generation *entity.Generation) error {
	panic("implement me")
}

func (r *feedRepo) ListGenerations(ctx context.Context) ([]*entity.Generation, error) {
	panic("implement me")
}

func (r *feedRepo) ListAllowedTypes(ctx context.Context) ([]string, error) {
	panic("implement me")
}

func (r *feedRepo) IsAllowedType(ctx context.Context, generationType string) (bool, error) {
	panic("implement me")
}
