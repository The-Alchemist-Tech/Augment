CREATE TABLE `funds` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(128) NOT NULL DEFAULT "",
  `units` int(10) NOT NULL DEFAULT 0
  `created_on` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `last_modified` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Index on queried fields
CREATE INDEX idx_funds_id ON funds(id);

CREATE INDEX idx_funds_name ON funds(name);