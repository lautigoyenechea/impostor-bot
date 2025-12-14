package main

import "sync"

type Server struct {
	Games map[string]*Game

	mu sync.RWMutex
}

func NewServer() *Server {
	return &Server{
		Games: make(map[string]*Game),
	}
}

func (s *Server) AddGame(k string, g *Game) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if g == nil {
		return
	}

	if _, ok := s.Games[k]; ok {
		return
	}

	s.Games[k] = g
}

func (s *Server) Game(k string) *Game {
	s.mu.Lock()
	defer s.mu.Unlock()

	if g, ok := s.Games[k]; ok {
		return g
	}
	return nil
}
