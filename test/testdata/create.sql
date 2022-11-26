CREATE TABLE IF NOT EXISTS balance
(
    id      BIGSERIAL PRIMARY KEY,
    user_id BIGINT                   NOT NULL,
    amount  DECIMAL(13, 4) DEFAULT 0 NOT NULL
);
CREATE TABLE IF NOT EXISTS reserve
(
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT                NOT NULL,
    service_id BIGINT                NOT NULL,
    order_id   BIGINT                NOT NULL,
    cost       DECIMAL(13, 4)        NOT NULL,
    is_profit  BOOLEAN DEFAULT false NOT NULL
);
CREATE TABLE IF NOT EXISTS history_operations
(
    id             BIGSERIAL PRIMARY KEY,
    sender_user_id BIGINT                                       NULL,
    user_id        BIGINT                                       NOT NULL,
    service_id     BIGINT                                       NULL,
    amount         DECIMAL(13, 4)                               NOT NULL,
    type           VARCHAR(32)                                  NOT NULL,
    datetime       TIMESTAMP DEFAULT (now() AT TIME ZONE 'utc') NOT NULL
);