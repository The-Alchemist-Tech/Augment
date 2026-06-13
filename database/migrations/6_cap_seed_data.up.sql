INSERT INTO cap (fund, buyer, seller, units, created_on)
VALUES
    (1, 2, 1, 400, NOW()), -- Units were transferred from the fund user to get some to the users.
    (1, 3, 1, 50, NOW()),
    (1, 4, 1, 5, NOW());
