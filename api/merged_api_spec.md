# VaporAuror 詳細 API 規格書 (API Specification) — Merged

這份文件專為**前端開發人員**撰寫。詳細列出所有 API 的傳入參數 (Request Body)、必要的標頭 (Headers)、以及各種成功與失敗情境的回傳格式 (Response JSON)。

> 所有路徑皆為同源相對路徑；瀏覽器經由 Caddy 入口代理 `/api/*` 到後端。
> 本文件由 `add_api_spec.md` 與 `bananaapi_spec.md` 合併而成，所有差異處已經過 **Go server 原始碼驗證**。

## ⚠️ 通用錯誤回應 (Global Errors)
在所有端點中，如果發生以下情況，後端會統一回傳對應的錯誤代碼，下方各端點的說明中將**不再贅述**這些基本錯誤：
- `400 Bad Request`: `{"error": "..."}` (傳入的 JSON 格式錯誤、缺少必填欄位 `binding:"required"`)
- `401 Unauthorized`: `{"error": "..."}` (未登入、JWT Token 缺失或無效)
- `403 Forbidden`: `{"error": "Forbidden: Requires <ROLE> role"}` (權限不足，例如一般玩家呼叫 ADMIN API)
- `500 Internal Server Error`: `{"error": "..."}` (資料庫連線失敗、伺服器內部錯誤)

---

## 1. 使用者與權限 (Users & Auth)

### `[POST] /api/auth/register` (註冊新帳號)
- **Headers**: 無
- **Request Body**:
  ```json
  {
    "username": "PlayerOne",
    "email": "player1@test.com",
    "password": "password123" // 必填，長度需 >= 6
  }
  ```
- **Responses**:
  - `201 Created`:
    ```json
    {
      "message": "Registration successful",
      "token": "eyJhbGciOi...",
      "user": {
        "id": 1,
        "username": "PlayerOne",
        "email": "player1@test.com",
        "role": "USERS"
      }
    }
    ```
    > 註冊成功後會自動登入，回傳 JWT Token 與使用者資訊。
  - `400 Bad Request`: `{"error": "..."}` (密碼太短或格式錯誤)
  - `500 Internal Server Error`: `{"error": "Failed to create user (username or email might already exist)"}`

### `[POST] /api/auth/login` (使用者登入)
- **Headers**: 無
- **Request Body**:
  ```json
  {
    "email": "player1@test.com",
    "password": "password123"
  }
  ```
- **Responses**:
  - `200 OK`:
    ```json
    {
      "message": "Login successful",
      "token": "eyJhbGciOi...",
      "user": {
        "id": 1,
        "username": "PlayerOne",
        "email": "player1@test.com",
        "role": "USERS"
      }
    }
    ```
  - `401 Unauthorized`: `{"error": "Invalid email or password"}`
  - `403 Forbidden`: `{"error": "This account is not active"}` (帳號已被停權或角色為 NULL)

### `[POST] /api/auth/logout` (登出)
- **Headers**: `Authorization: Bearer <token>`
- **Request Body**: 無
- **Responses**:
  - `200 OK`: `{"message": "Logged out successfully. Please remove your token."}`
  - **說明**: 由於採用 JWT 無狀態架構，後端只會回傳成功訊息，真正的登出必須由前端主動清除 Token。

### `[PUT] /api/users/profile` (修改個人資料)
- **Headers**: `Authorization: Bearer <token>`
- **Request Body**: (皆為選填，想改什麼傳什麼)
  ```json
  {
    "username": "NewName",
    "email": "new@test.com",
    "password": "newpassword"
  }
  ```
- **Responses**:
  - `200 OK`:
    ```json
    {
      "message": "Profile updated successfully",
      "user": {
        "id": 1,
        "username": "NewName",
        "email": "new@test.com",
        "role": "USERS"
      }
    }
    ```
  - `404 Not Found`: `{"error": "User not found"}`

### `[GET] /api/admin/users` (查看所有使用者清單)
- **Headers**: `Authorization: Bearer <admin_token>`
- **Responses**:
  - `200 OK`: `{"data": [ { "id": 1, "username": "PlayerOne", "role": "USERS" } ]}`
  - `403 Forbidden`: `{"error": "Forbidden: Requires ADMIN role"}`

### `[PUT] /api/admin/users/{id}/suspend` (切換帳號停權狀態)
- **Headers**: `Authorization: Bearer <admin_token>`
- **Request Body**: 無
- **說明**: 此端點為 **Toggle** 行為。若帳號為 `ACTIVE` 則切換為 `DEACTIVE`；若為 `DEACTIVE` 則切換回 `ACTIVE`。
- **Responses**:
  - `200 OK`: `{"message": "User account has been suspended"}`
  - `404 Not Found`: `{"error": "User not found"}`

### `[DELETE] /api/admin/users/{id}` (移除帳號)
- **Headers**: `Authorization: Bearer <admin_token>`
- **Request Body**: 無
- **Responses**:
  - `200 OK`: `{"message": "User completely removed"}`

### `[PUT] /api/admin/users/{id}/role` (更改帳號權限)
- **Headers**: `Authorization: Bearer <admin_token>`
- **Request Body**:
  ```json
  { "role": "DEVELOPER" } // 'USERS', 'CSR', 'DEVELOPER', 'ADMIN'
  ```
- **Responses**:
  - `200 OK`: `{"message": "User role updated successfully"}`

---

## 2. 商店與遊戲 (Store & Games)

### `[GET] /api/games` (瀏覽/搜尋遊戲)
- **Headers**: 無
- **Query Params** (可選): `?q=elden`
- **Responses**:
  - `200 OK`:
    ```json
    {
      "data": [
        {
          "game_id": 1,
          "title": "Elden Ring",
          "price": 1290,
          "overall_rating": 4.8
        }
      ]
    }
    ```

### `[GET] /api/games/{id}` (查看遊戲詳情)
- **Headers**: 無
- **Responses**:
  - `200 OK`: `{"data": { "game": {...}, "developer_name": "DevUser", "media": [...], "tags": [...], "reviews": [...] }}`
  - `404 Not Found`: `{"error": "Game not found"}`

### `[GET] /api/games/{id}/reviews` (查看遊戲評論)
- **Headers**: 無
- **Responses**:
  - `200 OK`: `[ { "review_id": 1, "content": "...", "attitude": "POSITIVE", "user": {...}, "replies": [...] } ]`
  - **注意**: 回傳格式為陣列 (非包在 `{"data": [...]}` 內)。

### `[GET] /api/developer/games` (查看自己的遊戲列表)
- **Headers**: `Authorization: Bearer <developer_token>`
- **Responses**:
  - `200 OK`: `{"data": [ { ...game_objects_with_media... } ]}`
  - **說明**: DEVELOPER 只會看到自己上架的遊戲；ADMIN 可查看全部遊戲。

### `[POST] /api/developer/games` (上架新遊戲)
- **Headers**: `Authorization: Bearer <developer_token>`
- **Request Body**:
  ```json
  {
    "title": "My Indie Game",    // 必填
    "price": 350.00,             // 必填，最小值 0
    "desc": "遊戲描述 (選填，支援 Markdown)"  // 選填
  }
  ```
- **Responses**:
  - `201 Created`: `{"message": "Game uploaded successfully", "game": {...}}`

### `[PUT] /api/developer/games/{id}` (編輯遊戲資訊)
- **Headers**: `Authorization: Bearer <developer_token>`
- **Request Body**:
  ```json
  {
    "price": 299.00,        // 選填，最小值 0
    "desc": "更新的遊戲描述" // 選填
  }
  ```
  > **注意**: `title` 無法透過此 API 修改；`price` 若傳 `0` 仍會被視為有效值並寫入。ADMIN 可編輯任何遊戲。
- **Responses**:
  - `200 OK`: `{"message": "Game updated successfully", "game": {...}}`
  - `403 Forbidden`: `{"error": "Forbidden: You can only edit your own games"}`
  - `404 Not Found`: `{"error": "Game not found"}`

### `[DELETE] /api/developer/games/{id}` (下架自己的遊戲)
- **Headers**: `Authorization: Bearer <developer_token>`
- **Responses**:
  - `200 OK`: `{"message": "Game deleted successfully"}`
  - `403 Forbidden`: `{"error": "Forbidden: You can only delete your own games"}`
  - `404 Not Found`: `{"error": "Game not found"}`

### `[DELETE] /api/admin/games/{id}` (強制下架遊戲)
- **Headers**: `Authorization: Bearer <admin_token>`
- **Responses**:
  - `200 OK`: `{"message": "Game deleted successfully by Admin"}`

### `[POST] /api/developer/games/{id}/media` (上傳遊戲素材)
- **Headers**: `Authorization: Bearer <developer_token>`
  > **注意**：使用 `fetch` 傳送 `FormData` 時，請**不要**手動設定 `Content-Type` 標頭，瀏覽器會自動補上帶有正確 `boundary` 的 `multipart/form-data` 標頭。
- **Request Body** (`multipart/form-data`):
  | 欄位名稱 | 類型 | 必填 | 說明 |
  |----------|------|------|------|
  | `file` | File | ✅ | 要上傳的圖片或遊戲檔案 |
  | `media_type` | String | 否 | `"media"` (圖片，預設) 或 `"game_file"` (遊戲檔案) |
- **儲存路徑**:
  - `media` → `assets/images/{game_id}/{sha256}.{ext}`，對外 URL `/media/images/{game_id}/{sha256}.{ext}`
  - `game_file` → `assets/game-files/{game_id}/{original_name}`，對外 URL `/downloads/{game_id}/{original_name}`
- **Responses**:
  - `201 Created`: `{"message": "Media uploaded successfully", "data": {...}, "file_url": "/media/images/..."}`
  - `400 Bad Request`: `{"error": "Missing file field"}`
  - `403 Forbidden`: `{"error": "Forbidden: You can only upload media for your own games"}`
  - `404 Not Found`: `{"error": "Game not found"}`

### `[DELETE] /api/developer/games/{id}/media/{media_id}` (刪除遊戲素材)
- **Headers**: `Authorization: Bearer <developer_token>`
- **Responses**:
  - `200 OK`: `{"message": "Media deleted successfully"}`
  - `403 Forbidden`: `{"error": "Forbidden: You can only manage your own games"}`
  - `404 Not Found`: `{"error": "Media not found"}` 或 `{"error": "Game not found"}`

### `[GET] /api/developer/games/{id}/stats` (查看遊戲銷售數據)
- **Headers**: `Authorization: Bearer <developer_token>`
- **Responses**:
  - `200 OK`:
    ```json
    {
      "stats": {
        "total_sales": 125,
        "total_revenue": 45000.50
      }
    }
    ```
  - `403 Forbidden`: `{"error": "Forbidden: You can only view stats for your own games"}`
  - `404 Not Found`: `{"error": "Game not found"}`

### `[GET] /api/tags` (查看所有可用標籤)
- **Headers**: 無
- **Responses**:
  - `200 OK`: `{"data": [ {"tag_id": 1, "tag_name": "RPG"} ]}`

### `[POST] /api/developer/tags` (建立新標籤)
- **Headers**: `Authorization: Bearer <developer_token>`
- **Request Body**: `{"tag_name": "Action"}`
- **Responses**:
  - `201 Created`: `{"message": "Tag created successfully", "data": { "tag_id": 1, "tag_name": "Action" }}`
  - `500 Internal Server Error`: `{"error": "Failed to create tag (might already exist)"}`

### `[POST] /api/developer/games/{id}/tags` (為遊戲貼標籤)
- **Headers**: `Authorization: Bearer <developer_token>`
- **Request Body**: `{"tag_id": 2}`
- **Responses**:
  - `200 OK`: `{"message": "Tag added to game"}`
  - `403 Forbidden`: `{"error": "Forbidden: Not your game"}`
  - `404 Not Found`: `{"error": "Game not found"}`

### `[DELETE] /api/developer/games/{id}/tags/{tag_id}` (移除遊戲標籤)
- **Headers**: `Authorization: Bearer <developer_token>`
- **Responses**:
  - `200 OK`: `{"message": "Tag removed from game"}`
  - `403 Forbidden`: `{"error": "Forbidden"}`
  - `404 Not Found`: `{"error": "Game not found"}`

---

## 3. 訂單、購物車與客服 (Transactions & Carts)

### `[GET] /api/protected/cart` (查看購物車內容)
- **Headers**: `Authorization: Bearer <token>`
- **Responses**:
  - `200 OK`: `{"data": [ { ...cart_items... } ]}`

### `[POST] /api/protected/cart` (放入購物車)
- **Headers**: `Authorization: Bearer <token>`
- **Request Body**: `{"game_id": 1}`
- **Responses**:
  - `200 OK`: `{"message": "Game added to cart successfully"}`
  - `400 Bad Request`: `{"error": "Game already in cart"}` 或 `{"error": "You already own this game"}`

### `[DELETE] /api/protected/cart/{game_id}` (移出購物車)
- **Headers**: `Authorization: Bearer <token>`
- **Responses**:
  - `200 OK`: `{"message": "Game removed from cart"}`

### `[POST] /api/protected/checkout` (結帳)
- **Headers**: `Authorization: Bearer <token>`
- **Request Body**: 無 (自動結算購物車內所有物品)
- **Responses**:
  - `200 OK`: `{"message": "Checkout successful. Games added to your library!"}`
  - `500 Internal Server Error`: `{"error": "Checkout failed: Cart is empty"}`

### `[GET] /api/protected/transactions` (查看購買紀錄)
- **Headers**: `Authorization: Bearer <token>`
- **Responses**:
  - `200 OK`: `{"data": [ { ...transactions... } ]}`

### `[POST] /api/social/refunds` (申請退款)
- **Headers**: `Authorization: Bearer <token>`
- **Request Body**:
  ```json
  {
    "transaction_item_id": 105,
    "reason": "遊戲有嚴重 Bug 無法執行"
  }
  ```
- **Responses**:
  - `201 Created`: `{"message": "Refund request submitted. A CSR will review it shortly."}`
  - `400 Bad Request`: `{"error": "A refund request already exists for this item"}`
  - `403 Forbidden`: `{"error": "Forbidden: Transaction item not found in your library"}`

### `[GET] /api/csr/refunds` (查看待處理退款)
- **Headers**: `Authorization: Bearer <csr_token>`
- **Responses**:
  - `200 OK`: `{"data": [ { ...pending_refunds... } ]}`

### `[PUT] /api/csr/refunds/{id}` (同意/拒絕退款)
- **Headers**: `Authorization: Bearer <csr_token>`
- **Request Body**:
  ```json
  {
    "status": "APPROVED", // 'APPROVED' 或 'REJECTED'
    "reject_reason": ""
  }
  ```
- **Responses**:
  - `200 OK`: `{"message": "Refund processed successfully"}`
  - `500 Internal Server Error`: `{"error": "Failed to process refund. Is it already processed?"}`

---

## 4. 遊戲庫與願望清單 (Library & Wishlist)

### `[GET] /api/protected/library` (顯示個人遊戲庫)
- **Headers**: `Authorization: Bearer <token>`
- **Responses**:
  - `200 OK`: `{"data": [ { "license_id": 1, "game_id": 5, "status": "ACTIVE" } ]}`

### `[GET] /api/protected/wishlist` (查看願望清單)
- **Headers**: `Authorization: Bearer <token>`
- **Responses**:
  - `200 OK`: `{"data": [ ...wishlist_items... ]}`

### `[POST] /api/protected/wishlist` (加入願望清單)
- **Headers**: `Authorization: Bearer <token>`
- **Request Body**: `{"game_id": 3}`
- **Responses**:
  - `200 OK`: `{"message": "Added to wishlist"}`
  - `500 Internal Server Error`: `{"error": "Failed to add to wishlist (might already exist)"}`

### `[DELETE] /api/protected/wishlist/{game_id}` (移除願望清單)
- **Headers**: `Authorization: Bearer <token>`
- **Responses**:
  - `200 OK`: `{"message": "Removed from wishlist"}`

### `[GET] /api/protected/library/{game_id}/play` (玩遊戲)
- **Headers**: `Authorization: Bearer <token>`
- **Responses**:
  - `200 OK`: `{"message": "Game launched successfully", "auth_token": "mock-play-token-12345"}`
  - `403 Forbidden`: `{"error": "You do not own this game or the license is inactive"}`

### `[GET] /api/protected/library/{game_id}/download` (下載遊戲)
- **Headers**: `Authorization: Bearer <token>`
- **Responses**:
  - `200 OK`: 直接回傳檔案串流 (binary)，附帶 `Content-Disposition: attachment; filename="{filename}"` 標頭。
    > **注意**：由於此端點需要 JWT Token 認證，前端**不能**直接使用 `window.location.href = url` 或單純的 `<a>` 標籤。前端必須使用帶有 `Authorization` 標頭的 `fetch` 取得回應，將其轉換為 `Blob`，再透過 `URL.createObjectURL(blob)` 建立臨時網址供 `<a>` 標籤觸發下載。
  - `403 Forbidden`: `{"error": "You do not own this game or the license is inactive"}`
  - `404 Not Found`: `{"error": "No downloadable game file is available"}`

---

## 5. 社交、評論與通訊 (Social & Reviews)

### `[POST] /api/social/games/{id}/reviews` (對遊戲發表評價)
- **Headers**: `Authorization: Bearer <token>`
- **Request Body**:
  ```json
  {
    "attitude": "POSITIVE", // 'POSITIVE' 或 'NEGATIVE'
    "content": "神作不解釋！"
  }
  ```
- **Responses**:
  - `201 Created`: `{"message": "Review posted successfully"}`
  - `403 Forbidden`: `{"error": "Forbidden: You must own the game to leave a review"}`
  - `404 Not Found`: `{"error": "Game not found"}`

### `[POST] /api/social/reviews/{review_id}/replies` (樓中樓回覆)
- **Headers**: `Authorization: Bearer <token>`
- **Request Body**:
  ```json
  {
    "parent_reply_id": 5, // 選填。如果有帶值，代表是「對回覆的回覆」
    "content": "我完全同意這篇評論！"
  }
  ```
- **Responses**:
  - `201 Created`: `{"message": "Reply posted successfully", "data": {...}}`
  - `404 Not Found`: `{"error": "Review not found"}`

### `[DELETE] /api/social/reviews/replies/{reply_id}` (刪除樓中樓回覆)
- **Headers**: `Authorization: Bearer <token>`
- **Responses**:
  - `200 OK`: `{"message": "Reply deleted successfully"}`
  - `403 Forbidden`: `{"error": "Forbidden: You can only delete your own replies"}`
  - `404 Not Found`: `{"error": "Reply not found"}`

### `[GET] /api/social/friends` (查看好友列表)
- **Headers**: `Authorization: Bearer <token>`
- **Responses**:
  - `200 OK`: `{"data": [ { "friendship_id": 1, "status": "ACCEPTED" } ]}`

### `[GET] /api/social/friends/requests` (查看待審核邀請)
- **Headers**: `Authorization: Bearer <token>`
- **Responses**:
  - `200 OK`: `{"data": [ { "friendship_id": 2, "status": "PENDING" } ]}`

### `[POST] /api/social/friends/request` (發送好友邀請)
- **Headers**: `Authorization: Bearer <token>`
- **Request Body**: (擇一即可)
  ```json
  { "receiver_id": 2 }
  ```
  或依使用者名稱查找：
  ```json
  { "username": "PlayerTwo" }
  ```
- **Responses**:
  - `201 Created`: `{"message": "Friend request sent"}`
  - `400 Bad Request`: `{"error": "Friend request already exists"}` 或 `{"error": "Cannot send a friend request to yourself"}` 或 `{"error": "receiver_id or username is required"}`
  - `404 Not Found`: `{"error": "User not found"}` (使用 username 查找時)

### `[PUT] /api/social/friends/request/{id}/accept` (接受好友邀請)
- **Headers**: `Authorization: Bearer <token>`
- **Responses**:
  - `200 OK`: `{"message": "Friend request accepted"}`
  - `403 Forbidden`: `{"error": "Forbidden: You are not the receiver"}`
  - `404 Not Found`: `{"error": "Friend request not found"}`

### `[PUT] /api/social/friends/request/{id}/decline` (拒絕好友邀請)
- **Headers**: `Authorization: Bearer <token>`
- **Responses**:
  - `200 OK`: `{"message": "Friend request declined"}`
  - `403 Forbidden`: `{"error": "Forbidden: You are not the receiver"}`
  - `404 Not Found`: `{"error": "Friend request not found"}`

### `[DELETE] /api/social/friends/request/{id}` (收回/解除好友)
- **Headers**: `Authorization: Bearer <token>`
- **Responses**:
  - `200 OK`: `{"message": "Friend request revoked / removed"}`
  - `403 Forbidden`: `{"error": "Forbidden"}`
  - `404 Not Found`: `{"error": "Friend request not found"}`

### `[POST] /api/social/messages` (傳輸文字訊息)
- **Headers**: `Authorization: Bearer <token>`
- **Request Body**:
  ```json
  {
    "receiver_id": 2,
    "content": "今晚一起打副本嗎？"
  }
  ```
- **Responses**:
  - `200 OK`: `{"message": "Message sent"}`

### `[GET] /api/social/messages/{user_id}` (讀取對話紀錄)
- **Headers**: `Authorization: Bearer <token>`
- **Responses**:
  - `200 OK`: `{"data": [ { ...message_history... } ]}`

### `[GET] /api/social/blacklist` (查看黑名單列表)
- **Headers**: `Authorization: Bearer <token>`
- **Responses**:
  - `200 OK`: `{"data": [ { "blacklist_id": 1, "blocked_id": 5 } ]}`

### `[POST] /api/social/blacklist` (加入黑名單)
- **Headers**: `Authorization: Bearer <token>`
- **Request Body**: `{"blocked_id": 5}` (也可使用 `{"user_id": 5}` 作為替代欄位名稱)
- **Responses**:
  - `201 Created`: `{"message": "User added to blacklist"}`
  - `400 Bad Request`: `{"error": "Cannot blacklist yourself"}` 或 `{"error": "blocked_id is required"}`

### `[DELETE] /api/social/blacklist/{user_id}` (移除黑名單)
- **Headers**: `Authorization: Bearer <token>`
- **Responses**:
  - `200 OK`: `{"message": "User removed from blacklist"}`
