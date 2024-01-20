package db

import (
	"database/sql"
	"fmt"

	"github.com/xedom/codeduel/types"
)

func (m *MariaDB) CreateUser(user *types.User) (error) {
	query := `INSERT INTO user (username, email, image_url)
		VALUES (?, ?, ?);
	;`
	_, err := m.db.Exec(query, user.Username, user.Email, user.ImageURL)
	if err != nil { return err }

	id, err := m.getLastInsertID()
	if err != nil { return err }

	user.ID = id
	return err
}

func (m *MariaDB) getLastInsertID() (int, error) {
	row, err := m.db.Query(`SELECT LAST_INSERT_ID();`)
	if err != nil { return 0, err }
	defer row.Close()

	var id int
	for row.Next() {
		if err := row.Scan(&id); err != nil { return 0, err }
	}
	if err := row.Err(); err != nil { return 0, err }

	return id, nil
}

func (m *MariaDB) CreateAuth(auth *types.AuthEntry) (error) {
	query := `INSERT INTO auth (user_id, provider, provider_id)
		VALUES (?, ?, ?);
	;`
	_, err := m.db.Exec(query, auth.UserID, auth.Provider, auth.ProviderID)
	if err != nil { return err }
	
	id, err := m.getLastInsertID()
	if err != nil { return err }

	auth.ID = id
	return err
}

func (m *MariaDB) GetAuthByProviderAndID(provider, providerID string) (*types.AuthEntry, error) {
	query := `SELECT * FROM auth WHERE provider = ? AND provider_id = ?;`
	rows, err := m.db.Query(query, provider, providerID)
	if err != nil { return nil, err }
	defer rows.Close()

	for rows.Next() { return m.parseAuth(rows) }
	if err := rows.Err(); err != nil { return nil, err }

	return nil, fmt.Errorf("auth with provider_id %s not found", providerID)
}

func (m *MariaDB) DeleteUser(id int) error {
	query := `DELETE FROM user WHERE id = ?;`
	res, err := m.db.Exec(query, id)

	if err != nil { return err }
	if rows, _ := res.RowsAffected(); rows == 0 {
		return fmt.Errorf("user with id %d not found", id)
	}

	return err
}

func (m *MariaDB) UpdateUser(user *types.User) error {
	query := `UPDATE user SET username = ?, email = ?, image_url = ? WHERE id = ?;`
	res, err := m.db.Exec(query, user.Username, user.Email, user.ImageURL, user.ID)
	if err != nil { return err }
	if rows, _ := res.RowsAffected(); rows == 0 {
		return fmt.Errorf("user with id %d not found", user.ID)
	}

	return err
}

func (m *MariaDB) GetUsers() ([]*types.User, error) {
	query := `SELECT * FROM user;`
	rows, err := m.db.Query(query)
	if err != nil { return nil, err }
	defer rows.Close()

	var users []*types.User
	for rows.Next() {
		user, err := m.parseUser(rows)
		if err != nil { return nil, err }
		users = append(users, user)
	}
	if err := rows.Err(); err != nil { return nil, err }

	return users, nil
}

func (m *MariaDB) GetUserByID(id int) (*types.User, error) {
	query := `SELECT * FROM user WHERE id = ?;`
	rows, err := m.db.Query(query, id)
	if err != nil { return nil, err }
	defer rows.Close()

	for rows.Next() { return m.parseUser(rows) }
	if err := rows.Err(); err != nil { return nil, err }

	return nil, fmt.Errorf("user with id %d not found", id)
}

func (m *MariaDB) InitUserTables() error {
	if err := m.createUserTable(); err != nil { return err }
	if err := m.createAuthTable(); err != nil { return err }

	return nil
}

func (m *MariaDB) createUserTable() error {
	query := `CREATE TABLE IF NOT EXISTS user (
		id INT unique AUTO_INCREMENT,
		username VARCHAR(50) NOT NULL,
		email VARCHAR(50) NOT NULL,
		image_url VARCHAR(255),
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		PRIMARY KEY (id),
		UNIQUE INDEX (username)
	);`
	_, err := m.db.Exec(query)
	return err
}

func (m *MariaDB) createAuthTable() error {
	query := `CREATE TABLE IF NOT EXISTS auth (
		id INT AUTO_INCREMENT,
		user_id INT NOT NULL,
		provider VARCHAR(50) NOT NULL,
		provider_id VARCHAR(50) NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		PRIMARY KEY (id),
		FOREIGN KEY (user_id) REFERENCES user(id),
		UNIQUE INDEX (provider_id)
	);`
	_, err := m.db.Exec(query)
	return err
}

func (m *MariaDB) parseUser(row *sql.Rows) (*types.User, error) {
	user := &types.User{}
	user_image_url := sql.NullString{}
	if err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user_image_url,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil { return nil, err }
	if user_image_url.Valid { user.ImageURL = user_image_url.String }

	return user, nil
}

func (m *MariaDB) parseAuth(row *sql.Rows) (*types.AuthEntry, error) {
	auth := &types.AuthEntry{}
	if err := row.Scan(
		&auth.ID,
		&auth.UserID,
		&auth.Provider,
		&auth.ProviderID,
		&auth.CreatedAt,
		&auth.UpdatedAt,
	); err != nil { return nil, err }

	return auth, nil
}
