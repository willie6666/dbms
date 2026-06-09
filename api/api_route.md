# API 規格與後端伺服器 (Go Server) 路由稽核報告

經過詳細的代碼追蹤與交叉比對，我將 API 規格書中的 53 支 API，與實際運行在 Go 後端伺服器的 `backend/routes/routes.go` 進行了一對一的比對，確認實作狀態。

## 📝 實作細節盤點

### 1. 公開路由 (Public Routes) - 免 JWT
| API 規格 | `routes.go` 註冊狀態 | 驗證 |
| :--- | :--- | :---: |
| `[POST] /api/auth/register` | `auth.POST("/register")` | ✅ |
| `[POST] /api/auth/login` | `auth.POST("/login")` | ✅ |
| `[GET] /api/games` | `api.GET("/games")` | ✅ |
| `[GET] /api/games/{id}` | `api.GET("/games/:id")` | ✅ |
| `[GET] /api/games/{id}/reviews` | `api.GET("/games/:id/reviews")` | ✅ |
| `[GET] /api/tags` | `api.GET("/tags")` | ✅ |

### 2. 一般登入保護 (Protected Routes) - 僅需 JWT
*這包含購物車、玩家遊戲庫、個人檔案，以及社群互動等。*
| API 規格 | `routes.go` 註冊狀態 | 驗證 |
| :--- | :--- | :---: |
| `[POST] /api/auth/logout` | `auth.POST("/logout")` | ✅ |
| `[PUT] /api/users/profile` | `users.PUT("/profile")` | ✅ |
| `[GET] /api/protected/cart` | `protected.GET("/cart")` | ✅ |
| `[POST] /api/protected/cart` | `protected.POST("/cart")` | ✅ |
| `[DELETE] /api/protected/cart/{id}` | `protected.DELETE("/cart/:game_id")` | ✅ |
| `[POST] /api/protected/checkout` | `protected.POST("/checkout")` | ✅ |
| `[GET] /api/protected/transactions` | `protected.GET("/transactions")` | ✅ |
| `[GET] /api/protected/refunds` | `protected.GET("/refunds")` | ✅ |
| `[GET] /api/protected/library` | `protected.GET("/library")` | ✅ |
| `[GET] /api/protected/wishlist` | `protected.GET("/wishlist")` | ✅ |
| `[POST] /api/protected/wishlist` | `protected.POST("/wishlist")` | ✅ |
| `[DELETE] /api/protected/wishlist/{id}` | `protected.DELETE("/wishlist/:game_id")` | ✅ |
| `[GET] /api/protected/library/.../play` | `protected.GET("/library/:game_id/play")` | ✅ |
| `[GET] /api/protected/library/.../download` | `protected.GET("/library/:game_id/download")` | ✅ |
| `[POST] /api/social/.../reviews` | `social.POST("/games/:id/reviews")` | ✅ |
| `[POST] /api/social/.../replies` | `social.POST("/reviews/:review_id/replies")` | ✅ |
| `[DELETE] /api/social/.../replies/{id}`| `social.DELETE("/reviews/replies/:reply_id")` | ✅ |
| `[POST] /api/social/refunds` | `social.POST("/refunds")` | ✅ |
| `[GET] /api/social/friends` | `social.GET("/friends")` | ✅ |
| `[GET] /api/social/friends/requests` | `social.GET("/friends/requests")` | ✅ |
| `[POST] /api/social/friends/request` | `social.POST("/friends/request")` | ✅ |
| `[PUT] /api/social/friends/.../accept` | `social.PUT("/friends/request/:id/accept")` | ✅ |
| `[PUT] /api/social/friends/.../decline`| `social.PUT("/friends/request/:id/decline")`| ✅ |
| `[DELETE] /api/social/friends/request/{id}` | `social.DELETE("/friends/request/:id")` | ✅ |
| `[GET] /api/social/blacklist` | `social.GET("/blacklist")` | ✅ |
| `[POST] /api/social/blacklist` | `social.POST("/blacklist")` | ✅ |
| `[DELETE] /api/social/blacklist/{id}` | `social.DELETE("/blacklist/:user_id")` | ✅ |
| `[POST] /api/social/messages` | `social.POST("/messages")` | ✅ |
| `[GET] /api/social/messages/{user_id}` | `social.GET("/messages/:user_id")` | ✅ |

### 3. 開發者路由 (Developer Routes) - 需 DEVELOPER 權限
| API 規格 | `routes.go` 註冊狀態 | 驗證 |
| :--- | :--- | :---: |
| `[GET] /api/developer/games` | `developer.GET("/games")` | ✅ |
| `[POST] /api/developer/games` | `developer.POST("/games")` | ✅ |
| `[PUT] /api/developer/games/{id}/publish` | `developer.PUT("/games/:id/publish")` | ✅ |
| `[PUT] /api/developer/games/{id}` | `developer.PUT("/games/:id")` | ✅ |
| `[DELETE] /api/developer/games/{id}` | `developer.DELETE("/games/:id")` | ✅ |
| `[POST] /api/developer/games/{id}/media` | `developer.POST("/games/:id/media")` | ✅ |
| `[DELETE] /api/developer/.../media/{id}` | `developer.DELETE("/games/:id/media/:media_id")` | ✅ |
| `[GET] /api/developer/games/{id}/stats` | `developer.GET("/games/:id/stats")` | ✅ |
| `[POST] /api/developer/tags` | `developer.POST("/tags")` | ✅ |
| `[POST] /api/developer/games/{id}/tags` | `developer.POST("/games/:id/tags")` | ✅ |
| `[DELETE] /api/developer/.../tags/{id}` | `developer.DELETE("/games/:id/tags/:tag_id")` | ✅ |

### 4. 管理員與客服路由 (Admin & CSR Routes)
| API 規格 | `routes.go` 註冊狀態 | 驗證 |
| :--- | :--- | :---: |
| `[GET] /api/admin/users` | `admin.GET("/users")` | ✅ |
| `[PUT] /api/admin/users/{id}/suspend` | `admin.PUT("/users/:id/suspend")` | ✅ |
| `[DELETE] /api/admin/users/{id}` | `admin.DELETE("/users/:id")` | ✅ |
| `[PUT] /api/admin/users/{id}/role` | `admin.PUT("/users/:id/role")` | ✅ |
| `[DELETE] /api/admin/games/{id}` | `admin.DELETE("/games/:id")` | ✅ |
| `[GET] /api/csr/refunds` | `csr.GET("/refunds")` | ✅ |
| `[PUT] /api/csr/refunds/{id}` | `csr.PUT("/refunds/:id")` | ✅ |

---
> [!NOTE]
> **補充說明**
> 1. API 文件中的路徑參數採用標準的 `{id}` 寫法，在 Go/Gin 框架中實作為 `:id`（如 `:game_id` 或 `:user_id`），此為框架約定俗成的語法對應，邏輯上等同於文件描述。
> 2. `auth`, `admin`, `csr`, `developer`, `protected`, `social` 各大 Router Group 皆已正確綁定對應的 `RequireRole()` 或 `RequireAuth()` 中介軟體防護。
