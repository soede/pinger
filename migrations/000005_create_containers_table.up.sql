CREATE TABLE containers (
    container_id SERIAL PRIMARY KEY,
    container_status VARCHAR(255) NOT NULL,
    addr VARCHAR(255) NOT NULL,
    p_duration FLOAT NOT NULL,
    pinged_at TIMESTAMP NOT NULL
);

