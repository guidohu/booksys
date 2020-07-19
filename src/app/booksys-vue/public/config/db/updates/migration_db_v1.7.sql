-- add recaptcha keys
REPLACE INTO configuration (property, value) VALUES ("recaptcha.publickey", NULL);
REPLACE INTO configuration (property, value) VALUES ("recaptcha.privatekey", NULL);

-- set new version in configuration
REPLACE INTO configuration (property, value) VALUES ("schema.version", "1.7");