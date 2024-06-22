-- Create items table
CREATE TABLE IF NOT EXISTS items (
    itemID          BIGSERIAL    UNIQUE NOT NULL,
    categoryID      BIGINT       NOT NULL,
    itemName        VARCHAR(128)  NOT NULL,
    itemDescription VARCHAR(255),
    price           MONEY        NOT NULL,
    archived        BOOLEAN     DEFAULT(FALSE),

    created_at  TIMESTAMP  DEFAULT now(),
    
    PRIMARY KEY (itemID)
);

-- Create stores table
CREATE TABLE IF NOT EXISTS stores (
    storeID     BIGSERIAL    UNIQUE NOT NULL,
    storeName   VARCHAR(255) NOT NULL,

    created_at  TIMESTAMP  DEFAULT now(),
    
    PRIMARY KEY (storeID)
);

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    userID      BIGSERIAL    UNIQUE NOT NULL,
    firstName   VARCHAR(64),
    middleName  VARCHAR(64),
    lastName    VARCHAR(64),
    birthDate   DATE,

    created_at  TIMESTAMP  DEFAULT now(),

    PRIMARY KEY (userID)
);

-- Create payments table
CREATE TABLE IF NOT EXISTS payments (
    paymentID   BIGSERIAL    UNIQUE NOT NULL,
    cardholder  VARCHAR(255) NOT NULL,
    cardnumber  VARCHAR(16)  NOT NULL,
    expiration  DATE         NOT NULL,

    created_at  TIMESTAMP  DEFAULT now(),

    PRIMARY KEY (paymentID)
);

-- Create POS table
CREATE TABLE IF NOT EXISTS pos (
    posID    BIGSERIAL UNIQUE NOT NULL,
    storeID  BIGINT    NOT NULL,
    hostname VARCHAR(32),

    created_at  TIMESTAMP  DEFAULT now(),

    PRIMARY KEY (posID),
    FOREIGN KEY (storeID) REFERENCES stores(storeID)
        ON UPDATE CASCADE
        ON DELETE CASCADE
);

-- Create transactions table
CREATE TABLE IF NOT EXISTS transactions (
    transactionID BIGSERIAL UNIQUE NOT NULL,
    status        VARCHAR(16) NOT NULL,
    posID         BIGINT      NOT NULL,
    storeID       BIGINT,
    userID        BIGINT,
    total         MONEY       DEFAULT 0,
    paymentID     BIGINT      DEFAULT 0,
    archived      BOOLEAN     DEFAULT(FALSE),
    startTime     TIMESTAMP,
    endTime       TIMESTAMP,

    created_at  TIMESTAMP  DEFAULT now(),

    PRIMARY KEY (transactionID),
    FOREIGN KEY (posID) REFERENCES pos(posID)
        ON UPDATE CASCADE
        ON DELETE SET NULL,
    FOREIGN KEY (storeID) REFERENCES stores(storeID)
        ON UPDATE CASCADE
        ON DELETE SET NULL,
    FOREIGN KEY (userID) REFERENCES users(userID)
        ON UPDATE CASCADE
        ON DELETE NO ACTION
    
    -- Removed to prevent issues with insertion
    -- FOREIGN KEY (paymentID) REFERENCES payments(paymentID)
    --     ON UPDATE CASCADE
    --     ON DELETE SET NULL
);

-- Create transaction items table
CREATE TABLE IF NOT EXISTS transaction_items (
    entryID       BIGSERIAL NOT NULL,
    transactionID BIGINT NOT NULL,
    itemID        BIGINT NOT NULL,
    quantity      INTEGER,

    created_at TIMESTAMP DEFAULT now(),

    PRIMARY KEY (entryID),
    FOREIGN KEY (transactionID) REFERENCES transactions(transactionID)
        ON UPDATE CASCADE
        ON DELETE CASCADE,
    FOREIGN KEY (itemID) REFERENCES items(itemID)
        ON UPDATE CASCADE
        ON DELETE CASCADE
);