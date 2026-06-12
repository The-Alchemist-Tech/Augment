DELETE FROM investors
WHERE (username, email, first_name, last_name) IN
(
    ("testUser1", "test1@test.com", "testFirst1", "testLast1"),
    ("testUser2", "test2@test.com", "testFirst2", "testLast2")
    ("testUser3", "test3@test.com", "testFirst3", "testLast3")
);