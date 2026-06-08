# g. 修改資料表結構的操作 (DDL: Data Definition Language)

本文件列出系統進行資料表建立、欄位修改或結構變更時所使用的 DDL 操作。在我們的 VaporAuror 專案中，實體資料表的定義主要由 `db/01_init_table.sql` 執行，同時 Go 的 GORM `AutoMigrate` 也會在啟動時自動檢查並補齊必要的綱要 (Schema)。

---

### 1. 建立新的資料表 (CREATE TABLE)
- **說明**：用於從無到有建立系統核心的實體表。此範例為建立遊戲主檔 (`games`)，並設定了主鍵 `SERIAL` (自動遞增)、字元長度限制以及預設值 `DEFAULT`。
- **觸發時機**：Docker Compose 啟動掛載初始化 SQL 或 GORM `AutoMigrate` 時。
- **原生 SQL 語法**：
  ```sql
  CREATE TABLE games (
      game_id SERIAL PRIMARY KEY,
      title VARCHAR(100) NOT NULL,
      developer_id INT NOT NULL,
      description TEXT,
      price DECIMAL(10, 2) NOT NULL,
      release_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
      status VARCHAR(20) DEFAULT 'ACTIVE'
  );
  ```

### 2. 新增外部鍵約束 (ALTER TABLE ... ADD CONSTRAINT)
- **說明**：這是修改現有資料表結構 (`ALTER TABLE`) 的重要應用。在建立完 `games` 表與 `users` 表後，必須透過 `ALTER TABLE` 補上外鍵約束，強制規定遊戲的 `developer_id` 必須真實存在於 `users` 表中，以確保資料完整性 (Referential Integrity)。
- **觸發時機**：資料表關聯建立時。
- **原生 SQL 語法**：
  ```sql
  ALTER TABLE games 
  ADD CONSTRAINT fk_games_developer 
  FOREIGN KEY (developer_id) 
  REFERENCES users(user_id)
  ON DELETE CASCADE;
  ```

### 3. 刪除整張資料表 (DROP TABLE)
- **說明**：這是極具破壞性的 DDL 操作，用於連根拔起銷毀整張資料表。`IF EXISTS` 可以避免表不存在時報錯，而 `CASCADE` 則會一併刪除那些依賴這張表的其他視圖 (Views) 或約束 (Constraints)。
- **觸發時機**：系統重置、清空測試資料或破壞性資料庫更新時。
- **原生 SQL 語法**：
  ```sql
  DROP TABLE IF EXISTS shopping_carts CASCADE;
  ```
