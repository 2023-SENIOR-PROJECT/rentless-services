CREATE TABLE IF NOT EXISTS reviews (
	id INT(6) UNSIGNED AUTO_INCREMENT PRIMARY KEY,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    author_id INT(6) NOT NULL,
    product_id VARCHAR(255) NOT NULL,
    rate TINYINT DEFAULT NULL,
    comment VARCHAR(255) DEFAULT NULL
);