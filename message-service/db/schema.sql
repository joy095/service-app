

First, create ENUM if not already done
CREATE TYPE connection_status AS ENUM (
    'pending', 'accepted', 'declined', 'blocked', 'withdrawn'
);

-- Now the fixed table with UUIDs
CREATE TABLE user_connections (
    requester_id UUID REFERENCES users(id) ON DELETE CASCADE,
    addressee_id UUID REFERENCES users(id) ON DELETE CASCADE,
    status connection_status NOT NULL DEFAULT 'pending',
    requested_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (requester_id, addressee_id),
    CHECK (requester_id <> addressee_id)
);


CREATE TABLE conversations (
    id SERIAL PRIMARY KEY,
    is_group BOOLEAN DEFAULT FALSE,
    title VARCHAR(255), -- optional for groups
    created_by UUID  REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE conversation_participants (
    conversation_id INTEGER REFERENCES conversations(id) ON DELETE CASCADE,
    user_id UUID  REFERENCES users(id) ON DELETE CASCADE,
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (conversation_id, user_id)
);

CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    conversation_id INTEGER REFERENCES conversations(id) ON DELETE CASCADE,
    sender_id UUID REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    reply_to_message_id INTEGER REFERENCES messages(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    edited_at TIMESTAMP
);