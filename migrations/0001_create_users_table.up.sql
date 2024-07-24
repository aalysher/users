CREATE TABLE users (
                       id VARCHAR(255) PRIMARY KEY,
                       first_name VARCHAR(255) NOT NULL,
                       last_name VARCHAR(255) NOT NULL,
                       age INTEGER NOT NULL,
                       email VARCHAR(255) UNIQUE NOT NULL,
                       created TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);