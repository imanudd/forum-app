-- +migrate Up
CREATE TABLE user_activities (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    post_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    is_liked BOOL NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by LONGTEXT NOT NULL,
    updated_by LONGTEXT NULL,

    CONSTRAINT user_activities_post_id_fkey FOREIGN KEY (post_id) REFERENCES posts(id),
    CONSTRAINT user_activities_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id)
);

-- +migrate Down
DROP TABLE user_activities;
		