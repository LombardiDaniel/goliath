BEGIN;

CREATE TABLE users ( 
    user_id SERIAL PRIMARY KEY,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    date_of_birth DATE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    UNIQUE (user_id, email)
);

-- set user.is_updated trigger
CREATE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at := NOW();
    RETURN NEW;
END;
$$ LANGUAGE PLpgSQL;

CREATE TRIGGER update_updated_at_trigger
BEFORE UPDATE ON users
FOR EACH ROW EXECUTE PROCEDURE update_updated_at();

-- unconfirmedUsers
CREATE TABLE unconfirmed_users ( 
    email VARCHAR(100) PRIMARY KEY,
    otp VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    date_of_birth DATE
);

-- organizations
CREATE TABLE organizations ( 
    organization_id CHAR(5) PRIMARY KEY,
    organization_name VARCHAR(100) NOT NULL,
    billing_plan_id INT,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    deleted_at TIMESTAMP,
    owner_user_id INT REFERENCES users (user_id) NOT NULL ,

    UNIQUE (organization_name, owner_user_id)
);

-- join users orgs
CREATE TABLE organizations_users (
    organization_id CHAR(5) REFERENCES organizations (organization_id) NOT NULL,
    user_id INT REFERENCES users (user_id) NOT NULL,
    is_admin BOOLEAN NOT NULL DEFAULT false,

    PRIMARY KEY (organization_id, user_id)
);

-- org invites
CREATE TABLE organization_invites (
    organization_id CHAR(5) REFERENCES organizations (organization_id) NOT NULL,
    user_id INT REFERENCES users (user_id) NOT NULL,
    is_admin BOOLEAN NOT NULL DEFAULT false,
    otp VARCHAR(255) NOT NULL UNIQUE,
    exp TIMESTAMP
);

CREATE FUNCTION delete_expired_invites()
RETURNS TRIGGER AS $$
BEGIN
    DELETE FROM organization_invites
    WHERE exp < NOW();
    
    RETURN NEW;
END;
$$ LANGUAGE PLpgSQL;

CREATE TRIGGER delete_expired_org_invites
AFTER INSERT ON organization_invites
FOR EACH STATEMENT EXECUTE FUNCTION delete_expired_invites();

-- password_resets
CREATE TABLE password_resets (
    user_id INT REFERENCES users (user_id) NOT NULL,
    otp VARCHAR(255) NOT NULL UNIQUE,
    exp TIMESTAMP NOT NULL
);

CREATE FUNCTION delete_expired_resets()
RETURNS TRIGGER AS $$
BEGIN
    DELETE FROM password_resets
    WHERE exp < NOW();

    RETURN NEW;
END;
$$ LANGUAGE PLpgSQL;

CREATE TRIGGER delete_expired_password_resets
AFTER INSERT ON password_resets
FOR EACH STATEMENT EXECUTE FUNCTION delete_expired_resets();

-- oauth
CREATE TABLE oauth_users ( 
    email VARCHAR(100) PRIMARY KEY,
    user_id INT,
    oauth_provider VARCHAR(20) NOT NULL,
    CONSTRAINT fk_oauth_users FOREIGN KEY (user_id, email) REFERENCES users(user_id, email)
);

COMMIT;