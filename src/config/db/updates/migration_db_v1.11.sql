-- add SMTP configuration options
ALTER TABLE user ADD column deleted tinyint(1) DEFAULT '0' AFTER comment;

-- set new version in configuration
REPLACE INTO configuration (property, value) VALUES ("schema.version", "1.11");