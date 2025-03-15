package db

import (
    "context"
    "fmt"
    "log"
)

func CreateUsersTable() error {
    // Create extensions
    extensionsQuery := `
    CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
    CREATE EXTENSION IF NOT EXISTS "citext";`

    if _, err := DB.Exec(context.Background(), extensionsQuery); err != nil {
        return fmt.Errorf("error creating extensions: %v", err)
    }

    // Create users table with timestamp columns
    tableQuery := `
    CREATE TABLE IF NOT EXISTS users (
        id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
        username CITEXT UNIQUE NOT NULL,
        email CITEXT UNIQUE NOT NULL,
        password_hash TEXT NOT NULL,
        refresh_token TEXT,
        otp_hash TEXT,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );`

    if _, err := DB.Exec(context.Background(), tableQuery); err != nil {
        return fmt.Errorf("error creating users table: %v", err)
    }

    // Create update timestamp function
    functionQuery := `
    CREATE OR REPLACE FUNCTION trigger_set_updated_at()
    RETURNS TRIGGER AS $$
    BEGIN
        NEW.updated_at = CURRENT_TIMESTAMP;
        RETURN NEW;
    END;
    $$ LANGUAGE plpgsql;`

    if _, err := DB.Exec(context.Background(), functionQuery); err != nil {
        return fmt.Errorf("error creating update timestamp function: %v", err)
    }

    // Drop existing trigger if any
    dropTriggerQuery := `DROP TRIGGER IF EXISTS trigger_set_updated_at ON users;`
    if _, err := DB.Exec(context.Background(), dropTriggerQuery); err != nil {
        return fmt.Errorf("error dropping existing trigger: %v", err)
    }

    // Create new trigger
    triggerQuery := `
    CREATE TRIGGER trigger_set_updated_at
        BEFORE UPDATE ON users
        FOR EACH ROW
        EXECUTE FUNCTION trigger_set_updated_at();`

    if _, err := DB.Exec(context.Background(), triggerQuery); err != nil {
        return fmt.Errorf("error creating trigger: %v", err)
    }

    log.Println("Database migration completed successfully!")
    return nil
}