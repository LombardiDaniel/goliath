BEGIN;

CREATE TABLE users ( 
    user_id SERIAL PRIMARY KEY,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    date_of_birth DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);
CREATE UNIQUE INDEX users_email_idx ON users (email);

-- set user.is_updated trigger
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at := CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_updated_at_trigger
BEFORE UPDATE ON users
FOR EACH ROW EXECUTE PROCEDURE update_updated_at();

-- unconfirmedUsers
CREATE TABLE unconfirmed_users ( 
    email VARCHAR(100) PRIMARY KEY,
    otp VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    date_of_birth DATE
);
CREATE UNIQUE INDEX unconfirmed_users_otp_idx ON unconfirmed_users (otp);

-- organizations
CREATE TABLE organizations ( 
    organization_id CHAR(5) PRIMARY KEY,
    organization_name VARCHAR(100) NOT NULL,
    billing_plan_id INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,
    owner_user_id INT NOT NULL,
    FOREIGN KEY (owner_user_id) REFERENCES users (user_id)
);
CREATE UNIQUE INDEX organizations_organization_nae_owner_user_id_constraint ON organizations (organization_name, owner_user_id);

-- join users orgs
CREATE TABLE organizations_users (
    organization_id CHAR(5) REFERENCES organizations (organization_id) NOT NULL,
    user_id INT REFERENCES users (user_id) NOT NULL,
    is_admin BOOLEAN NOT NULL DEFAULT false
);
CREATE UNIQUE INDEX organizations_users_idx ON organizations_users (organization_id, user_id);

-- org invites
CREATE TABLE organization_invites (
    organization_id CHAR(5) REFERENCES organizations (organization_id) NOT NULL,
    user_id INT REFERENCES users (user_id) NOT NULL,
    is_admin BOOLEAN NOT NULL DEFAULT false,
    invite_otp VARCHAR(255) NOT NULL,
    invite_exp TIMESTAMP
);

COMMIT;