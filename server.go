package main

type Server struct {
	Games []*Game
}

func (s *Server) AddGame(g *Game) {
	if g != nil {
		s.Games = append(s.Games, g)
	}
}
