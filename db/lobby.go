package db

import (
	"strings"

	"github.com/xedom/codeduel/types"
)

func (m *MariaDB) CreateLobby(lobby *types.Lobby) error {
	query := `INSERT INTO lobby (uuid, challenge_id, owner_id, max_players, game_duration, allowed_languages)
		VALUES (?, ?, ?, ?, ?, ?);
	;`

	allowedLanguages := ""
	for i, lang := range lobby.AllowedLanguages {
		if i != 0 {
			allowedLanguages += ","
		}
		allowedLanguages += lang
	}

	_, err := m.db.Exec(
		query,
		lobby.UniqueId,
		lobby.ChallengeId,
		lobby.OwnerId,
		lobby.MaxPlayers,
		lobby.GameDuration,
		allowedLanguages,
	)

	if err != nil {
		return err
	}

	id, err := m.getLastInsertID()
	if err != nil {
		return err
	}

	lobby.Id = id
	return err
}

func (m *MariaDB) CreateLobbyUserSubmission(userSubmission *types.LobbyUser) error {
	query := `INSERT INTO lobby_user (lobby_id, user_id, code, language, tests_passed, submission_date)
		VALUES (?, ?, ?, ?, ?, ?);
	;`
	_, err := m.db.Exec(
		query,
		userSubmission.LobbyId,
		userSubmission.UserId,
		userSubmission.Code,
		userSubmission.Language,
		userSubmission.TestsPassed,
		userSubmission.SubmissionDate,
	)
	if err != nil {
		return err
	}

	id, err := m.getLastInsertID()
	if err != nil {
		return err
	}

	userSubmission.Id = id
	// TODO: Update userSubmission.CreatedAt and userSubmission.UpdatedAt
	return err
}

func (m *MariaDB) GetLobbyByUniqueId(uniqueId string) (*types.Lobby, error) {
	query := `SELECT id, uuid, challenge_id, owner_id, status, max_players, game_duration, allowed_languages, created_at, updated_at
		FROM lobby WHERE uuid = ?;`

	row := m.db.QueryRow(query, uniqueId)
	allowLanguages := ""
	lobby := &types.Lobby{}
	err := row.Scan(
		&lobby.Id,
		&lobby.UniqueId,
		&lobby.ChallengeId,
		&lobby.OwnerId,
		&lobby.Status,
		&lobby.MaxPlayers,
		&lobby.GameDuration,
		&allowLanguages,
		&lobby.CreatedAt,
		&lobby.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	lobby.AllowedLanguages = []string{}
	for _, lang := range allowLanguages {
		lobby.AllowedLanguages = append(lobby.AllowedLanguages, string(lang))
	}

	if err := row.Err(); err != nil {
		return nil, err
	}

	return lobby, nil
}

func (m *MariaDB) EndLobby(lobbyUniqueId string) error {
	query := `UPDATE lobby SET status = 'closed' WHERE uuid = ?;`
	_, err := m.db.Exec(query, lobbyUniqueId)
	if err != nil {
		return err
	}

	return nil
}

func (m *MariaDB) GetLobbyResults(lobbyUniqueId string) (*types.LobbyResults, error) {
	query := `SELECT l.id, l.uuid, l.challenge_id, l.owner_id, l.status, l.max_players, l.game_duration, l.allowed_languages, l.created_at, l.updated_at,
		u.id, u.lobby_id, u.user_id, u.code, u.language, u.tests_passed, u.submission_date, u.created_at, u.updated_at
		FROM lobby l
		JOIN lobby_user u ON l.id = u.lobby_id
		WHERE l.uuid = ?;`

	rows, err := m.db.Query(query, lobbyUniqueId)
	if err != nil {
		return nil, err
	}

	lobby := &types.Lobby{}
	results := []types.LobbyUserResult{}
	for rows.Next() {
		allowLanguages := ""
		user := types.LobbyUserResult{}
		err := rows.Scan(
			&lobby.Id,
			&lobby.UniqueId,
			&lobby.ChallengeId,
			&lobby.OwnerId,
			&lobby.Status,
			&lobby.MaxPlayers,
			&lobby.GameDuration,
			&allowLanguages,
			&lobby.CreatedAt,
			&lobby.UpdatedAt,

			&user.Id,
			&user.LobbyId,
			&user.UserId,
			&user.Code,
			&user.Language,
			&user.TestsPassed,
			&user.SubmissionDate,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		lobby.AllowedLanguages = []string{}
		langSplided := strings.Split(allowLanguages, ",")
		for _, lang := range langSplided {
			lobby.AllowedLanguages = append(lobby.AllowedLanguages, string(lang))
		}

		results = append(results, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	return &types.LobbyResults{
		Lobby:   *lobby,
		Results: results,
	}, nil
}

// -- Init Tables --
func (m *MariaDB) InitLobbyTables() error {
	if err := m.createTableLobby(); err != nil {
		return err
	}
	if err := m.createTableLobbyUser(); err != nil {
		return err
	}

	return nil
}

func (m *MariaDB) createTableLobby() error {
	query := `CREATE TABLE IF NOT EXISTS lobby (
		id INT AUTO_INCREMENT,
		uuid VARCHAR(255) NOT NULL,
		challenge_id INT NOT NULL,
		owner_id INT NOT NULL,
		status VARCHAR(50) DEFAULT 'open',
		
		mode_id INT NOT NULL,
		max_players INT NOT NULL,
		game_duration INT NOT NULL,
		allowed_languages VARCHAR(255) NOT NULL,

		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

		PRIMARY KEY (id),
		FOREIGN KEY (challenge_id) REFERENCES challenge(id),
		FOREIGN KEY (owner_id) REFERENCES user(id),
		UNIQUE INDEX (id),
		UNIQUE INDEX (uuid)
	);`
	_, err := m.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (m *MariaDB) createTableLobbyUser() error {
	query := `CREATE TABLE IF NOT EXISTS lobby_user (
		id INT AUTO_INCREMENT,
		lobby_id INT NOT NULL,
		user_id INT NOT NULL,
		
		code TEXT NOT NULL,
		language VARCHAR(50) NOT NULL,
		tests_passed INT NOT NULL,
		submission_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

		PRIMARY KEY (id),
		UNIQUE INDEX (id),
		UNIQUE INDEX (lobby_id, user_id),
		FOREIGN KEY (lobby_id) REFERENCES lobby(id),
		FOREIGN KEY (user_id) REFERENCES user(id)
	);`
	_, err := m.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
