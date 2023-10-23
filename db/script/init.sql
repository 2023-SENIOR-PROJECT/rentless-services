CREATE TABLE IF NOT EXISTS users (
	id INT(6) UNSIGNED AUTO_INCREMENT PRIMARY KEY,
	firstname VARCHAR(30) NOT NULL,
	lastname VARCHAR(30) NOT NULL,
	age INT(3) NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS auths (
    email VARCHAR(50) PRIMARY KEY,
    pwd VARCHAR(255) NOT NULL,
    token VARCHAR(255) NOT NULL,
    user_id INT(6) UNSIGNED NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- COPY PatientInfo FROM '/home/data/PatientInfo.csv' DELIMITER ',' CSV HEADER;
-- COPY Datasource FROM '/home/data/datasource.csv' DELIMITER ',' CSV HEADER;
-- COPY province FROM '/home/data/province.csv' DELIMITER ',' CSV HEADER;