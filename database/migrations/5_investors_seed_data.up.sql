INSERT INTO investors (id, username, email, first_name, last_name, created_on, last_modified)
VALUES
(1, "fund", "fund@fund.com", "fund", "fund", NOW(), NOW()), -- acts as a way for initial investors have shares by having this as the seller
(2, "testUser1", "test1@test.com", "testFirst1", "testLast1", NOW(), NOW()),
(3, "testUser2", "test2@test.com", "testFirst2", "testLast2", NOW(), NOW()),
(4, "testUser3", "test3@test.com", "testFirst3", "testLast3", NOW(), NOW());