DELETE FROM cap
WHERE (fund, buyer, seller, units) IN (
    (1, 1, 0, 400),
    (1, 2, 0, 50),
    (1, 3, 0, 5)
);
