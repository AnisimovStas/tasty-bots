package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"tasty-bots/internal/tastybot"
)

type Storage struct {
	db *sql.DB
}

func New(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("cant open db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("cant connect to db: %w", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Add(ctx context.Context, bot *tastybot.Bot) (int, error) {
	q := `INSERT INTO bots (tastyToken, baseUrl,status,casesCount ) VALUES (?,?,?,?)`

	res, err := s.db.ExecContext(ctx, q, bot.TastyToken, bot.BaseUrl, bot.Status, bot.CasesCount)
	if err != nil {
		return 0, fmt.Errorf("error during adding bot: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error during adding bot: %w", err)
	}

	return int(id), nil
}

func (s *Storage) PickById(ctx context.Context, id int) (*tastybot.Bot, error) {
	q := `SELECT * FROM bots WHERE id = ? LIMIT 1`

	var bot tastybot.Bot

	err := s.db.QueryRowContext(ctx, q, id).Scan(&bot.Id, &bot.TastyToken, &bot.BaseUrl, &bot.Status, &bot.CasesCount)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("error during select %w", err)
	}

	return &bot, nil
}

// PickAll pick all bots.
func (s *Storage) PickAll(ctx context.Context) ([]tastybot.Bot, error) {
	q := `SELECT * FROM bots`

	var bots []tastybot.Bot

	rows, err := s.db.QueryContext(ctx, q)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error during select %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var bot tastybot.Bot
		if err := rows.Scan(&bot.Id, &bot.TastyToken, &bot.BaseUrl, &bot.Status, &bot.CasesCount); err != nil {
			return nil, fmt.Errorf("error during select %w", err)
		}
		bots = append(bots, bot)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during select %w", err)
	}

	return bots, nil
}

// ChangeStatusById select bot by id and change status
func (s *Storage) ChangeStatusById(ctx context.Context, id int, status string) error {
	q := /*sql */ `UPDATE bots SET status = ? WHERE id = ?`
	r, err := s.db.ExecContext(ctx, q, status, id)
	if err != nil {
		return fmt.Errorf("error during changing status for bot id: %v,  %w", id, err)
	}
	count, _ := r.RowsAffected()

	if count == 0 {
		return fmt.Errorf("error during changing status for bot id: %v,  %w", id, "can't find bot with that id")
	}

	return nil
}

func (s *Storage) IncreaseCaseCountById(ctx context.Context, id int) error {
	q := /*sql */ `UPDATE bots SET casesCount = (SELECT casesCount FROM bots WHERE id = ?)+1 WHERE id = ?`
	r, err := s.db.ExecContext(ctx, q, id, id)
	if err != nil {
		return fmt.Errorf("error during changing status for bot id: %v,  %w", id, err)
	}
	count, _ := r.RowsAffected()

	if count == 0 {
		return fmt.Errorf("error during changing status for bot id: %v,  %w", id, "can't find bot with that id")
	}

	return nil
}

func (s *Storage) Init(ctx context.Context) error {
	q := "CREATE TABLE IF NOT EXISTS bots (id INTEGER PRIMARY KEY, tastyToken TEXT, baseUrl TEXT, status TEXT, casesCount INTEGER)"

	_, err := s.db.ExecContext(ctx, q)

	if err != nil {
		return fmt.Errorf("error during create table %w", err)
	}
	return nil
}
