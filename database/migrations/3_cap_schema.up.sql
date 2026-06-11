CREATE TABLE `cap` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `fund` int(11) unsigned NOT NULL,
  `buyer` int(11) unsigned NOT NULL,
  `seller` int(11) unsigned NOT NULL,
  `units` DECIMAL(14,4) NOT NULL,
  `created_on` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `last_modified` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),

  -- Define Foreign Keys
  CONSTRAINT fk_cap_fund 
        FOREIGN KEY (fund) 
        REFERENCES funds(id),

  CONSTRAINT fk_cap_buyer 
        FOREIGN KEY (buyer) 
        REFERENCES investors(id),

  CONSTRAINT fk_cap_seller 
        FOREIGN KEY (seller) 
        REFERENCES investors(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Index on queried fields
-- I would add indexes on fund, buyer, and seller as these are
-- likely to be used in queries in a more full featured app
CREATE INDEX idx_cap_fund ON cap(fund);