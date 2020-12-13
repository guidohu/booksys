--
-- Remove all tables if they exist
--

DROP TABLE IF EXISTS `boat_engine_hours`;
DROP TABLE IF EXISTS `boat_fuel`;
DROP TABLE IF EXISTS `boat_maintenance`;
DROP TABLE IF EXISTS `boat_oil`;
DROP TABLE IF EXISTS `browser_session`;
DROP TABLE IF EXISTS `expenditure`;
DROP TABLE IF EXISTS `expenditure_type`;
DROP TABLE IF EXISTS `heat`;
DROP TABLE IF EXISTS `invitation`;
DROP TABLE IF EXISTS `invitation_status`;
DROP TABLE IF EXISTS `password_reset`;
DROP TABLE IF EXISTS `payment`;
DROP TABLE IF EXISTS `pricing`;
DROP TABLE IF EXISTS `session`;
DROP TABLE IF EXISTS `session_type`;
DROP TABLE IF EXISTS `user`;
DROP TABLE IF EXISTS `user_role`;
DROP TABLE IF EXISTS `user_status`;
DROP TABLE IF EXISTS `user_to_session`;
DROP TABLE IF EXISTS `configuration`;

--
-- Table structure for table `boat_engine_hours`
--

CREATE TABLE `boat_engine_hours` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `timestamp` datetime DEFAULT NULL,
  `before_hours` float(10,1) DEFAULT NULL,
  `after_hours` float(10,1) DEFAULT NULL,
  `delta_hours` float(10,1) DEFAULT NULL,
  `type` int(11) DEFAULT NULL,
  `user_id` int(11) DEFAULT NULL,
  `comment` text CHARACTER SET utf8,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

--
-- Table structure for table `boat_fuel`
--

CREATE TABLE `boat_fuel` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `timestamp` datetime DEFAULT NULL,
  `user_id` int(11) DEFAULT NULL,
  `engine_hours` float(10,1) DEFAULT NULL,
  `liters` float(10,2) DEFAULT NULL,
  `cost_chf` float(10,2) DEFAULT NULL,
  `cost_chf_brutto` float(10,2) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

--
-- Table structure for table `boat_maintenance`
--

CREATE TABLE `boat_maintenance` (
  `id` mediumint(9) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `engine_hours` float(10,1) NOT NULL,
  `description` text CHARACTER SET utf8,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

--
-- Table structure for table `boat_oil`
--

CREATE TABLE `boat_oil` (
  `id` mediumint(9) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `oil_l` float(10,2) DEFAULT NULL,
  `engine_hours` float(10,2) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

--
-- Table structure for table `browser_session`
--

CREATE TABLE `browser_session` (
  `session_secret` char(64) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `valid_thru` datetime DEFAULT NULL,
  `last_activity` datetime DEFAULT NULL,
  `user_id` int(11) DEFAULT NULL,
  `username` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `first_name` varchar(255) CHARACTER SET utf8 DEFAULT NULL,
  `last_name` varchar(255) CHARACTER SET utf8 DEFAULT NULL,
  `user_status` int(11) DEFAULT NULL,
  `user_role_id` int(11) NOT NULL,
  `session_data` text CHARACTER SET utf8,
  PRIMARY KEY (`session_secret`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

--
-- Table structure for table `expenditure`
--

CREATE TABLE `expenditure` (
  `id` mediumint(9) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) DEFAULT NULL,
  `type_id` int(11) DEFAULT NULL,
  `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `amount_chf` float(10,2) DEFAULT NULL,
  `comment` text CHARACTER SET utf8,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

--
-- Table structure for table `expenditure_type`
--

CREATE TABLE `expenditure_type` (
  `id` int(11) NOT NULL DEFAULT '0',
  `name` text CHARACTER SET utf8,
  `comment` text CHARACTER SET utf8,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

INSERT INTO `expenditure_type` (`id`, `name`, `comment`) VALUES
(0, 'fuel', 'Costs for fuel'),
(1, 'maintenance', 'Costs for maintenance'),
(2, 'material', 'Material which is needed for the boat'),
(3, 'invest', 'Investment in new Features, ...'),
(4, 'session', 'Payments for sessions'),
(5, 'other', 'Other costs'),
(6, 'membership fee', 'Membership Fee'),
(7, 'salary', 'money for towing people or do whatever in the name of the community'),
(8, 'owners refund', 'Payback of investments from owners');

--
-- Table structure for table `heat`
--

CREATE TABLE `heat` (
  `id` mediumint(9) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) DEFAULT NULL,
  `session_id` int(11) DEFAULT NULL,
  `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `duration_s` int(11) DEFAULT NULL,
  `cost_chf` float(10,2) DEFAULT NULL,
  `comment` TEXT DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

--
-- Table structure for table `invitation`
--

CREATE TABLE `invitation` (
  `id` mediumint(9) NOT NULL AUTO_INCREMENT,
  `session_id` int(11) DEFAULT NULL,
  `user_id` int(11) DEFAULT NULL,
  `inviter_id` int(11) DEFAULT NULL,
  `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `status` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

--
-- Table structure for table `invitation_status`
--

CREATE TABLE `invitation_status` (
  `id` int(11) NOT NULL DEFAULT '0',
  `name` text CHARACTER SET utf8,
  `comment` text CHARACTER SET utf8,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

INSERT INTO `invitation_status` (`id`, `name`, `comment`) VALUES
(1, 'full', 'Session is full by now'),
(2, 'query', 'The invitation is still valid'),
(3, 'accepted', 'User accepted the invitation'),
(4, 'declined', 'User declined the invitation'),
(5, 'unknown', 'No invitation sent yet');

--
-- Table structure for table `password_reset`
--

CREATE TABLE `password_reset` (
  `id` mediumint(9) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `token` varchar(255) NOT NULL,
  `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `valid` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

--
-- Table structure for table `payment`
--

CREATE TABLE `payment` (
  `id` mediumint(9) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) DEFAULT NULL,
  `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `amount_chf` float(10,2) DEFAULT NULL,
  `type_id` int(11) DEFAULT NULL,
  `comment` text CHARACTER SET utf8,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

--
-- Table structure for table `pricing`
--

CREATE TABLE `pricing` (
  `id` mediumint(9) NOT NULL AUTO_INCREMENT,
  `user_status_id` int(11) DEFAULT NULL,
  `price_chf_min` float(10,2) DEFAULT NULL,
  `comment` text CHARACTER SET utf8,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

INSERT INTO `pricing` (`id`, `user_status_id`, `price_chf_min`, `comment`) VALUES
(1, 1, 2.80, 'Guest Price'),
(2, 2, 2.80, 'Member Price'),
(3, 3, 1.30, 'Boat Community Price'),
(4, 4, 0.00, 'Course Price'),
(5, 5, 1.30, 'Partner of Boat Community');

--
-- Table structure for table `session`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `session` (
  `id` mediumint(9) NOT NULL AUTO_INCREMENT,
  `date` date DEFAULT NULL,
  `start` time DEFAULT NULL,
  `start_time` timestamp DEFAULT CURRENT_TIMESTAMP,
  `end` time DEFAULT NULL,
  `end_time` timestamp DEFAULT CURRENT_TIMESTAMP,
  `title` text CHARACTER SET utf8,
  `comment` text CHARACTER SET utf8,
  `type` int(11) DEFAULT '0',
  `free` int(11) DEFAULT NULL,
  `creator_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

--
-- Table structure for table `session_type`
--

CREATE TABLE `session_type` (
  `id` mediumint(9) NOT NULL,
  `name` text CHARACTER SET utf8,
  `comment` text CHARACTER SET utf8,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

INSERT INTO `session_type` (`id`, `name`, `comment`) VALUES
(0, 'default', 'A normal session'),
(1, 'course', 'A course session');

--
-- Table structure for table `user`
--

CREATE TABLE `user` (
  `id` mediumint(9) NOT NULL AUTO_INCREMENT,
  `username` varchar(255) DEFAULT NULL,
  `password_salt` int(11) DEFAULT NULL,
  `password_hash` text CHARACTER SET utf8,
  `first_name` varchar(255) CHARACTER SET utf8 DEFAULT NULL,
  `last_name` varchar(255) CHARACTER SET utf8 DEFAULT NULL,
  `address` varchar(255) CHARACTER SET utf8 DEFAULT NULL,
  `city` varchar(255) CHARACTER SET utf8 DEFAULT NULL,
  `plz` int(11) DEFAULT NULL,
  `mobile` varchar(255) DEFAULT NULL,
  `email` varchar(255) DEFAULT NULL,
  `license` tinyint(1) DEFAULT '0',
  `status` int(11) DEFAULT '0',
  `locked` tinyint(1) DEFAULT '1',
  `comment` text CHARACTER SET utf8,
  `deleted` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

--
-- Table structure for table `user_role`
--

CREATE TABLE `user_role` (
  `id` int(11) NOT NULL,
  `name` text COLLATE utf8_bin NOT NULL,
  `description` text COLLATE utf8_bin NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

INSERT INTO `user_role` (`id`, `name`, `description`) VALUES
(1, 'guest', 'Guest User Permissions'),
(2, 'member', 'Member User Permissions'),
(3, 'admin', 'Administrator User Permissions');

--
-- Table structure for table `user_status`
--

CREATE TABLE `user_status` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` text CHARACTER SET utf8,
  `description` text CHARACTER SET utf8,
  `user_role_id` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

INSERT INTO `user_status` (`id`, `name`, `description`, `user_role_id`) VALUES
(1, 'Guest', 'Guest Access', 1),
(2, 'Member', 'Member Access', 2),
(3, 'Admin', 'Administrator and Boat Community Access', 3),
(4, 'Course', 'Course Status', 1),
(5, 'Partner', 'Member Permissions', 2);

--
-- Table structure for table `user_to_session`
--

CREATE TABLE `user_to_session` (
  `id` mediumint(9) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) DEFAULT NULL,
  `session_id` int(11) DEFAULT NULL,
  `time` time DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

--
-- Table structure for table `configuration`
--

CREATE TABLE `configuration` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `property` VARCHAR(255) NOT NULL UNIQUE,
  `value` TEXT,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

INSERT INTO configuration (property, value) VALUES 
("schema.version",                  "1.9"           ),
("browser.session.timeout.default", "10800"         ),
("browser.session.timeout.max",     "604800"        ),
("location.longitude",              "8.542939"      ),
("location.latitude",               "47.367658"     ),
("location.gmt_offset",             "1"             ),
("location.timezone",               "Europe/Berlin" ),
("business.day.start",              "08:00:00"      ),
("business.day.end",                "21:00:00"      ),
("business.day.startatsunrise",     "false"         ),
("business.day.endatsunset",        "false"         ),
("session.cancel.graceperiod",      "86400"         ),
("recaptcha.privatekey",            NULL            ),
("recaptcha.publickey",             NULL            ),
("currency",                        "CHF"           ),
("location.address",                NULL            ),
("location.map",                    NULL            ),
("payment.account.owner",           NULL            ),
("payment.account.iban",            NULL            ),
("payment.account.bic",             NULL            ),
("payment.account.comment",         NULL            ),
("smtp.sender",                     NULL            ),
("smtp.server",                     NULL            ),
("smtp.username",                   NULL            ),
("smtp.password",                   NULL            );
