package server

import (
	cHttp "docker/internal/containers/delivery/http"
	cRepository "docker/internal/containers/repository"
	cUseCase "docker/internal/containers/usecase"
)

func (s *Server) MapHandlers() error {
	cRepo := cRepository.NewContainersRepository(s.db)

	//TODO: middlewares
	containersUC := cUseCase.NewContainersUseCase(s.cfg, cRepo)

	containersHandlers := cHttp.NewContainerHandlers(s.cfg, containersUC)

	cHttp.MapContainersRoutes(s.mux, containersHandlers)
	return nil
}
