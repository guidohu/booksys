-- add comment fields to heats
ALTER TABLE heat ADD column comment TEXT AFTER cost_chf;

-- set new version in configuration
REPLACE INTO configuration (property, value) VALUES ("schema.version", "1.6");