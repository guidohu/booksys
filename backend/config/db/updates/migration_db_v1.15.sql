-- add configuration option
REPLACE INTO configuration (property, value) VALUES ("engine.hour.format", "hh.h");

-- set new version in configuration
REPLACE INTO configuration (property, value) VALUES ("schema.version", "1.15");