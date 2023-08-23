/*
    Seed user table
*/
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL,
    pwd VARCHAR(32) NOT NULL
);

/*
    Seed user data 
*/
INSERT INTO users(username, email, pwd) VALUES
    ('alice', 'alice@mail.com', '6384e2b2184bcbf58eccf10ca7a6563c'),
    ('bob', 'bob@mail.com', '9f9d51bc70ef21ca5c14f307980a29d8');