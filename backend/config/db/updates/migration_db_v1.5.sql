-- migrate configuration settings to the database
INSERT INTO configuration (property, value) VALUES
("browser.session.timeout.default", "10800"),
("browser.session.timeout.max",     "604800"),
("location.longitude",              "8.542939"),
("location.latitude",               "47.367658"),
("location.gmt_offset",             "1"),
("location.timezone",               "Europe/Berlin"),
("business.day.start",              "08:00:00"),
("business.day.end",                "21:00:00"),
("business.day.startatsunrise",     "false"),
("business.day.endatsunset",        "false"),
("session.cancel.graceperiod",      "86400");

-- update the new version in configuration
REPLACE INTO configuration (property, value) VALUES ("schema.version", "1.5");