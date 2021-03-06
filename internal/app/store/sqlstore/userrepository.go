package sqlstore

import (
	"database/sql"

	"github.com/DalerBakhriev/social_network/internal/app/model"
	"github.com/DalerBakhriev/social_network/internal/app/store"
)

// UserRepository ...
type UserRepository struct {
	store *Store
}

// Create ...
func (r *UserRepository) Create(u *model.User) error {

	if err := u.BeforeCreate(); err != nil {
		return err
	}

	_, err := r.store.db.Exec(
		`INSERT INTO users (email, name, surname, age, sex, interests, city, encrypted_password)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		u.Email,
		u.Name,
		u.Surname,
		u.Age,
		u.Sex,
		u.Interests,
		u.City,
		u.EncryptedPassword,
	)

	return err
}

// FindByEmail ...
func (r *UserRepository) FindByEmail(email string) (*model.User, error) {

	u := &model.User{}
	if err := r.store.db.QueryRow(
		`SELECT id,
				email,
				name,
				surname,
				age,
				sex,
				interests,
				city,
				encrypted_password
		 FROM users
		 WHERE email = ?`,
		email,
	).Scan(
		&u.ID,
		&u.Email,
		&u.Name,
		&u.Surname,
		&u.Age,
		&u.Sex,
		&u.Interests,
		&u.City,
		&u.EncryptedPassword,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return u, nil
}

// Find ...
func (r *UserRepository) Find(id int) (*model.User, error) {

	u := &model.User{}
	if err := r.store.db.QueryRow(
		`SELECT id,
				email,
				name,
				surname,
				age,
				sex,
				interests,
				city,
				encrypted_password
		 FROM users
		 WHERE id = ?`,
		id,
	).Scan(
		&u.ID,
		&u.Email,
		&u.Name,
		&u.Surname,
		&u.Age,
		&u.Sex,
		&u.Interests,
		&u.City,
		&u.EncryptedPassword,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return u, nil
}

// Update ...
func (r *UserRepository) Update(u *model.User) error {

	if err := u.BeforeCreate(); err != nil {
		return err
	}

	_, err := r.store.db.Exec(
		`UPDATE users
		 SET name = ?,
			 surname = ?,
			 age = ?,
			 sex = ?,
			 interests = ?,
			 city = ?
		 WHERE id = ?`,
		u.Name,
		u.Surname,
		u.Age,
		u.Sex,
		u.Interests,
		u.City,
		u.ID,
	)

	return err
}

// GetTopUsers ...
func (r *UserRepository) GetTopUsers(n int) ([]*model.User, error) {

	rows, err := r.store.db.Query(
		`SELECT id,
		        name,
				surname,
				sex,
				age,
				city,
				interests
		 FROM users
		 ORDER BY name
		 limit ?`,
		n,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*model.User, 0)
	for rows.Next() {

		user := &model.User{}
		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Surname,
			&user.Sex,
			&user.Age,
			&user.City,
			&user.Interests,
		); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

// GetFriendsList ...
func (r *UserRepository) GetFriendsList(id int) ([]*model.User, error) {

	rows, err := r.store.db.Query(
		`SELECT id,
		        name,
				surname,
				age,
				sex,
				city,
				interests
		 FROM users
		 WHERE id IN (SELECT friend_id
					  FROM friends
					  WHERE user_id = ?
					    AND is_accepted = true)`,
		id,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*model.User, 0)
	for rows.Next() {
		user := &model.User{}
		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Surname,
			&user.Age,
			&user.Sex,
			&user.City,
			&user.Interests,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// GetFriendsRequests ...
func (r *UserRepository) GetFriendsRequests(id int) ([]*model.User, error) {

	rows, err := r.store.db.Query(
		`SELECT id,
		        name,
				surname,
				age,
				sex,
				city,
				interests
		 FROM users
		 WHERE id IN (SELECT friend_id
					  FROM friends
					  WHERE user_id = ?
					    AND is_accepted = false)`,
		id,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*model.User, 0)
	for rows.Next() {
		user := &model.User{}
		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Surname,
			&user.Age,
			&user.Sex,
			&user.City,
			&user.Interests,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// RequestWasAlreadySent ...
func (r *UserRepository) RequestWasAlreadySent(fromID, toID int) bool {

	var userID, friendID int
	err := r.store.db.QueryRow(
		`SELECT user_id, friend_id
		 FROM friends
		 WHERE user_id = ? AND friend_id = ?`,
		fromID,
		toID,
	).Scan(&userID, &friendID)

	if err != nil {
		return false
	}

	return true
}

// SendFriendRequest ...
func (r *UserRepository) SendFriendRequest(fromID, toID int) error {

	if requestwasAlreadySent := r.RequestWasAlreadySent(fromID, toID); requestwasAlreadySent {
		return store.ErrFriendRequestWasAlreadySent
	}

	_, err := r.store.db.Exec(
		`INSERT INTO friends (user_id, friend_id, is_accepted)
		 VALUES (?, ?, ?), (?, ?, ?)`,
		fromID, toID, false,
		toID, fromID, false,
	)

	return err
}

// AcceptFriendRequest ...
func (r *UserRepository) AcceptFriendRequest(fromID, toID int) error {

	_, err := r.store.db.Exec(
		`UPDATE friends
		 SET is_accepted = true
		 WHERE user_id IN (?, ?) AND friend_id IN (?, ?)`,
		fromID, toID, toID, fromID,
	)

	return err
}
