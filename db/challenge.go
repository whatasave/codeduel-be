package db

import (
	"database/sql"

	"github.com/xedom/codeduel/types"
)

func (m *MariaDB) createChallengeTable() error {
	query := `CREATE TABLE IF NOT EXISTS ` + "`challenge`" + ` (
		id INT unique AUTO_INCREMENT,
		owner_id INT NOT NULL,
		title VARCHAR(50) NOT NULL,
		description VARCHAR(255) NOT NULL,
		content TEXT NOT NULL,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		
		PRIMARY KEY (id),
		FOREIGN KEY (owner_id) REFERENCES user(id),
		UNIQUE INDEX (id)
	);`
	_, err := m.db.Exec(query)
	return err
}

func (m *MariaDB) GetChallenges() (*[]types.Challenge, error) {
	query := `SELECT * FROM ` + "`challenge`"
	rows, err := m.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	challenges := &[]types.Challenge{}
	for rows.Next() {
		var challenge types.Challenge
		err := rows.Scan(&challenge.ID, &challenge.OwnerID, &challenge.Title, &challenge.Description, &challenge.Content, &challenge.CreatedAt, &challenge.UpdatedAt)
		if err != nil {
			return nil, err
		}
		*challenges = append(*challenges, challenge)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	
	return challenges, nil
}

func (m *MariaDB) GetChallengeByID(id int) (*types.Challenge, error) {
	query := `SELECT * FROM ` + "`challenge`" + ` WHERE id = ?;`
	rows, err := m.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		return m.parseChallenge(rows)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return nil, nil
}

func (m *MariaDB) CreateChallenge(challenge *types.Challenge) error {
	query := `INSERT INTO ` + "`challenge`" + ` (owner_id, title, description, content) VALUES (?, ?, ?, ?);`
	res, err := m.db.Exec(query, challenge.OwnerID, challenge.Title, challenge.Description, challenge.Content)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	challenge.ID = int(id)
	return nil
}

func (m *MariaDB) UpdateChallenge(challenge *types.Challenge) error {
	query := `UPDATE ` + "`challenge`" + ` SET title = ?, description = ?, content = ? WHERE id = ?;`
	_, err := m.db.Exec(query, challenge.Title, challenge.Description, challenge.Content, challenge.ID)
	return err
}

func (m *MariaDB) DeleteChallenge(id int) error {
	query := `DELETE FROM ` + "`challenge`" + ` WHERE id = ?;`
	_, err := m.db.Exec(query, id)
	return err
}

func (m *MariaDB) parseChallenge(rows *sql.Rows) (*types.Challenge, error) {
	var challenge types.Challenge
	err := rows.Scan(&challenge.ID, &challenge.OwnerID, &challenge.Title, &challenge.Description, &challenge.Content, &challenge.CreatedAt, &challenge.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &challenge, nil
}

func (m *MariaDB) GetChallengesByOwnerID(ownerID int) (*[]types.Challenge, error) {
	query := `SELECT * FROM ` + "`challenge`" + ` WHERE owner_id = ?;`
	rows, err := m.db.Query(query, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	challenges := &[]types.Challenge{}
	for rows.Next() {
		var challenge types.Challenge
		err := rows.Scan(&challenge.ID, &challenge.OwnerID, &challenge.Title, &challenge.Description, &challenge.Content, &challenge.CreatedAt, &challenge.UpdatedAt)
		if err != nil {
			return nil, err
		}
		*challenges = append(*challenges, challenge)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	
	return challenges, nil
}
