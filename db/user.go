package db

import (
	"database/sql"
	"fmt"

	"github.com/xedom/codeduel/types"
)

func (m *MariaDB) CreateUser(user *types.User) error {
	query := `INSERT INTO user (username, email, avatar)
		VALUES (?, ?, ?);
	;`
	_, err := m.db.Exec(query, user.Username, user.Email, user.Avatar)
	if err != nil {
		return err
	}

	id, err := m.getLastInsertID()
	if err != nil {
		return err
	}

	user.ID = id
	return err
}

func (m *MariaDB) getLastInsertID() (int, error) {
	row, err := m.db.Query(`SELECT LAST_INSERT_ID();`)
	if err != nil {
		return 0, err
	}

	var id int
	for row.Next() {
		if err := row.Scan(&id); err != nil {
			return 0, err
		}
	}
	if err := row.Err(); err != nil {
		return 0, err
	}

	return id, row.Close()
}

func (m *MariaDB) CreateAuth(auth *types.AuthEntry) error {
	query := `INSERT INTO auth (user_id, provider, provider_id) VALUES (?, ?, ?);`
	_, err := m.db.Exec(query, auth.UserID, auth.Provider, auth.ProviderID)
	if err != nil {
		return err
	}

	id, err := m.getLastInsertID()
	if err != nil {
		return err
	}

	auth.ID = id
	return err
}

func (m *MariaDB) GetAuthByProviderAndID(provider, providerID string) (*types.AuthEntry, error) {
	query := `SELECT * FROM auth WHERE provider = ? AND provider_id = ?;`
	rows, err := m.db.Query(query, provider, providerID)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		if err := rows.Close(); err != nil {
			fmt.Println("DB(GetAuthByProviderAndID):", err)
		}
	}(rows)

	for rows.Next() {
		return m.parseAuth(rows)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("auth with provider_id %s not found", providerID)
}

func (m *MariaDB) DeleteUser(id int) error {
	query := `DELETE FROM user WHERE id = ?;`
	res, err := m.db.Exec(query, id)

	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return fmt.Errorf("user with id %d not found", id)
	}

	return err
}

func (m *MariaDB) DeleteUserByUsername(username string) error {
	query := `DELETE FROM user WHERE username = ?;`
	res, err := m.db.Exec(query, username)

	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return fmt.Errorf("DB(DeleteUserByUsername): user with id %s not found", username)
	}

	return err
}

func (m *MariaDB) UpdateUser(user *types.User) error {
	query := `UPDATE user SET username = ?, email = ?, avatar = ? WHERE id = ?;`
	res, err := m.db.Exec(query, user.Username, user.Email, user.Avatar, user.ID)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return fmt.Errorf("user with id %d not found", user.ID)
	}

	return err
}

func (m *MariaDB) GetUsers() ([]*types.UserResponse, error) {
	query := `SELECT user.name, username, avatar, background_img, bio, created_at FROM user;`
	rows, err := m.db.Query(query)
	if err != nil {
		return nil, err
	}

	var users []*types.UserResponse
	for rows.Next() {
		user, err := m.parseUserResponse(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, rows.Close()
}

func (m *MariaDB) GetUserByID(id int) (*types.User, error) {
	query := `SELECT * FROM user WHERE id = ?;`
	rows, err := m.db.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("DB(GetUserByID): %s", err.Error())
	}
	defer func(rows *sql.Rows) {
		if err := rows.Close(); err != nil {
			fmt.Println("DB(GetUserByID):", err)
		}
	}(rows)

	for rows.Next() {
		return m.parseUser(rows)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("DB(GetUserByID): %s", err.Error())
	}

	return nil, fmt.Errorf("DB(GetUserByID): user with id %d not found", id)
}

func (m *MariaDB) GetUserByUsername(username string) (*types.User, error) {
	query := `SELECT * FROM user WHERE username = ?;`
	rows, err := m.db.Query(query, username)
	if err != nil {
		return nil, fmt.Errorf("DB(GetUserByUsername): %s", err.Error())
	}
	defer func(rows *sql.Rows) {
		if err := rows.Close(); err != nil {
			fmt.Println("DB(GetUserByUsername):", err)
		}
	}(rows)

	for rows.Next() {
		return m.parseUser(rows)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("DB(GetUserByUsername): %s", err.Error())
	}

	return nil, fmt.Errorf("DB(GetUserByUsername): user with username %s not found", username)
}

func (m *MariaDB) GetUserStats(id int) ([]*types.UserStatsParsed, error) {
	// query := `SELECT * FROM user_stats WHERE user_id = ?;`
	query := `SELECT
		user_stats.id,
		stats.name,
		user_stats.stat,
		user_stats.created_at,
		user_stats.updated_at
	FROM user_stats
	JOIN stats ON user_stats.stats_id = stats.id
	WHERE user_id = ?;`
	rows, err := m.db.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("DB(GetUserStats:0): %s", err.Error())
	}

	var stats []*types.UserStatsParsed

	for rows.Next() {
		stat := &types.UserStatsParsed{}
		if err := rows.Scan(
			&stat.ID,
			&stat.Name,
			&stat.Stat,
			&stat.CreatedAt,
			&stat.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("DB(GetUserStats:1): %s", err.Error())
		}
		stats = append(stats, stat)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("DB(GetUserStats:2): %s", err.Error())
	}

	return stats, rows.Close()
}

func (m *MariaDB) InitUserTables() error {
	if err := m.createUserTable(); err != nil {
		return err
	}
	if err := m.createAuthTable(); err != nil {
		return err
	}
	if err := m.createStatsTable(); err != nil {
		return err
	}
	if err := m.createUserStatsTable(); err != nil {
		return err
	}

	return nil
}

func (m *MariaDB) createUserTable() error {
	query := `CREATE TABLE IF NOT EXISTS user (
		id INT unique AUTO_INCREMENT,
		username VARCHAR(50) NOT NULL,
		name VARCHAR(50) DEFAULT '',
		email VARCHAR(50) NOT NULL,
		avatar VARCHAR(255),
		background_img VARCHAR(255) DEFAULT '',
		bio TEXT DEFAULT (''),
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

func (m *MariaDB) createUserStatsTable() error {
	query := `CREATE TABLE IF NOT EXISTS user_stats (
		id INT AUTO_INCREMENT,
		user_id INT NOT NULL,
		stats_id INT NOT NULL,
		stat VARCHAR(50) NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		PRIMARY KEY (id),
		FOREIGN KEY (user_id) REFERENCES user(id),
		FOREIGN KEY (stats_id) REFERENCES stats(id)
	);`
	_, err := m.db.Exec(query)
	return err
}

func (m *MariaDB) createStatsTable() error {
	query := `CREATE TABLE IF NOT EXISTS stats (
		id INT AUTO_INCREMENT,
		name VARCHAR(50) NOT NULL,

		PRIMARY KEY (id),
		UNIQUE INDEX (id),
		UNIQUE INDEX (name)
	);`
	_, err := m.db.Exec(query)
	if err != nil {
		return err
	}

	defaultStats := []string{"Games", "Wins", "Top 3"}

	for _, stat := range defaultStats {
		query := `INSERT IGNORE INTO stats (name) VALUES (?);`
		_, err := m.db.Exec(query, stat)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *MariaDB) parseUser(row *sql.Rows) (*types.User, error) {
	user := &types.User{}
	user_avatar := sql.NullString{}
	if err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Name,
		&user.Email,
		&user.Avatar,
		&user.BackgroundImg,
		&user.Bio,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return nil, err
	}
	if user_avatar.Valid {
		user.Avatar = user_avatar.String
	}

	return user, nil
}

func (m *MariaDB) parseUserResponse(row *sql.Rows) (*types.UserResponse, error) {
	user := &types.UserResponse{}
	user_avatar := sql.NullString{}
	if err := row.Scan(
		&user.Name,
		&user.Username,
		&user.Avatar,
		&user.BackgroundImg,
		&user.Bio,
		&user.CreatedAt,
	); err != nil {
		return nil, err
	}
	if user_avatar.Valid {
		user.Avatar = user_avatar.String
	}

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
	); err != nil {
		return nil, err
	}

	return auth, nil
}
