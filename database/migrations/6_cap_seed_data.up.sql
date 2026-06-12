INSERT INTO cap (fund, buyer, seller, units, created_on, last_modified)
VALUES
    (1, 2, 1, 400, NOW(), NOW()), -- Units were transferred from the fund user to get some to the users.
    (1, 3, 1, 50, NOW(), NOW()),
    (1, 4, 1, 5, NOW(), NOW());
