CREATE TABLE IF NOT EXISTS contents (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    content_title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    content_hastags TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255) NOT NULL,
    updated_by VARCHAR(255) NOT NULL,
    CONSTRAINT fk_user_id_content FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE    
);