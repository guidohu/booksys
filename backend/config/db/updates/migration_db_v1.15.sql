-- add configuration option
REPLACE INTO configuration (property, value) VALUES ("engine.hour.format", "hh.h");

-- change data format for boat engine hours to support hh:mm more precisely
ALTER TABLE boat_engine_hours MODIFY before_hours DECIMAL(10,5) DEFAULT NULL;
ALTER TABLE boat_engine_hours MODIFY after_hours DECIMAL(10,5) DEFAULT NULL;
ALTER TABLE boat_engine_hours MODIFY delta_hours DECIMAL(10,5) DEFAULT NULL;

-- set new version in configuration
REPLACE INTO configuration (property, value) VALUES ("schema.version", "1.15");