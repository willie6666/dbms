# 1. 帳號生命週期與權限管理 (Account Lifecycle)

這份文件詳細描述了 VaporAuror 平台中，從使用者註冊開始，一直到個人資料維護、權限變更，甚至是帳號被停權或永久刪除的完整閉環業務邏輯。

---

## 1. 註冊與身分分派 (Registration & Role Assignment)

- **起點**：訪客 (GUEST) 進入 `login.html` 並切換到註冊表單。
- **流程**：
  1. 訪客填寫 `Username`, `Email`, `Password`，並可勾選「註冊為開發者 (Developer)」。
  2. 前端透過 `POST /api/auth/register` 將資料送至後端。
  3. 後端檢查 `Email` 是否已存在。若無重複，對密碼進行 Bcrypt 雜湊處理。
  4. 依據勾選狀態，將 `Role` 設定為 `'USER'` 或 `'DEVELOPER'`，寫入 `users` 表單。
  5. **注意**：客服 (`CSR`) 與 管理員 (`ADMIN`) 角色無法自行註冊，必須由現有管理員在後台手動賦予。
- **終點**：回傳 `201 Created`，前端顯示成功訊息，引導使用者至登入頁面。

---

## 2. 登入與 Token 核發 (Login & JWT)

- **起點**：訪客在 `login.html` 輸入 `Email` 與 `Password`。
- **流程**：
  1. 前端呼叫 `POST /api/auth/login`。
  2. 後端比對 `Email` 與 Bcrypt 雜湊。
  3. **狀態檢查**：如果 `users.permission` 是 `'DEACTIVE'` (停權) 或 `'DELETED'` (刪除)，則拒絕登入，回傳 `403 Forbidden`。
  4. 登入成功後，後端產生一組 JWT Token (內含 `user_id`, `role`, `permission`) 並回傳。
  5. 前端將 Token 與基本資料儲存至 `localStorage`。
- **終點**：依據 `role`，前端自動跳轉：
  - `USER` / `CSR` -> 跳轉至商店首頁 `index.html`。
  - `DEVELOPER` -> 跳轉至開發者後台 `dev_dashboard.html`。
  - `ADMIN` -> 跳轉至系統管理員後台 `admin_dashboard.html`。

---

## 3. 個人資料修改 (Profile Management)

- **起點**：已登入的使用者進入 `settings.html`。
- **流程**：
  1. 畫面載入時，前端呼叫 `GET /api/users/profile` 取得目前的 `Username`, `Email`, `Bio` (個人簡介), `AvatarURL`。
  2. 使用者可以修改這三種屬性，送出後前端呼叫 `PUT /api/users/profile`。
  3. 後端驗證若修改 `Email`，是否與他人重複。
  4. 更新 `users` 資料表對應欄位。
- **終點**：更新成功，前端重新整理 `localStorage` 中的快取資料，確保右上角的頭像與名稱即時更新。

---

## 4. 帳號停權與永久刪除 (Suspension & Deletion)

- **起點**：擁有 `ADMIN` 權限的系統管理員，進入 `admin_dashboard.html` 點擊「停權」或「永久刪除」。
- **流程 (停權 Suspend)**：
  1. 呼叫 `PUT /api/admin/users/:id/suspend`。
  2. 後端將 `users.permission` 在 `'ACTIVE'` 與 `'DEACTIVE'` 之間切換。
  3. 影響：該玩家下次登入時會被阻擋；已登入的 Token 若過期後也無法再更新。
- **流程 (永久刪除 Delete)**：
  1. 呼叫 `DELETE /api/admin/users/:id`。
  2. 後端進行 **Soft Delete (軟刪除)**，將 `users.permission` 設為 `'DELETED'`。
  3. 為了保護已經購買該開發者遊戲的玩家，**系統會觸發連鎖反應 (Cascade)**：
     - 若被刪除者是 `DEVELOPER`，系統會自動找出他名下的所有遊戲。
     - 將這些遊戲的狀態 `games.status` 強制設為 `'TAKEN_DOWN'` (下架)。
     - **更進一步**：將所有玩家對於這些遊戲的購買授權 `game_licenses.status` 設定為 `'REVOKED'` (撤銷)。
- **終點**：該帳號永遠無法登入；若為開發者，其心血與販售中的遊戲會徹底從平台上消失，玩家的遊戲庫也會同步清空這款遊戲。
