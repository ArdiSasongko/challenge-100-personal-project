CREATE TABLE IF NOT EXISTS comments (
    id SERIAL PRIMARY KEY,
    content_id INT NOT NULL,
    user_id INT NOT NULL,
    comment TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255) NOT NULL,
    updated_by VARCHAR(255) NOT NULL,
    CONSTRAINT fk_user_id_comments FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_content_id_comments FOREIGN KEY (content_id) REFERENCES contents(id) ON DELETE CASCADE    
);