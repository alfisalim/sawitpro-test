package handler

import (
	"github.com/SawitProRecruitment/UserService/middlewares"
	"github.com/SawitProRecruitment/UserService/repository"
)

type Server struct {
	Repository repository.RepositoryInterface
	Validator  middlewares.CustomValidatorInterface
}

type NewServerOptions struct {
	Repository repository.RepositoryInterface
	Validator  middlewares.CustomValidatorInterface
}

func NewServer(opts NewServerOptions) *Server {
	return &Server{
		Repository: opts.Repository,
		Validator:  opts.Validator,
	}
}
