# 1. 帳號生命週期與個人設定 (Account Lifecycle)

本文件描述玩家或開發者從註冊、登入、一直到修改個人基本資料的完整閉環業務邏輯。此流程對應前端的 `login.html`, `register.html` 與 `profile.html`，並由後端 `auth_controller.go` 與 `user_controller.go` 處理。

---

## 1. 帳號註冊 (Registration)

- **起點**：使用者在 `register.html` 填寫註冊表單。
- **輸入資料**：`username`, `email`, `password`，以及一個可選的「我是開發者 (IsDeveloper)」勾選框。
- **流程與後端驗證 (`POST /api/auth/register`)**：
  1. **基礎驗證**：後端檢查密碼長度必須 `>= 6`，且 Email 必須包含 `@`。
  2. **密碼加密**：使用 Bcrypt 將密碼雜湊化後存入資料庫。
  3. **角色分派 (Role Assignment)**：
     - 若 `is_developer` 為 `true`，賦予 `DEVELOPER` 角色。
     - 若為 `false`，賦予 `USERS` 角色。
  4. **預設狀態**：帳號建立時，系統預設將 `Status` 設為 `OFFLINE`，`Permission` 設為 `ACTIVE`。
  5. **資料庫寫入**：寫入 `users` 表。若 `username` 或 `email` 已被占用，會由 GORM 的 Unique 限制擋下並回傳錯誤。
- **終點 (自動登入)**：註冊成功後，後端會直接核發並回傳 JWT Token，前端收到後存入 `localStorage` 並自動導向商店首頁 (`index.html`)，完成無縫登入體驗。

---

## 2. 帳號登入與權限檢查 (Login & Permission Check)

- **起點**：使用者在 `login.html` 輸入 `email` 與 `password`。
- **流程與後端驗證 (`POST /api/auth/login`)**：
  1. **信箱比對**：檢查 `users` 表中是否有對應的 Email。
  2. **密碼驗證**：透過 Bcrypt 比對密碼。
  3. **核心權限防護 (Permission Check)**：
     - 系統會檢查使用者的 `Permission` 欄位。
     - 若 `user.Permission != "ACTIVE"` (例如被管理員標記為 `DEACTIVE` 停權，或 `DELETED` 刪除)，登入將被**強制拒絕**並回傳 HTTP 403 `This account is not active`。
     - 若 `user.Role == "NULL"` 也會被拒絕登入。
  4. **核發 Token**：驗證全數通過後，核發包含 `user_id` 與 `role` 的 JWT Token。
- **終點**：前端儲存 Token，依照 `Role` 決定導覽列顯示的特殊入口 (例如「開發者中心」或「管理後台」)，並導向商店首頁。

---

## 3. 個人資料修改 (Profile Update)

- **起點**：已登入的使用者進入 `profile.html`。
- **流程與後端驗證 (`PUT /api/users/profile`)**：
  1. **身分驗證**：必須帶有合法的 JWT Token (`RequireAuth` Middleware)。
  2. **輸入檢查**：使用者可選擇性修改 `username`, `email` 或 `password`。
     - 若修改密碼，長度必須 `>= 6`。
  3. **唯一性防撞檢查 (Uniqueness Check)**：
     - 若修改 `username`，後端會查詢是否已有**其他**使用者使用該名稱 (排除自己)。若有則回傳 `Username already taken`。
     - 若修改 `email`，同理查詢是否有其他使用者使用。
  4. **更新資料**：將變更寫入資料庫 (`DB.Save`)。
- **終點**：前端收到更新成功的提示，並更新 `localStorage` 中快取的使用者資訊，介面右上角的名稱同步變更。

---

## 4. 登出 (Logout)

- **流程**：本系統採用無狀態 JWT，因此「登出」主要由前端清除 `localStorage` 中的 Token 與使用者資訊。後端的 `POST /api/auth/logout` 僅作為一個象徵性的確認端點，未來可擴充作為 Token 黑名單 (Blacklist) 機制的掛載點。
