-- +migrate Up
CREATE TABLE users (
	id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	email VARCHAR(250) NOT NULL UNIQUE,
	username VARCHAR(250) NOT NULL UNIQUE,
	password VARCHAR(500) NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	created_by LONGTEXT NOT NULL,
	updated_by LONGTEXT NULL
);

-- +migrate Down
DROP TABLE users;