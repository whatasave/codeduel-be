package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/xedom/codeduel/types"
	"github.com/xedom/codeduel/utils"
)

func (m *MariaDB) GetChallenges() (*[]types.Challenge, error) {
	query := "SELECT * FROM `challenge`"
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
		err := rows.Scan(&challenge.Id, &challenge.OwnerId, &challenge.Title, &challenge.Description, &challenge.Content, &challenge.CreatedAt, &challenge.UpdatedAt)
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
	query := "SELECT * FROM `challenge` WHERE id = ?;"
	rows, err := m.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("%s DB(GetChallengeByID): %s", utils.GetLogTag("DB"), err)
		}
	}()

	for rows.Next() {
		return m.parseChallenge(rows)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return nil, nil
}

func (m *MariaDB) GetChallengeByIDFull(id int) (*types.ChallengeFull, error) {
	query := `
	SELECT
	c.id, c.title, c.description, c.content, c.created_at, c.updated_at,
	c.owner_id, u.name, u.username, u.avatar
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
	)
	if err != nil {
		return nil, err
	}

	// TODO: Get test cases
	challenge.TestCases = []types.TestCase{
		{Input: "1 2", Output: "3"},
		{Input: "2 3", Output: "5"},
		{Input: "3 4", Output: "7"},
		{Input: "4 5", Output: "9"},
		{Input: "5 6 7", Output: "18"},
		{Input: "6 7 8", Output: "21"},
		{Input: "7 8 9", Output: "24"},
		{Input: "8 9 10", Output: "27"},
	}
	challenge.HiddenTestCases = []types.TestCase{
		{Input: "100 1", Output: "101"},
		{Input: "5 6 ", Output: "11"},
		{Input: " 6 7", Output: "13"},
		{Input: " 7 8 ", Output: "15"},
		{Input: "7 8 2 3", Output: "20"},
	}

	return challenge, nil
}

func (m *MariaDB) GetChallengesID() ([]int, error) {
	query := "SELECT id FROM `challenge`;"
	rows, err := m.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("%s DB(GetChallengesID): %s", utils.GetLogTag("DB"), err)
		}
	}()

	ids := []int{}
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ids, nil
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

func (m *MariaDB) GetChallengesByOwnerID(ownerID int) (*[]types.Challenge, error) {
	query := "SELECT * FROM `challenge` WHERE owner_id = ?;"
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
		err := rows.Scan(&challenge.Id, &challenge.OwnerId, &challenge.Title, &challenge.Description, &challenge.Content, &challenge.CreatedAt, &challenge.UpdatedAt)
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
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		
		PRIMARY KEY (id),
		FOREIGN KEY (owner_id) REFERENCES user(id),
		UNIQUE INDEX (id)
	);`
	_, err := m.db.Exec(query)
	return err
}

// -- Utils --
func (m *MariaDB) parseChallenge(rows *sql.Rows) (*types.Challenge, error) {
	var challenge types.Challenge
	err := rows.Scan(&challenge.Id, &challenge.OwnerId, &challenge.Title, &challenge.Description, &challenge.Content, &challenge.CreatedAt, &challenge.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &challenge, nil
}
