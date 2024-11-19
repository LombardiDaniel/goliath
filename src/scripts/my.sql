START TRANSACTION;

-- users table
CREATE TABLE users ( 
    user_id INT AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    date_of_birth DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);
CREATE UNIQUE INDEX users_email_idx ON users (email);

-- MySQL does not use triggers for the "ON UPDATE" timestamp field like PostgreSQL
-- The 'updated_at' field is automatically handled by MySQL's `ON UPDATE CURRENT_TIMESTAMP`

-- unconfirmed_users table
CREATE TABLE unconfirmed_users ( 
    email VARCHAR(100) PRIMARY KEY,
    otp VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    date_of_birth DATE
);
CREATE UNIQUE INDEX unconfirmed_users_otp_idx ON unconfirmed_users (otp);

-- organizations table
CREATE TABLE organizations ( 
    organization_id CHAR(5) PRIMARY KEY,
    organization_name VARCHAR(100) NOT NULL,
    billing_plan_id INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP NULL,
    owner_user_id INT NOT NULL,
    FOREIGN KEY (owner_user_id) REFERENCES users (user_id)
);
CREATE UNIQUE INDEX organizations_organization_name_owner_user_id_constraint ON organizations (organization_name, owner_user_id);

-- organizations_users table
CREATE TABLE organizations_users (
    organization_id CHAR(5) NOT NULL,
    user_id INT NOT NULL,
    is_admin BOOLEAN NOT NULL DEFAULT false,
    PRIMARY KEY (organization_id, user_id),
    FOREIGN KEY (organization_id) REFERENCES organizations (organization_id),
    FOREIGN KEY (user_id) REFERENCES users (user_id)
);
CREATE UNIQUE INDEX organizations_users_idx ON organizations_users (organization_id, user_id);

-- organization_invites table
CREATE TABLE organization_invites (
    organization_id CHAR(5) NOT NULL,
    user_id INT NOT NULL,
    is_admin BOOLEAN NOT NULL DEFAULT false,
    invite_otp VARCHAR(255) NOT NULL,
    invite_exp TIMESTAMP,
    PRIMARY KEY (organization_id, user_id),
    FOREIGN KEY (organization_id) REFERENCES organizations (organization_id),
    FOREIGN KEY (user_id) REFERENCES users (user_id)
);

COMMIT;
