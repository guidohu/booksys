-- add configuration option
REPLACE INTO configuration (property, value) VALUES ("fuel.payment.type", "instant");

-- add new column to directly pay for fuel or pay via bills
ALTER TABLE boat_fuel ADD COLUMN contributes_to_balance int(8) DEFAULT 1;
INSERT INTO `expenditure_type` (`id`, `name`, `comment`) 
    VALUES (9, 'fuel bill', 'Costs for fuel through bill');

-- set new version in configuration
REPLACE INTO configuration (property, value) VALUES ("schema.version", "1.16");