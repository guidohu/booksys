-- add configuration option
REPLACE INTO configuration (property, value) VALUES ("engine.hour.format", "hh.h");

-- change data format for boat engine hours to support hh:mm more precisely
ALTER TABLE boat_engine_hours MODIFY before_hours DECIMAL(10,5) DEFAULT NULL;
ALTER TABLE boat_engine_hours MODIFY after_hours DECIMAL(10,5) DEFAULT NULL;
ALTER TABLE boat_engine_hours MODIFY delta_hours DECIMAL(10,5) DEFAULT NULL;

-- change data format for boat fuel hours to support hh:mm more precisely
ALTER TABLE boat_fuel MODIFY engine_hours DECIMAL(10,5) DEFAULT NULL;
ALTER TABLE boat_fuel MODIFY liters DECIMAL(10,3) DEFAULT NULL;
ALTER TABLE boat_fuel MODIFY cost_chf DECIMAL(10,3) DEFAULT NULL;
ALTER TABLE boat_fuel MODIFY cost_chf_brutto DECIMAL(10,3) DEFAULT NULL;

-- change data format for boat maintenance hours to support hh:mm more precisely
ALTER TABLE boat_maintenance MODIFY engine_hours DECIMAL(10,5) NOT NULL;

-- change tables to not use float
ALTER TABLE expenditure MODIFY amount_chf DECIMAL(10,3) DEFAULT NULL;
ALTER TABLE pricing MODIFY price_chf_min  DECIMAL(10,3) DEFAULT NULL;
ALTER TABLE payment MODIFY amount_chf     DECIMAL(10,3) DEFAULT NULL;
ALTER TABLE heat MODIFY cost_chf          DECIMAL(10,3) DEFAULT NULL;

-- remove table boat_oil
DROP TABLE IF EXISTS boat_oil;

-- set new version in configuration
REPLACE INTO configuration (property, value) VALUES ("schema.version", "1.15");