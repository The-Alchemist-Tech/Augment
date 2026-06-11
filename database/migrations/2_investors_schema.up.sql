CREATE TABLE `investors` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `username` VARCHAR(32) NOT NULL DEFAULT "",
  `email` VARCHAR(255) NOT NULL DEFAULT "",
  `first_name` VARCHAR(32) NOT NULL DEFAULT "",
  `last_name` VARCHAR(32) NOT NULL DEFAULT "",
  `created_on` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `last_modified` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Index on queried fields
CREATE INDEX idx_investors_id ON investors(id);
CREATE INDEX idx_investors_email ON investors(email);