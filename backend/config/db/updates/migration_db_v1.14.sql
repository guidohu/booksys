-- add configuration
REPLACE INTO configuration (property, value) VALUES ("logo.file", NULL);

-- set new version in configuration
REPLACE INTO configuration (property, value) VALUES ("schema.version", "1.14");