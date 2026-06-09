-- ==========================================
-- 10. 新增三個遊戲與其對應的 Media
-- ==========================================
WITH inserted_games AS (
    INSERT INTO games (developer_id, title, description, price, overall_rating, status)
    VALUES 
        (2, 'Shake Game', '一款關於搖晃的有趣遊戲', 10.00, 4.5, 'ACTIVE'),
        (2, 'Water Fly Adventure', '水上飛行冒險', 20.00, 4.8, 'ACTIVE'),
        (2, 'Wing Simulator', '展翅翱翔模擬器', 30.00, 4.2, 'ACTIVE')
    RETURNING game_id, title
)
INSERT INTO game_media (game_id, file_url, media_type)
SELECT 
    game_id,
    CASE 
        WHEN title = 'Shake Game' THEN '/media/images/' || game_id || '/shake.gif'
        WHEN title = 'Water Fly Adventure' THEN '/media/images/' || game_id || '/water_fly.gif'
        WHEN title = 'Wing Simulator' THEN '/media/images/' || game_id || '/wing.gif'
    END,
    'media'
FROM inserted_games;
