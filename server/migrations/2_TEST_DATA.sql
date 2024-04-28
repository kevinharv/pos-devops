-- Create some items
INSERT INTO items (categoryID, itemName, itemDescription, price) VALUES (1, 'Dell XPS 13', 'Laptop Computer', 1499.99);

-- Create some stores
INSERT INTO stores (storeName) VALUES ('Store 1');
INSERT INTO stores (storeName) VALUES ('Store 2');
INSERT INTO stores (storeName) VALUES ('Store 3');

-- Create some users
INSERT INTO users (firstName, middleName, lastName, birthDate) VALUES ('Joseph', 'Lee', 'Lin', '2001/07/07');

-- Create some payments
INSERT INTO payments (cardholder, cardnumber, expiration) VALUES ('Foo A Bar', '123456789098765', '2026/04/01');

-- Create some POS terminals
INSERT INTO pos (storeID, hostname) VALUES (1, 'S1POS1');
INSERT INTO pos (storeID, hostname) VALUES (1, 'S1POS2');
INSERT INTO pos (storeID, hostname) VALUES (2, 'S2POS1');
INSERT INTO pos (storeID, hostname) VALUES (2, 'S2POS2');
INSERT INTO pos (storeID, hostname) VALUES (3, 'S3POS1');
INSERT INTO pos (storeID, hostname) VALUES (3, 'S3POS2');