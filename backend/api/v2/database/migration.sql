-- Collection of commands that were needed to make the database setup work with the new code
ALTER TABLE browser_session MODIFY session_secret varchar(256) COLLATE utf8_unicode_ci NOT NULL;

-- change session types

-- set is_discounted in case cost_chf_brutto is no NULL and cost values are not the same