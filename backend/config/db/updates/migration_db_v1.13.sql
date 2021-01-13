-- remove legacy fields in login
ALTER TABLE user DROP column password;
ALTER TABLE browser_session DROP column remember;

-- set new version in configuration
REPLACE INTO configuration (property, value) VALUES ("schema.version", "1.13");