-- 1. 建立 users 資料表
CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    registration_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_visit_ip VARCHAR(45),
    role VARCHAR(20) NOT NULL CHECK (role IN ('NULL', 'USERS', 'CSR', 'ADMIN', 'DEVELOPER')),
    status VARCHAR(20) DEFAULT 'OFFLINE' CHECK (status IN ('ONLINE', 'OFFLINE')),
    permission VARCHAR(20) DEFAULT 'ACTIVE' CHECK (permission IN ('ACTIVE', 'DEACTIVE', 'DELETED'))
);

-- 2. 建立 games 資料表
CREATE TABLE games (
    game_id SERIAL PRIMARY KEY,
    developer_id INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    price DECIMAL(10, 2) NOT NULL DEFAULT 0.00 CHECK (price >= 0),
    overall_rating DECIMAL(3, 2) DEFAULT 0.00
);

-- 3. 建立 tags 資料表
CREATE TABLE tags (
    tag_id SERIAL PRIMARY KEY,
    tag_name VARCHAR(50) UNIQUE NOT NULL
);

-- 4. 建立 game_tags 中介資料表
CREATE TABLE game_tags (
    game_id INT REFERENCES games(game_id) ON DELETE CASCADE,
    tag_id INT REFERENCES tags(tag_id) ON DELETE CASCADE,
    PRIMARY KEY (game_id, tag_id)
);

-- 5. 建立 game_media 資料表
CREATE TABLE game_media (
    media_id SERIAL PRIMARY KEY,
    game_id INT NOT NULL REFERENCES games(game_id) ON DELETE CASCADE,
    file_url VARCHAR(500) NOT NULL,
    media_type VARCHAR(20) DEFAULT 'media' CHECK (media_type IN ('media', 'game_file'))
);

-- 6. 建立 transactions 資料表
CREATE TABLE transactions (
    transaction_id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    total_amount DECIMAL(10, 2) NOT NULL DEFAULT 0.00 CHECK (total_amount >= 0),
    transaction_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    receipt_number VARCHAR(100) UNIQUE NOT NULL
);

-- 7. 建立 transaction_items 資料表
CREATE TABLE transaction_items (
    item_id SERIAL PRIMARY KEY,
    transaction_id INT NOT NULL REFERENCES transactions(transaction_id) ON DELETE CASCADE,
    game_id INT NOT NULL REFERENCES games(game_id) ON DELETE CASCADE,
    purchase_price DECIMAL(10, 2) NOT NULL DEFAULT 0.00 CHECK (purchase_price >= 0)
);

-- 8. 建立 shopping_carts 資料表
CREATE TABLE shopping_carts (
    cart_id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    game_id INT NOT NULL REFERENCES games(game_id) ON DELETE CASCADE,
    added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 9. 建立 game_licenses 資料表
CREATE TABLE game_licenses (
    license_id SERIAL PRIMARY KEY,
    game_id INT NOT NULL REFERENCES games(game_id) ON DELETE CASCADE,
    user_id INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    transaction_item_id INT NOT NULL REFERENCES transaction_items(item_id) ON DELETE CASCADE,
    acquired_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(20) DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'REVOKED', 'FROZEN'))
);

-- 10. 建立 wish_lists 資料表
CREATE TABLE wish_lists (
    wishlist_id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    game_id INT NOT NULL REFERENCES games(game_id) ON DELETE CASCADE,
    added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, game_id)
);


-- 11. 建立 refund_requests 資料表
CREATE TABLE refund_requests (
    refund_id SERIAL PRIMARY KEY,
    buyer_id INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    transaction_item_id INT NOT NULL REFERENCES transaction_items(item_id) ON DELETE CASCADE,
    handled_by INT REFERENCES users(user_id) ON DELETE SET NULL,
    reason TEXT NOT NULL,
    reject_reason TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    resolved_at TIMESTAMP,
    status VARCHAR(20) DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'APPROVED', 'REJECTED'))
);

-- 12. 評論 
CREATE TABLE reviews (
    review_id SERIAL PRIMARY KEY,
    game_id INT NOT NULL REFERENCES games(game_id) ON DELETE CASCADE,
    user_id INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    attitude VARCHAR(20) NOT NULL CHECK (attitude IN ('POSITIVE', 'NEGATIVE')),
    status VARCHAR(20) DEFAULT 'VISIBLE' CHECK (status IN ('VISIBLE', 'HIDDEN', 'DELETED'))
);

-- 13. 評論回覆
CREATE TABLE review_replies (
    review_reply_id SERIAL PRIMARY KEY,
    review_id INT NOT NULL REFERENCES reviews(review_id) ON DELETE CASCADE,
    user_id INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    parent_reply_id INT REFERENCES review_replies(review_reply_id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    status VARCHAR(20) DEFAULT 'VISIBLE' CHECK (status IN ('VISIBLE', 'HIDDEN', 'DELETED'))
);

-- 14. 好友關係
CREATE TABLE friendships (
    friendship_id SERIAL PRIMARY KEY,
    sender_id INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    receiver_id INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(20) DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'ACCEPTED', 'DECLINED')),
    UNIQUE (sender_id, receiver_id)
);

-- 15. 私訊
CREATE TABLE messages (
    message_id SERIAL PRIMARY KEY,
    sender_id INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    receiver_id INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    sent_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_read BOOLEAN DEFAULT FALSE
);

-- 16. 黑名單 (BLACKLIST)
CREATE TABLE blacklists (
    blacklist_id SERIAL PRIMARY KEY,
    blocker_id INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    blocked_id INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (blocker_id, blocked_id)
);

