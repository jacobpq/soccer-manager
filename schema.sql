CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL
);
DROP TABLE IF EXISTS sessions;
CREATE TABLE sessions (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    refresh_token VARCHAR(255) UNIQUE NOT NULL,
    refresh_expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE teams (
    id SERIAL PRIMARY KEY,
    user_id INT UNIQUE REFERENCES users(id),
    name VARCHAR(255) UNIQUE NOT NULL,
    country VARCHAR(100) NOT NULL,
    budget DECIMAL(15, 2) DEFAULT 5000000
);
CREATE TABLE players (
    id SERIAL PRIMARY KEY,
    team_id INT REFERENCES teams(id),
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    country VARCHAR(100),
    age INT,
    position VARCHAR(50),
    value DECIMAL(15, 2) DEFAULT 1000000,
    market_value DECIMAL(15, 2),
    on_transfer_list BOOLEAN DEFAULT FALSE
);