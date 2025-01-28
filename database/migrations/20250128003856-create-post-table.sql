-- +migrate Up
CREATE TABLE posts (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    post_title VARCHAR(250) NOT NULL,
    post_content LONGTEXT NOT NULL,
    post_hastags LONGTEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by LONGTEXT NOT NULL,
    updated_by LONGTEXT NULL

    CONSTRAINT fk_user_id_post FOREIGN KEY (user_id) REFERENCES users(id)
);


-- +migrate Down
DROP TABLE posts;
		