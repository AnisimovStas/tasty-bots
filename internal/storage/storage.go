package storage

import "tasty-bots/internal/tastybot"

type Storage struct {
	Bots []*tastybot.User
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Add(bot *tastybot.User) {
	s.Bots = append(s.Bots, bot)
}

func (s *Storage) GetByField(field string) (*tastybot.User, bool) {
	for _, bot := range s.Bots {
		if bot.Name == field {
			return bot, true
		}
	}

	return nil, false

}

func (s *Storage) GetAll() []*tastybot.User {
	return s.Bots
}
