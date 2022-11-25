-- DROP TABLE IF EXISTS balance CASCADE;
-- DROP TABLE IF EXISTS reserve CASCADE;
-- DROP TABLE IF EXISTS history_operations CASCADE;

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
-- CREATE TABLE usr
-- (
--     id     BIGSERIAL PRIMARY KEY,
--     name   VARCHAR(64)              NOT NULL,
--     amount DECIMAL(13, 4) DEFAULT 0 NOT NULL
-- );
-- CREATE TABLE service
-- (
--     id    BIGSERIAL PRIMARY KEY,
--     name  VARCHAR(100)  NOT NULL,
--     price DECIMAL(10, 2) NOT NULL,
-- );
-- CREATE TABLE order
-- (
--     id BIGSERIAL PRIMARY KEY,
-- );
