package main

import (
	"errors"
	"sync"
)

var (
	ErrRunningGame = errors.New("there is a running game")
	ErrGameNull    = errors.New("game is null")
)

type Server struct {
	Games map[string]*Game

	mu sync.RWMutex
}

func NewServer() *Server {
	return &Server{
		Games: make(map[string]*Game),
	}
}

func (s *Server) AddGame(k string, g *Game) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if g == nil {
		return ErrGameNull
	}

	if existingGame, ok := s.Games[k]; ok {
		if !existingGame.Ended {
			return ErrRunningGame
		}
	}

	s.Games[k] = g
	return nil
}

func (s *Server) Game(k string) *Game {
	s.mu.Lock()
	defer s.mu.Unlock()

	if g, ok := s.Games[k]; ok {
		return g
	}
	return nil
}
