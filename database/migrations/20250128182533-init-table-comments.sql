-- +migrate Up
CREATE TABLE comments (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    post_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    comment_content LONGTEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by LONGTEXT NOT NULL,
    updated_by LONGTEXT NULL,

    CONSTRAINT fk_post_id_comment FOREIGN KEY (post_id) REFERENCES posts(id),
    CONSTRAINT fk_user_id_comment FOREIGN KEY (user_id) REFERENCES users(id)
); 

-- +migrate Down
DROP TABLE comments;
		