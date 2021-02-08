// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	http "net/http"

	mock "github.com/stretchr/testify/mock"
)

// Handler is an autogenerated mock type for the Handler type
type Handler struct {
	mock.Mock
}

// CancelGeneration provides a mock function with given fields: w, r
func (_m *Handler) CancelGeneration(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}

// GenerateFeed provides a mock function with given fields: w, r
func (_m *Handler) GenerateFeed(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}

// ListGenerationTypes provides a mock function with given fields: w, r
func (_m *Handler) ListGenerationTypes(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}

// ListGenerations provides a mock function with given fields: w, r
func (_m *Handler) ListGenerations(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}
