# 7. 平台最高管理員操作 (Admin Operations)

本文件描述擁有最高權限的 `ADMIN` (系統管理員) 在平台上的終極治理流程。Admin 的權限凌駕於其他所有流程之上，可以對使用者、遊戲、評論進行強制干預。

---

## 1. 使用者角色指派與剝奪 (Role Management)

- **起點**：Admin 在 `admin_dashboard.html` 觀看全站帳號列表 (呼叫 `GET /api/admin/users`)。
- **流程**：
  1. 對某位普通 `USER`，Admin 可以從下拉選單選擇變更為 `CSR` (客服人員)。
  2. 呼叫 `PUT /api/admin/users/:id/role`，送出 `{ "role": "CSR" }`。
  3. 後端更新 `users.role`。
- **終點**：該玩家下次登入時，將能看見「客服中心」的按鈕，並具有處理退款的權限。反之亦可降級回 `USER`。

---

## 2. 帳號停權管理 (Account Suspension)

- **起點**：Admin 發現某位玩家有違規行為，但在 `admin_dashboard.html` 點擊「停權」按鈕。
- **流程**：
  1. 呼叫 `PUT /api/admin/users/:id/suspend`。
  2. 此端點為 Toggle 行為，若原本為 `ACTIVE` 會變更為 `DEACTIVE`；若原本為 `DEACTIVE` 則解除停權變為 `ACTIVE`。
- **終點**：被停權的使用者將無法登入系統 (在 Login 階段會被擋下)。

---

## 3. 強制下架違規遊戲 (Force Takedown & License Revocation)

- **起點**：Admin 發現某款遊戲 (`games.status = 'ACTIVE'`) 涉及侵權、惡意軟體或嚴重違規，在 `admin_dashboard.html` 對該遊戲點擊「強制下架」。
- **流程**：
  1. 呼叫 `DELETE /api/admin/games/:id` (雖然是 DELETE，但為確保歷史交易明細不斷鏈，實施的是軟刪除)。
  2. 後端將 `games.status` 強制設定為 `'TAKEN_DOWN'`。
  3. **終極撤銷授權 (Revoke Licenses)**：
     - Admin 下架與「開發者自主下架」最大的不同在於：Admin 下架帶有懲罰與保護玩家的性質。
     - 系統會透過 `game_id` 到 `game_licenses` 中尋找所有擁有此遊戲的買家。
     - 將所有人的授權狀態 `status` 一次性全部變更為 `'REVOKED'`。
- **終點**：
  - 商店首頁再也找不到該遊戲。
  - 玩家的購物車若有此遊戲，會被卡住無法結帳。
  - **最重要的是**：已經買過這款遊戲的玩家，他們的「我的遊戲庫」將立刻失去這款遊戲的存取權限。

---

## 4. 全局惡意帳號殲滅 (Account Annihilation)

- **起點**：Admin 對某位屢勸不聽的開發者點擊「永久刪除」。
- **流程**：
  1. 呼叫 `DELETE /api/admin/users/:id`。
  2. 將 `users.permission` 設為 `'DELETED'` (永遠無法登入)。
  3. **自動觸發上述的強制下架機制**：
     - 如果該使用者是 `DEVELOPER`，系統會找出他建立的所有遊戲。
     - 執行 `UPDATE games SET status = 'TAKEN_DOWN'`。
     - 再對這些遊戲執行 `UPDATE game_licenses SET status = 'REVOKED'`。
- **終點**：單次點擊即將該惡劣開發者、其散佈的所有遊戲、以及玩家庫存中的授權徹底淨空，實現平台的最高安全閉環。
