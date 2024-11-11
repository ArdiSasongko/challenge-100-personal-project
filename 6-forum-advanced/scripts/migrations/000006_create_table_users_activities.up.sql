CREATE TABLE IF NOT EXISTS users_activities (
    id SMALLSERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    content_id INT NOT NULL,
    is_liked BOOLEAN NOT NULL DEFAULT FALSE,
    is_saved BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255) NOT NULL,
    updated_by VARCHAR(255) NOT NULL,
    CONSTRAINT fk_user_id_users_activities FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_content_id_users_activities FOREIGN KEY (content_id) REFERENCES contents(id) ON DELETE CASCADE
);