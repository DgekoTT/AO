CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create Users table
CREATE TABLE IF NOT EXISTS users (
    ID uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    Status VARCHAR(255) DEFAULT 'registered',
    RegistrationProvider VARCHAR(24) DEFAULT 'email',
    AccessToken TEXT,
    RefreshToken TEXT,
    Email VARCHAR(255) NOT NULL,
    HashedPassword BYTEA
    );

CREATE TABLE IF NOT EXISTS roles (
    RoleID serial PRIMARY KEY,
    Role VARCHAR(32) UNIQUE NOT NULL
    );

CREATE TABLE IF NOT EXISTS user_roles (
    UserID uuid,
    RoleID INT,
    PRIMARY KEY (UserID, RoleID),
    FOREIGN KEY (UserID) REFERENCES users(ID) ON DELETE CASCADE,
    FOREIGN KEY (RoleID) REFERENCES roles(RoleID) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS verify_infos (
    UserID uuid PRIMARY KEY,
    verificationToken TEXT,
    FOREIGN KEY (UserID) REFERENCES users (ID)
    );


CREATE TABLE IF NOT EXISTS oauth_providers (
    UserID uuid PRIMARY KEY,
    accountID VARCHAR(255),
    FOREIGN KEY (UserID) REFERENCES users (ID)
    );

CREATE INDEX IF NOT EXISTS idx_email ON users (Email);