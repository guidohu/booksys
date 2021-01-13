-- add SMTP configuration options
REPLACE INTO configuration (property, value) VALUES ("smtp.sender", NULL);
REPLACE INTO configuration (property, value) VALUES ("smtp.server", NULL);
REPLACE INTO configuration (property, value) VALUES ("smtp.username", NULL);
REPLACE INTO configuration (property, value) VALUES ("smtp.password", NULL);

-- set new version in configuration
REPLACE INTO configuration (property, value) VALUES ("schema.version", "1.9");