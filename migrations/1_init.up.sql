CREATE TABLE IF NOT EXISTS transactions
(
    id BIGSERIAL PRIMARY KEY,
    external_id VARCHAR(64) NOT NULL UNIQUE,
    block_id BIGINT NOT NULL,
    timestamp TIMESTAMP NOT NULL
);

CREATE TABLE delegate_messages
(
    id BIGSERIAL PRIMARY KEY,
    transaction_id BIGINT NOT NULL REFERENCES transactions(id),
    delegator_address VARCHAR(64) NOT NULL,
    validator_address VARCHAR(64) NOT NULL,
    amount BIGINT NOT NULL,
    currency VARCHAR(64) NOT NULL
);

