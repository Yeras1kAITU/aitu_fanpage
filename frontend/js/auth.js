class AuthManager {
    constructor() {
        this.currentUser = null;
        this.loadUserFromStorage();
    }

    loadUserFromStorage() {
        const userData = localStorage.getItem('user');
        if (userData) {
            this.currentUser = JSON.parse(userData);
        }
    }

    saveUserToStorage(user) {
        this.currentUser = user;
        localStorage.setItem('user', JSON.stringify(user));
    }

    clearUser() {
        this.currentUser = null;
        localStorage.removeItem('user');
        localStorage.removeItem('token');
    }

    async login(email, password) {
        try {
            const response = await fetch(`${API_BASE}/api/auth/login`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ email, password })
            });

            if (!response.ok) {
                const error = await response.json();
                throw new Error(error.message || 'Login failed');
            }

            const data = await response.json();
            setAuthToken(data.token);
            this.saveUserToStorage(data.user);
            return data.user;
        } catch (error) {
            console.error('Login error:', error);
            throw error;
        }
    }

    async register(userData) {
        try {
            const response = await fetch(`${API_BASE}/api/auth/register`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(userData)
            });

            if (!response.ok) {
                const error = await response.json();
                throw new Error(error.message || 'Registration failed');
            }

            const data = await response.json();
            setAuthToken(data.token);
            this.saveUserToStorage(data.user);
            return data.user;
        } catch (error) {
            console.error('Registration error:', error);
            throw error;
        }
    }

    async logout() {
        try {
            await fetchWithAuth('/api/auth/logout', { method: 'POST' });
        } catch (error) {
        } finally {
            this.clearUser();
            window.location.href = '/index.html';
        }
    }

    async updateProfile(profileData) {
        try {
            const response = await fetchWithAuth('/api/users/me', {
                method: 'PUT',
                body: JSON.stringify(profileData)
            });

            if (!response.ok) {
                throw new Error('Failed to update profile');
            }

            const user = await response.json();
            this.saveUserToStorage(user);
            return user;
        } catch (error) {
            console.error('Update profile error:', error);
            throw error;
        }
    }

    async changePassword(currentPassword, newPassword) {
        try {
            const response = await fetchWithAuth('/api/users/me/password', {
                method: 'PUT',
                body: JSON.stringify({
                    current_password: currentPassword,
                    new_password: newPassword
                })
            });

            if (!response.ok) {
                const error = await response.json();
                throw new Error(error.message || 'Password change failed');
            }

            return true;
        } catch (error) {
            console.error('Change password error:', error);
            throw error;
        }
    }

    isAuthenticated() {
        return !!this.currentUser && !!getAuthToken();
    }

    getUser() {
        return this.currentUser;
    }

    hasRole(role) {
        return this.currentUser?.role === role;
    }

    can(permission) {
        if (!this.currentUser) return false;

        switch(permission) {
            case 'create_post':
                return this.currentUser.is_active;
            case 'manage_posts':
                return ['admin', 'moderator'].includes(this.currentUser.role);
            case 'manage_users':
                return ['admin', 'moderator'].includes(this.currentUser.role);
            case 'view_admin':
                return this.currentUser.role === 'admin';
            default:
                return false;
        }
    }
}

const authManager = new AuthManager();

async function checkAuthStatus() {
    const authLinks = document.getElementById('auth-links');
    const createPostBtn = document.getElementById('create-post-btn');
    const adminLink = document.getElementById('admin-link');
    const userProfileSidebar = document.getElementById('user-profile-sidebar');

    const token = getAuthToken();

    if (token && authManager.isAuthenticated()) {
        const user = authManager.getUser();

        if (authLinks) {
            authLinks.innerHTML = `
                <div class="nav-dropdown">
                    <button class="nav-dropbtn">
                        <img src="${user.profile_image || 'assets/default-avatar.svg'}" 
                             alt="${user.display_name}" class="nav-avatar">
                        ${user.display_name} <i class="fas fa-caret-down"></i>
                    </button>
                    <div class="nav-dropdown-content">
                        <a href="dashboard.html"><i class="fas fa-tachometer-alt"></i> Dashboard</a>
                        <a href="profile.html"><i class="fas fa-user-circle"></i> Profile</a>
                        <hr>
                        <a href="#" id="logout-btn"><i class="fas fa-sign-out-alt"></i> Logout</a>
                    </div>
                </div>
            `;

            document.getElementById('logout-btn')?.addEventListener('click', (e) => {
                e.preventDefault();
                authManager.logout();
            });
        }

        if (createPostBtn && authManager.can('create_post')) {
            createPostBtn.style.display = 'inline-block';
        }

        if (adminLink && authManager.can('view_admin')) {
            adminLink.style.display = 'block';
        }

        if (userProfileSidebar) {
            userProfileSidebar.innerHTML = `
                <div class="sidebar-profile">
                    <img src="${user.profile_image || 'assets/default-avatar.svg'}" 
                         alt="${user.display_name}" class="profile-avatar">
                    <h4>${user.display_name}</h4>
                    <span class="user-role ${user.role}">${getRoleLabel(user.role)}</span>
                    <div class="profile-stats">
                        <div class="stat">
                            <span class="stat-value">${user.post_count || 0}</span>
                            <span class="stat-label">Posts</span>
                        </div>
                        <div class="stat">
                            <span class="stat-value">${user.like_count || 0}</span>
                            <span class="stat-label">Likes</span>
                        </div>
                        <div class="stat">
                            <span class="stat-value">${user.comment_count || 0}</span>
                            <span class="stat-label">Comments</span>
                        </div>
                    </div>
                    <a href="profile.html" class="btn btn-small btn-outline mt-2">
                        <i class="fas fa-edit"></i> Edit Profile
                    </a>
                </div>
            `;
        }
    } else {
        if (authLinks) {
            authLinks.innerHTML = `
                <a href="login.html" class="nav-link"><i class="fas fa-sign-in-alt"></i> Login</a>
                <a href="register.html" class="nav-link"><i class="fas fa-user-plus"></i> Register</a>
            `;
        }

        if (userProfileSidebar) {
            userProfileSidebar.innerHTML = `
                <p>Please login to see your profile</p>
                <a href="login.html" class="btn btn-small">Login</a>
            `;
        }
    }
}

document.addEventListener('DOMContentLoaded', () => {
    checkAuthStatus();
});