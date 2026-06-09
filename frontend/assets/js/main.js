// ============================================================
// VaporAuror — main.js
// 全域狀態與共用元件
// ============================================================

// ============================================================
// 頁面初始化
// ============================================================
document.documentElement.setAttribute('data-theme', 'dark');

document.addEventListener("DOMContentLoaded", () => {
    renderHeader();

    // 搜尋框事件（首頁）
    const searchForm = document.getElementById('search-form');
    if (searchForm) {
        searchForm.addEventListener('submit', (e) => {
            e.preventDefault();
            const keyword = document.getElementById('search-input').value.trim();
            if (keyword) {
                window.location.href = `/pages/store/search.html?q=${encodeURIComponent(keyword)}`;
            }
        });
    }
});

document.addEventListener('click', (event) => {
    const burger = event.target.closest('.navbar-burger');
    if (!burger) return;
    const target = document.getElementById(burger.dataset.target);
    if (!target) return;
    burger.classList.toggle('is-active');
    target.classList.toggle('is-active');
});

// ============================================================
// 渲染共用導覽列
// ============================================================
function renderHeader() {
    const header = document.getElementById('global-header');
    if (!header) return;

    // 登入頁與註冊頁只顯示 Logo
    const path = window.location.pathname;
    if (path.includes('login.html') || path.includes('register.html')) {
        header.innerHTML = `
            <nav class="navbar va-navbar" role="navigation" aria-label="main navigation">
                <div class="navbar-brand">
                    <a href="/" class="navbar-item logo">VaporAuror</a>
                </div>
            </nav>`;
        return;
    }

    const currentRole = localStorage.getItem('userRole') || 'GUEST';
    const userDataStr = localStorage.getItem('currentUser');
    const username = userDataStr ? JSON.parse(userDataStr).username : '玩家';

    let navItems = `<a class="navbar-item" href="/">商店首頁</a>`;

    if (currentRole !== 'GUEST') {
        navItems += `
            <a class="navbar-item" href="/pages/user/library">遊戲庫</a>
            <a class="navbar-item" href="/pages/user/social">社群</a>
        `;
        if (currentRole === 'DEVELOPER') {
            navItems += `<a class="navbar-item has-text-warning" href="/pages/dashboard/dev_dashboard">開發者中心</a>`;
        } else if (currentRole === 'CSR') {
            navItems += `<a class="navbar-item has-text-info" href="/pages/dashboard/csr_dashboard">客服中心</a>`;
        } else if (currentRole === 'ADMIN') {
            navItems += `<a class="navbar-item has-text-info" href="/pages/dashboard/csr_dashboard">客服中心</a>`;
            navItems += `<a class="navbar-item has-text-danger" href="/pages/dashboard/admin_dashboard">管理後台</a>`;
        }
    }

    let accountHtml = '';
    if (currentRole === 'GUEST') {
        accountHtml += `
            <div class="navbar-item">
                <div class="buttons">
                    <a href="/pages/auth/login" class="button is-primary"><strong>登入 / 註冊</strong></a>
                </div>
            </div>`;
    } else {
        accountHtml += `
            <a href="/pages/user/cart" class="navbar-item">購物車</a>
            <div class="navbar-item has-dropdown is-hoverable">
                <a class="navbar-link">${escapeHtml(username)}</a>
                <div class="navbar-dropdown is-right">
                    <a class="navbar-item" href="/pages/user/profile">基本資料</a>
                    <a class="navbar-item" href="/pages/user/history">購買紀錄</a>
                    <a class="navbar-item" href="/pages/user/refund_request">退款申請</a>
                    <a class="navbar-item" href="/pages/user/wishlist">願望清單</a>
                    <hr class="navbar-divider">
                    <a class="navbar-item" href="#" id="logout-btn">登出</a>
                </div>
            </div>
        `;
    }

    header.innerHTML = `
        <nav class="navbar va-navbar" role="navigation" aria-label="main navigation">
            <div class="navbar-brand">
                <a href="/" class="navbar-item logo">VaporAuror</a>
                <a role="button" class="navbar-burger" aria-label="menu" aria-expanded="false" data-target="global-navbar-menu">
                    <span aria-hidden="true"></span>
                    <span aria-hidden="true"></span>
                    <span aria-hidden="true"></span>
                    <span aria-hidden="true"></span>
                </a>
            </div>
            <div id="global-navbar-menu" class="navbar-menu">
                <div class="navbar-start">${navItems}</div>
                <div class="navbar-end">${accountHtml}</div>
            </div>
        </nav>`;

    // 登出事件
    const logoutBtn = document.getElementById('logout-btn');
    if (logoutBtn) {
        logoutBtn.addEventListener('click', async (e) => {
            e.preventDefault();
            try {
                await apiLogout();
            } catch (_) {
                // 即使 API 失敗也清除本地狀態
                localStorage.removeItem('token');
                localStorage.removeItem('currentUser');
                localStorage.removeItem('userRole');
            }
            window.location.href = '/';
        });
    }
}

// ============================================================
// 渲染遊戲卡片列表
// ============================================================
function renderGames(games) {
    const container = document.getElementById('game-list');
    if (!container) return;

    if (!games) {
        // Skeleton Loading 骨架屏狀態
        container.innerHTML = '';
        for (let i = 0; i < 4; i++) {
            const skeleton = document.createElement('div');
            skeleton.className = 'game-list-card card';
            skeleton.innerHTML = `
                <div class="game-thumbnail skeleton"></div>
                <div class="game-list-info">
                    <div class="skeleton skeleton-title"></div>
                    <div class="skeleton skeleton-text"></div>
                    <div class="skeleton skeleton-text" style="width: 80%;"></div>
                    <div class="game-list-tags" style="margin-top: 10px;">
                        <span class="tag skeleton" style="width: 40px; height: 16px;"></span>
                        <span class="tag skeleton" style="width: 50px; height: 16px;"></span>
                    </div>
                </div>
            `;
            container.appendChild(skeleton);
        }
        return;
    }

    if (games.length === 0) {
        container.innerHTML = `
            <div class="empty-state">
                <p>找不到符合條件的遊戲</p>
                <a href="/" class="button is-light" style="margin-top:10px;">回首頁</a>
            </div>`;
        return;
    }

    container.innerHTML = '';
    games.forEach(game => {
        const card = document.createElement('div');
        card.className = 'game-list-card card';

        const tagsHtml = (game.tags || []).map(t => {
            const tagName = typeof t === 'string' ? t : (t.tag_name || t.name || t);
            return `<span class="tag is-rounded">${escapeHtml(String(tagName))}</span>`;
        }).join('');
        const priceHtml = game.price === 0
            ? `<span class="game-list-price free">免費</span>`
            : `<span class="game-list-price">NT$ ${game.price.toLocaleString()}</span>`;

        let coverHtml = `<span style="font-size:13px;color:#8f98a0;">[${escapeHtml(game.title)} 封面]</span>`;
        if (game.media && game.media.length > 0) {
            const thumb = game.media.find(m => m.media_type === 'thumbnail') || game.media.find(m => m.media_type === 'media');
            if (thumb) {
                coverHtml = `<img src="${thumb.file_url}" alt="cover" style="width:100%; height:100%; object-fit:cover; border-radius: 8px;">`;
            }
        }

        card.innerHTML = `
            <div class="card-image game-thumbnail">${coverHtml}</div>
            <div class="card-content game-list-info">
                <p class="title is-5 game-list-title">${escapeHtml(game.title)}</p>
                <p class="content game-list-desc">${escapeHtml(game.desc || '')}</p>
                <div class="game-list-tags tags">${tagsHtml}</div>
                <p>${priceHtml}</p>
            </div>
        `;

        card.addEventListener('click', () => {
            const gameId = game.game_id || game.id;
            window.location.href = `/pages/store/game_detail?id=${gameId}`;
        });

        container.appendChild(card);
    });
}

// ============================================================
// 工具函式
// ============================================================
function escapeHtml(str) {
    if (typeof str !== 'string') return '';
    return str
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#39;');
}

function getCurrentUser() {
    const str = localStorage.getItem('currentUser');
    return str ? JSON.parse(str) : null;
}

function getCurrentRole() {
    return localStorage.getItem('userRole') || 'GUEST';
}

function getGameCoverUrl(game) {
    const media = (game && game.media) || (game && game.Media) || [];
    const image = media.find(m => m.media_type === 'thumbnail') || media.find(m => m.media_type === 'media');
    return image ? image.file_url : '';
}

function renderMarkdown(markdown) {
    const source = String(markdown || '').trim();
    if (!source) return '<p class="has-text-grey">尚未填寫介紹</p>';
    if (!window.marked || !window.DOMPurify) {
        return `<p>${escapeHtml(source).replace(/\n/g, '<br>')}</p>`;
    }
    return DOMPurify.sanitize(marked.parse(source));
}

function requireLogin(redirectTo = '/pages/auth/login.html') {
    if (!localStorage.getItem('token')) {
        alert('請先登入！');
        window.location.href = redirectTo;
        return false;
    }
    return true;
}

// ============================================================
// 全域 Toast Notification 系統
// type: 'info', 'success', 'error'
// ============================================================
function showToast(message, type = 'info') {
    let container = document.getElementById('toast-container');
    if (!container) {
        container = document.createElement('div');
        container.id = 'toast-container';
        document.body.appendChild(container);
    }

    const toast = document.createElement('div');
    toast.className = `toast-message toast-${type}`;
    
    toast.innerHTML = `<span>${escapeHtml(message)}</span>`;
    container.appendChild(toast);

    // Animate in
    setTimeout(() => toast.classList.add('show'), 10);

    // Auto remove after 3s
    setTimeout(() => {
        toast.classList.remove('show');
        setTimeout(() => toast.remove(), 400); // Wait for transition
    }, 3000);
}
