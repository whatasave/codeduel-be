package db

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/xedom/codeduel/types"
	"github.com/xedom/codeduel/utils"
)

func (m *MariaDB) GetChallenges() (*[]types.Challenge, error) {
	query := "SELECT id, owner_id, title, description, content, tests, created_at, updated_at FROM `challenge`;"
	rows, err := m.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("%s DB(GetChallenges): %s", utils.GetLogTag("DB"), err)
		}
	}()

	challenges := &[]types.Challenge{}
	for rows.Next() {
		var challenge types.Challenge
		var testCases string

		if err := rows.Scan(
			&challenge.Id,
			&challenge.OwnerId,

			&challenge.Title,
			&challenge.Description,
			&challenge.Content,

			&testCases,

			&challenge.CreatedAt,
			&challenge.UpdatedAt,
		); err != nil {
			return nil, err
		}

		if err := json.Unmarshal([]byte(testCases), &challenge.TestCases); err != nil {
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
	query := "SELECT id, owner_id, title, description, content, tests, created_at, updated_at FROM `challenge` WHERE id = ? LIMIT 1;"
	row := m.db.QueryRow(query, id)
	if row == nil {
		return nil, fmt.Errorf("challenge not found")
	}

	challenge := &types.Challenge{}
	var testCases string

	if err := row.Scan(
		&challenge.Id,
		&challenge.OwnerId,

		&challenge.Title,
		&challenge.Description,
		&challenge.Content,

		&testCases,

		&challenge.CreatedAt,
		&challenge.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("row.Scan | id: %d | %w", id, err)
	}

	if err := json.Unmarshal([]byte(testCases), &challenge.TestCases); err != nil {
		return nil, err
	}

	return challenge, nil
}

func (m *MariaDB) GetChallengesByOwnerID(ownerID int) (*[]types.Challenge, error) {
	query := "SELECT id, owner_id, title, description, content, tests, created_at, updated_at FROM `challenge` WHERE owner_id = ?;"
	rows, err := m.db.Query(query, ownerID)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("%s DB(GetChallengesByOwnerID): %s", utils.GetLogTag("DB"), err)
		}
	}()

	challenges := &[]types.Challenge{}
	for rows.Next() {
		var challenge types.Challenge
		var testCases string
		err := rows.Scan(
			&challenge.Id,
			&challenge.OwnerId,

			&challenge.Title,
			&challenge.Description,
			&challenge.Content,

			&testCases,

			&challenge.CreatedAt,
			&challenge.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal([]byte(testCases), &challenge.TestCases)
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

func (m *MariaDB) GetRandomChallengeFull() (*types.ChallengeFull, error) {
	query := `
	SELECT
	c.id, c.title, c.description, c.content, c.created_at, c.updated_at,
	c.owner_id, u.name, u.username, u.avatar,
	c.tests, c.tests_hidden
	FROM challenge c
	JOIN user u ON c.owner_id = u.id
	ORDER BY RAND()
	LIMIT 1;`

	row := m.db.QueryRow(query)
	if row == nil {
		return nil, fmt.Errorf("challenge not found")
	}

	if row.Err() != nil {
		return nil, row.Err()
	}

	challenge := &types.ChallengeFull{}
	var testCases, hiddenTestCases string

	err := row.Scan(
		&challenge.Id,
		&challenge.Title,
		&challenge.Description,
		&challenge.Content,
		&challenge.CreatedAt,
		&challenge.UpdatedAt,
		&challenge.Owner.Id,
		&challenge.Owner.Name,
		&challenge.Owner.Username,
		&challenge.Owner.Avatar,
		&testCases,
		&hiddenTestCases,
	)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(testCases), &challenge.TestCases)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(hiddenTestCases), &challenge.HiddenTestCases)
	if err != nil {
		return nil, err
	}

	return challenge, nil
}

func (m *MariaDB) GetChallengeByIDFull(id int) (*types.ChallengeFull, error) {
	query := `
	SELECT
	c.id, c.title, c.description, c.content, c.created_at, c.updated_at,
	c.owner_id, u.name, u.username, u.avatar,
	c.tests, c.tests_hidden
	FROM challenge c
	JOIN user u ON c.owner_id = u.id
	WHERE c.id = ?
	LIMIT 1;`

	row := m.db.QueryRow(query, id)
	if row == nil {
		return nil, fmt.Errorf("challenge not found")
	}

	if row.Err() != nil {
		return nil, row.Err()
	}

	challenge := &types.ChallengeFull{}
	var testCases, hiddenTestCases string

	err := row.Scan(
		&challenge.Id,
		&challenge.Title,
		&challenge.Description,
		&challenge.Content,
		&challenge.CreatedAt,
		&challenge.UpdatedAt,
		&challenge.Owner.Id,
		&challenge.Owner.Name,
		&challenge.Owner.Username,
		&challenge.Owner.Avatar,
		&testCases,
		&hiddenTestCases,
	)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(testCases), &challenge.TestCases)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(hiddenTestCases), &challenge.HiddenTestCases)
	if err != nil {
		return nil, err
	}

	return challenge, nil
}

func (m *MariaDB) CreateChallenge(challenge *types.Challenge) error {
	query := "INSERT INTO `challenge` (owner_id, title, description, content) VALUES (?, ?, ?, ?);"
	res, err := m.db.Exec(query, challenge.OwnerId, challenge.Title, challenge.Description, challenge.Content)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	challenge.Id = int(id)
	return nil
}

func (m *MariaDB) UpdateChallenge(challenge *types.Challenge) error {
	query := "UPDATE `challenge` SET title = ?, description = ?, content = ? WHERE id = ?;"
	_, err := m.db.Exec(query, challenge.Title, challenge.Description, challenge.Content, challenge.Id)
	return err
}

func (m *MariaDB) DeleteChallenge(id int) error {
	query := "DELETE FROM `challenge` WHERE id = ?;"
	_, err := m.db.Exec(query, id)
	return err
}

// -- Init Tables --
func (m *MariaDB) InitChallengeTables() []MigrationFunc {
	return []MigrationFunc{
		m.createTableChallenge,
	}
}

func (m *MariaDB) createTableChallenge() error {
	query := `CREATE TABLE IF NOT EXISTS ` + "`challenge`" + ` (
		id INT unique AUTO_INCREMENT,
		owner_id INT NOT NULL,
		title VARCHAR(50) NOT NULL,
		description VARCHAR(255) NOT NULL,
		content LONGTEXT NOT NULL,

		tests JSON NULL DEFAULT NULL,
		tests_hidden JSON NULL DEFAULT NULL,

		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		
		PRIMARY KEY (id),
		FOREIGN KEY (owner_id) REFERENCES user(id),
		UNIQUE INDEX (id)
	);`
	_, err := m.db.Exec(query)
	return err
}
