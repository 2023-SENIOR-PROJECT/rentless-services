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
CREATE TABLE IF NOT EXISTS reviews (
	id INT(6) UNSIGNED AUTO_INCREMENT PRIMARY KEY,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    author_id INT(6) NOT NULL,
    product_id VARCHAR(255) NOT NULL,
    rate TINYINT DEFAULT NULL,
    comment VARCHAR(255) DEFAULT NULL
);

-- Reviews for Casio FX-115ESPLUS Scientific Calculator
INSERT INTO reviews (created_at, author_id, product_id, rate, comment) VALUES
    (CURRENT_TIMESTAMP, 1, '652f72d76659132b2d8647d6', 4, 'Great calculator!'),
    (CURRENT_TIMESTAMP, 2, '652f72d86659132b2d8647d8', 5, 'Excellent features and functionality.'),
    (CURRENT_TIMESTAMP, 3, '652f72d96659132b2d8647da', 3, 'Decent calculator, but a bit complex to use.'),
    (CURRENT_TIMESTAMP, 4, '6536a96365ac4edd1ef8c031', 4, 'Good value for money.');

-- Reviews for Premium Wireless Headphones
INSERT INTO reviews (created_at, author_id, product_id, rate, comment) VALUES
    (CURRENT_TIMESTAMP, 1, '654ced974fd690fff6d4aea5', 5, 'Amazing sound quality!'),
    (CURRENT_TIMESTAMP, 2, '654cee2b3048ff2d6692f2a8', 4, 'Comfortable to wear for extended periods.'),
    (CURRENT_TIMESTAMP, 3, '654cee2e3048ff2d6692f2aa', 3, 'Average performance, not worth the price.'),
    (CURRENT_TIMESTAMP, 4, '654d021796f1773e5c2350b8', 5, 'Best headphones I have ever owned!'),
    (CURRENT_TIMESTAMP, 1, '654d11df588f8e387e2edf7e', 4, 'Stylish design and good battery life.');

-- Reviews for Casio FX-991EX Scientific Calculator
INSERT INTO reviews (created_at, author_id, product_id, rate, comment) VALUES
    (CURRENT_TIMESTAMP, 2, '65538a67614e18ad5c04da94', 5, 'A must-have for engineering students.'),
    (CURRENT_TIMESTAMP, 3, '65538f82614e18ad5c04da98', 4, 'Solid performance and reliable calculations.'),
    (CURRENT_TIMESTAMP, 4, '65538a67614e18ad5c04da94', 3, 'Difficult to navigate, but gets the job done.');



-- COPY PatientInfo FROM '/home/data/PatientInfo.csv' DELIMITER ',' CSV HEADER;
-- COPY Datasource FROM '/home/data/datasource.csv' DELIMITER ',' CSV HEADER;
-- COPY province FROM '/home/data/province.csv' DELIMITER ',' CSV HEADER;