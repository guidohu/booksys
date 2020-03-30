-- add recaptcha keys
REPLACE INTO configuration (property, value) VALUES ("currency", "CHF");
REPLACE INTO configuration (property, value) VALUES ("location.map", NULL);
REPLACE INTO configuration (property, value) VALUES ("location.address", NULL);
REPLACE INTO configuration (property, value) VALUES ("payment.account.owner", NULL);
REPLACE INTO configuration (property, value) VALUES ("payment.account.iban",  NULL);
REPLACE INTO configuration (property, value) VALUES ("payment.account.bic", NULL);
REPLACE INTO configuration (property, value) VALUES ("payment.account.comment", NULL);

-- set new version in configuration
REPLACE INTO configuration (property, value) VALUES ("schema.version", "1.8");