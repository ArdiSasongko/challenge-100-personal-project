CREATE TABLE IF NOT EXISTS images (
    id SMALLSERIAL PRIMARY KEY,
    content_id INT NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    file_url VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255) NOT NULL,
    updated_by VARCHAR(255) NOT NULL,
    CONSTRAINT fk_content_id_images FOREIGN KEY (content_id) REFERENCES contents(id) ON DELETE CASCADE
);