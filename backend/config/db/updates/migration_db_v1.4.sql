-- remove friendship status tables
DROP table IF EXISTS friend;
DROP table IF EXISTS friendship_status;

-- add save timefields to sessions
ALTER TABLE session ADD column start_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP AFTER start;
ALTER TABLE session ADD column end_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP AFTER end;
UPDATE session SET start_time = CONCAT(date, ' ', start) WHERE start IS NOT NULL;
UPDATE session SET start_time = DATE_SUB(start_time, INTERVAL 1 HOUR) WHERE start IS NOT NULL;
UPDATE session SET end_time = CONCAT(date, ' ', end) WHERE start IS NOT NULL;
UPDATE session SET end_time = DATE_SUB(end_time, INTERVAL 1 HOUR) WHERE start IS NOT NULL;

-- create configuration table
CREATE TABLE configuration (
    id INT AUTO_INCREMENT PRIMARY KEY,
    property VARCHAR(255) NOT NULL UNIQUE,
    value TEXT
);

-- set new version in configuration
INSERT INTO configuration (property, value) VALUES ("schema.version", "1.4");