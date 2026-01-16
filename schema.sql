CREATE TABLE point_root (
    user_id VARCHAR(20) PRIMARY KEY,
    point_num INT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- 予約テーブル
CREATE TABLE point_reservations (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(20) NOT NULL,
    point_amount INT NOT NULL,
    execute_at TIMESTAMP NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    idempotency_key VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    INDEX idx_execute_status (execute_at, status)
);

-- 冪等性確保用トランザクションテーブル
CREATE TABLE point_transactions (
    id VARCHAR(36) PRIMARY KEY,
    idempotency_key VARCHAR(100) UNIQUE NOT NULL,
    user_id VARCHAR(20) NOT NULL,
    point_amount INT NOT NULL,
    created_at TIMESTAMP NOT NULL
);
