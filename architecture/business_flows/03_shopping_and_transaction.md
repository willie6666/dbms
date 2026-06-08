# 3. 玩家購物車與結帳金流 (Shopping & Transaction)

本文件描述玩家在商店尋找遊戲、加入購物車、一直到成功結帳並取得遊戲授權的完整閉環。

---

## 1. 商店瀏覽與願望清單 (Browse & Wishlist)

- **起點**：玩家瀏覽 `index.html` 或 `game_detail.html`。
- **流程 (搜尋)**：
  - 玩家可以透過 `GET /api/games?q={keyword}&tag={tag}` 進行過濾，過濾器只會返回 `status = 'ACTIVE'` 且 `price >= min_price` 的遊戲。
- **流程 (加入願望清單)**：
  1. 玩家對感興趣的遊戲點擊「加入願望清單」，呼叫 `POST /api/protected/wishlist`。
  2. 後端檢查是否已存在於 `wish_lists` 表。若無，則寫入資料。
- **終點**：玩家可以在 `wishlist.html` 看到所有收藏的遊戲。若收藏的遊戲被管理員下架，清單中會出現紅色的「已下架」警告。

---

## 2. 購物車管理 (Cart Management)

- **起點**：玩家在商店或願望清單中點擊「加入購物車」。
- **流程**：
  1. 前端呼叫 `POST /api/protected/cart`。
  2. **防護檢查 1**：後端檢查該遊戲的 `status` 是否為 `'ACTIVE'`。若非 ACTIVE (例如 DRAFT 或 TAKEN_DOWN)，則阻擋加入。
  3. **防護檢查 2**：檢查玩家是否已經擁有這款遊戲 (在 `game_licenses` 中尋找 `ACTIVE` 的紀錄)。若已擁有，拒絕加入。
  4. 寫入 `shopping_carts` 資料表。
- **管理購物車**：
  - 玩家開啟購物車頁面，前端呼叫 `GET /api/protected/cart` 列出所有商品。
  - 若玩家想移除特定商品，可點擊移除並呼叫 `DELETE /api/protected/cart/:game_id`。
- **終點**：前端右上角的購物車圖示數字增加或減少，隨時保持與資料庫同步。

---

## 3. 結帳與發放授權 (Checkout & Licensing)

- **起點**：玩家進入 `cart.html`，確認總金額後點擊「確認結帳」。
- **流程**：
  1. 前端向 `POST /api/protected/checkout` 發出請求。
  2. 後端利用資料庫的交易機制 (`DB.Begin()`) 開啟一個 Transaction 確保 ACID 屬性。
  3. **重新驗證**：
     - 檢查購物車內是否有非 `ACTIVE` 的遊戲，若有則直接 Rollback 回傳錯誤。
  4. **紀錄建立**：
     - 在 `transactions` 表建立一筆交易主檔 (包含總金額與時間)。
     - 針對每一款遊戲，在 `transaction_items` 表建立明細 (包含單款遊戲的購買價格 `purchase_price`)。
  5. **發放授權**：
     - 針對每一款購買的遊戲，在 `game_licenses` 中新增一筆紀錄，綁定 `user_id`, `game_id` 與 `transaction_item_id`，狀態為 `'ACTIVE'`。
     - 授權發放完畢後，清空該玩家的 `shopping_carts`。
  6. 交易提交 (`DB.Commit()`)。
- **終點**：結帳成功，玩家被引導至「我的遊戲庫 (`library.html`)」，可看到剛買的遊戲。
