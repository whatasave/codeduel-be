package db

import (
	"log"
	"strings"

	"github.com/xedom/codeduel/types"
	"github.com/xedom/codeduel/utils"
)

func (m *MariaDB) CreateLobby(lobby *types.Lobby) error {
	query := `INSERT INTO lobby (uuid, challenge_id, owner_id, mode, max_players, game_duration, allowed_languages)
		VALUES (?, ?, ?, ?, ?, ?, ?);
	;`

	allowedLanguages := ""
	for i, lang := range lobby.AllowedLanguages {
		if i != 0 {
			allowedLanguages += ","
		}
		allowedLanguages += lang
	}

	res, err := m.db.Exec(
		query,
		lobby.UniqueId,
		lobby.ChallengeId,
		lobby.OwnerId,
		lobby.Mode,
		lobby.MaxPlayers,
		lobby.GameDuration,
		allowedLanguages,
	)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	lobby.Id = int(id)

	for _, userId := range lobby.UsersId {
		if err := m.CreateLobbyUser(userId, lobby.Id); err != nil {
			return err
		}
	}

	return err
}

func (m *MariaDB) CreateLobbyUser(userId int, lobbyId int) error {
	query := `INSERT INTO lobby_user (lobby_id, user_id) VALUES (?, ?);`
	_, err := m.db.Exec(query, lobbyId, userId)

	return err
}

func (m *MariaDB) UpdateLobbyUserSubmission(userSubmission *types.LobbyUser) error {
	// query := `INSERT INTO lobby_user (lobby_id, user_id, code, language, tests_passed, submitted_at)
	// 	VALUES (?, ?, ?, ?, ?, ?);
	// ;`
	query := `UPDATE lobby_user SET code = ?, language = ?, tests_passed = ?, submitted_at = ?
		WHERE lobby_id = ? AND user_id = ?;`

	_, err := m.db.Exec(
		query,
		userSubmission.Code,
		userSubmission.Language,
		userSubmission.TestsPassed,
		userSubmission.SubmittedAt,
		userSubmission.LobbyId,
		userSubmission.UserId,
	)

	return err
}

func (m *MariaDB) GetLobbyByUniqueId(uniqueId string) (*types.Lobby, error) {
	query := `SELECT id, uuid, challenge_id, owner_id, ended, max_players, game_duration, allowed_languages, created_at, updated_at
		FROM lobby WHERE uuid = ?;`

	row := m.db.QueryRow(query, uniqueId)
	allowLanguages := ""
	lobby := &types.Lobby{}
	err := row.Scan(
		&lobby.Id,
		&lobby.UniqueId,
		&lobby.ChallengeId,
		&lobby.OwnerId,
		&lobby.Ended,
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

func (m *MariaDB) GetLobbyResults(lobbyUniqueId string) (*types.LobbyResults, error) {
	query := `SELECT
		l.id, l.uuid, l.challenge_id, l.owner_id, l.ended, l.mode, l.max_players, l.game_duration, l.allowed_languages, l.created_at, l.updated_at,
		u.id, u.lobby_id, u.user_id, u.code, u.language, u.tests_passed, u.show_code, u.submitted_at, u.created_at, u.updated_at
		FROM lobby l
		JOIN lobby_user u ON l.id = u.lobby_id
		WHERE l.uuid = ?;`

	rows, err := m.db.Query(query, lobbyUniqueId)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("%s DB(GetLobbyResults): %s", utils.GetLogTag("DB"), err)
		}
	}()

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
			&lobby.Ended,
			&lobby.Mode,
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
			&user.ShowCode,
			&user.SubmittedAt,
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

	return &types.LobbyResults{
		Lobby:   *lobby,
		Results: results,
	}, nil
}

func (m *MariaDB) EndLobby(lobbyUniqueId string) error {
	query := `UPDATE lobby SET ended = TRUE WHERE uuid = ?;`
	_, err := m.db.Exec(query, lobbyUniqueId)

	return err
}

func (m *MariaDB) UpdateShareLobbyCode(lobbyId int, userId int, showCode bool) error {
	query := `UPDATE lobby_user SET show_code = ? WHERE lobby_id = ? AND user_id = ?;`
	_, err := m.db.Exec(query, showCode, lobbyId, userId)

	return err
}

func (m *MariaDB) GetMatchByUsername(username string) ([]*types.SingleMatchResult, error) {
	query := `
	SELECT 
		l.id AS lobby_id, l.uuid AS lobby_uuid, l.created_at AS lobby_created_at,
		l.mode AS lobby_mode, l.max_players AS lobby_max_players, l.game_duration AS lobby_game_duration, l.allowed_languages AS lobby_allowed_languages,
		l.challenge_id, ch.title AS challenge_title, ch.description AS challenge_description,
		own.id AS challenge_owner_id, own.username AS challenge_owner_username, own.name AS challenge_owner_name, own.avatar AS challenge_owner_avatar,
		us.id AS player_id, us.username AS player_username, us.name AS player_name, us.avatar AS player_avatar,
		u.code AS player_code, u.language AS player_language, u.tests_passed AS player_tests_passed, u.show_code AS player_show_code, u.submitted_at AS player_submitted_at
	FROM lobby l
	JOIN lobby_user u ON l.id = u.lobby_id AND l.ended = 1
	JOIN user us ON u.user_id = us.id
	JOIN challenge ch ON l.challenge_id = ch.id
	JOIN user own ON ch.owner_id = own.id
	WHERE us.username = ?;`

	rows, err := m.db.Query(query, username)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("%s DB(GetMatchByUsername): %s", utils.GetLogTag("DB"), err)
		}
	}()

	matches := []*types.SingleMatchResult{}
	for rows.Next() {
		match := &types.SingleMatchResult{}
		allowLanguages := ""
		err := rows.Scan(
			&match.Match.Id,
			&match.Match.UniqueId,
			&match.Match.CreatedAt,
			&match.Match.Mode,
			&match.Match.MaxPlayers,
			&match.Match.Duration,
			&allowLanguages,

			&match.Challenge.Id,
			&match.Challenge.Title,
			&match.Challenge.Description,

			&match.Challenge.Owner.Id,
			&match.Challenge.Owner.Username,
			&match.Challenge.Owner.Name,
			&match.Challenge.Owner.Avatar,

			&match.Player.Id,
			&match.Player.Username,
			&match.Player.Name,
			&match.Player.Avatar,
			&match.Player.Code,
			&match.Player.Language,
			&match.Player.TestsPassed,
			&match.Player.ShowCode,
			&match.Player.SubmittedAt,
		)
		if err != nil {
			return nil, err
		}

		match.Match.AllowedLanguages = []string{}
		languages := strings.Split(allowLanguages, ",")
		for _, lang := range languages {
			match.Match.AllowedLanguages = append(match.Match.AllowedLanguages, string(lang))
		}

		matches = append(matches, match)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return matches, nil
}

// -- Init Tables --
func (m *MariaDB) InitLobbyTables() []MigrationFunc {
	return []MigrationFunc{
		m.createTableMode,
		m.createTableLanguage,
		m.createTableLobby,
		m.createTableLobbyUser,
	}
}

func (m *MariaDB) createTableLobby() error {
	query := `CREATE TABLE IF NOT EXISTS lobby (
		id INT AUTO_INCREMENT,
		uuid VARCHAR(255) NOT NULL,
		challenge_id INT NOT NULL,
		owner_id INT NOT NULL,
		ended BOOLEAN NOT NULL DEFAULT FALSE,
		
		mode VARCHAR(50) NOT NULL,
		max_players INT NOT NULL,
		game_duration INT NOT NULL,
		allowed_languages VARCHAR(255) NOT NULL,

		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

		PRIMARY KEY (id),
		FOREIGN KEY (challenge_id) REFERENCES challenge(id),
		FOREIGN KEY (owner_id) REFERENCES user(id),
		UNIQUE INDEX (id),
		UNIQUE INDEX (uuid)
	);`
	_, err := m.db.Exec(query)

	return err
}

func (m *MariaDB) createTableLobbyUser() error {
	query := `CREATE TABLE IF NOT EXISTS lobby_user (
		id INT AUTO_INCREMENT,
		lobby_id INT NOT NULL,
		user_id INT NOT NULL,
		
		code TEXT,
		language VARCHAR(50),
		tests_passed INT NOT NULL DEFAULT 0,
		show_code BOOLEAN NOT NULL DEFAULT FALSE,
		match_rank INT,
		submitted_at TIMESTAMP,

		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

		PRIMARY KEY (id),
		FOREIGN KEY (lobby_id) REFERENCES lobby(id),
		FOREIGN KEY (user_id) REFERENCES user(id),
		UNIQUE INDEX (id),
		UNIQUE INDEX (lobby_id, user_id)
	);`
	_, err := m.db.Exec(query)
	return err
}

func (m *MariaDB) createTableMode() error {
	query := `CREATE TABLE IF NOT EXISTS mode (
		id INT unique AUTO_INCREMENT,
		name VARCHAR(50) NOT NULL,
		description VARCHAR(255) NOT NULL,
		
		PRIMARY KEY (id),
		UNIQUE INDEX (id)
	);`
	_, err := m.db.Exec(query)
	if err != nil {
		return err
	}

	queryDefaultValues := `INSERT IGNORE INTO mode
	(id, name, description) VALUES	
	(1, 'speed', 'The shortest time wins.'),
	(2, 'size', 'The shortest code wins.'),
	(3, 'efficiency', 'The most efficient code wins.'),
	(4, 'memory', 'The most memory efficient code wins.'),
	(5, 'readability', 'The most readable code wins.'),
	(6, 'style', 'The most stylish code wins.');`

	_, err = m.db.Exec(queryDefaultValues)
	return err
}

func (m *MariaDB) createTableLanguage() error {
	query := `CREATE TABLE IF NOT EXISTS language (
		id INT unique AUTO_INCREMENT,
		name VARCHAR(50) NOT NULL,

		PRIMARY KEY (id),
		UNIQUE INDEX (id),
		UNIQUE INDEX (name)
	);`
	_, err := m.db.Exec(query)
	if err != nil {
		return err
	}

	queryDefaultValues := `INSERT IGNORE INTO language
	(id, name) VALUES
	(0, 'c'),
	(1, 'cpp'),
	(2, 'java'),
	(3, 'js'),
	(4, 'go'),
	(5, 'rs'),
	(6, 'rb'),
	(7, 'py');`
	_, err = m.db.Exec(queryDefaultValues)
	return err
}
