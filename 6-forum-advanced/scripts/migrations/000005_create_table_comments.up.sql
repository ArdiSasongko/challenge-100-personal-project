CREATE TABLE IF NOT EXISTS comments (
    id SMALLSERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    content_id INT NOT NULL,
    comment_body TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255) NOT NULL,
    updated_by VARCHAR(255) NOT NULL,
    CONSTRAINT fk_user_id_comment FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_content_id_comment FOREIGN KEY (content_id) REFERENCES contents(id) ON DELETE CASCADE
);