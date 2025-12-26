package store

import (
	"context"
	"database/sql"
)

type actorStore struct {
	db *sql.DB
}

// CREATE
func (s *actorStore) Create(ctx context.Context, actor *Actor) error {
	query := `
		INSERT INTO actor (name)
		VALUES ($1)
		RETURNING id
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, actor.Name).Scan(&actor.ID)
	if err != nil {
		return err
	}

	return nil
}

// UPDATE
func (s *actorStore) Update(ctx context.Context, actor *Actor) error {
	query := `
		UPDATE actor
		SET name = $1
		WHERE id = $2
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, actor.Name, actor.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// GetList implements Actor.
func (s *actorStore) GetList(ctx context.Context) ([]Actor, error) {
	query := `
		SELECT id, name
		FROM actor
		ORDER BY id DESC
		LIMIT 20
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var actors []Actor
	for rows.Next() {
		var actor Actor
		if err := rows.Scan(&actor.ID, &actor.Name); err != nil {
			return nil, err
		}
		actors = append(actors, actor)
	}

	return actors, nil
}

func (s *actorStore) Delete(ctx context.Context, id int64) error {
	query := `
		DELETE FROM actor
		WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// alternative func
// func (s *actorStore) Updatee(ctx context.Context, actor Actor) error {
// 	query := `
// 		UPDATE actor
// 		SET name = $1
// 		WHERE id = $2
// 		RETURNING id

// 	`
// 	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
// 	defer cancel()

// 	err := s.db.QueryRowContext(
// 		ctx,
// 		query,
// 		actor.Name,
// 		actor.ID,
// 	).Scan(&actor.ID)
// 	if err != nil {
// 		switch {
// 		case errors.Is(err, sql.ErrNoRows):
// 			return ErrNotFound
// 		default:
// 			return err
// 		}
// 	}

// 	return nil
// }
