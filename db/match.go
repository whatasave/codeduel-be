package db

// -- Init Tables --
func (m *MariaDB) InitMatchTables() error {
	if err := m.createTableStatus(); err != nil {
		return err
	}
	if err := m.createTableMatchStatus(); err != nil {
		return err
	}
	if err := m.createTableMode(); err != nil {
		return err
	}
	if err := m.createTableLanguage(); err != nil {
		return err
	}
	if err := m.createTableChallenge(); err != nil {
		return err
	}

	if err := m.createTableMatch(); err != nil {
		return err
	}
	if err := m.createTableMatchUserLink(); err != nil {
		return err
	}

	return nil
}

func (m *MariaDB) createTableMatch() error {
	query := `CREATE TABLE IF NOT EXISTS ` + "`match`" + ` (
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
		FOREIGN KEY (challenge_id) REFERENCES ` + "`challenge`" + `(id),
		FOREIGN KEY (owner_id) REFERENCES user(id),
		FOREIGN KEY (mode_id) REFERENCES mode(id),
		UNIQUE INDEX (id),
		UNIQUE INDEX (uuid)
	);`
	_, err := m.db.Exec(query)
	return err
}

func (m *MariaDB) createTableMatchUserLink() error {
	query := `CREATE TABLE IF NOT EXISTS match_user_link (
		id INT AUTO_INCREMENT,
		match_id INT NOT NULL,
		user_id INT NOT NULL,

		code TEXT NOT NULL,
		language_id INT NOT NULL,
		language VARCHAR(50) NOT NULL,
		tests_passed INT NOT NULL,
		submission_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		
		status_id INT NOT NULL,
		duration INT NOT NULL,
		` + "`rank`" + ` INT NOT NULL,
		
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

		PRIMARY KEY (id),
		FOREIGN KEY (match_id) REFERENCES ` + "`match`" + `(id),
		FOREIGN KEY (user_id) REFERENCES user(id),
		FOREIGN KEY (status_id) REFERENCES status(id),
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
	(6, 'style', 'The most stylish code wins.');
	`
	_, err = m.db.Exec(queryDefaultValues)
	return err
}

func (m *MariaDB) createTableStatus() error {
	query := `CREATE TABLE IF NOT EXISTS status (
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

	queryDefaultValues := `INSERT IGNORE INTO status (id, name) VALUES
	(0, 'not ready'),
	(1, 'ready'),
	(2, 'in match'),
	(3, 'abandoned'),
	(4, 'finished');`
	_, err = m.db.Exec(queryDefaultValues)
	return err
}

func (m *MariaDB) createTableMatchStatus() error {
	query := `CREATE TABLE IF NOT EXISTS status (
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

	queryDefaultValues := `INSERT IGNORE INTO status
	(id, name) VALUES
	(0, 'starting'),
	(1, 'ongoing'),
	(2, 'finished');`
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
	(4, 'golang'),
	(5, 'rust'),
	(6, 'ruby'),
	(7, 'python');`
	_, err = m.db.Exec(queryDefaultValues)
	return err
}
