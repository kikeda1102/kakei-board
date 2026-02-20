CREATE TABLE expenses (
    id         VARCHAR(36)  NOT NULL,
    amount     BIGINT       NOT NULL,
    category   VARCHAR(64)  NOT NULL,
    memo       VARCHAR(512) NOT NULL DEFAULT '',
    date       DATE         NOT NULL,
    created_at DATETIME(6)  NOT NULL DEFAULT (UTC_TIMESTAMP(6)),
    PRIMARY KEY (id),
    INDEX idx_expenses_date (date DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
