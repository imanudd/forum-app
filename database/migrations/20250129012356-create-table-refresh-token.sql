-- +migrate Up
CREATE TABLE refresh_tokens(
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    refresh_token TEXT NOT NULL,
    expired_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by LONGTEXT NOT NULL,
    updated_by LONGTEXT NULL,

    CONSTRAINT refresh_token_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id)
);

-- +migrate Down
DROP TABLE refresh_tokens;
		