package service

import (
	"os"
	"task7/pkg/repository"
)

type Authorization interface {
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
	IsAdmin (tokenString string) (bool, string, int)
}

type Update interface {
	UpdateComicTable() (UpdateResponse, error)
	UpdateIndexTable() (string, error)
	signalHandler(c chan os.Signal, resultComicsSlice []repository.Comic, written *[]int)
}

type Search interface {
	SearchComics(searchQuery string) ([]string, string, int)
}

type Service struct {
	Authorization
	Update
	Search
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Update:        NewUpdateService(repos.Update),
		Search:        NewSearchService(repos.Search),
	}
}
