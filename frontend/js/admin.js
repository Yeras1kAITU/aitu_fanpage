class AdminManager {
    constructor() {
        this.currentPage = 1;
        this.usersPerPage = 10;
        this.totalUsers = 0;
        this.currentUser = null;
    }

    async getSystemStats() {
        try {
            const response = await fetchWithAuth('/api/admin/stats');
            return response.json();
        } catch (error) {
            console.error('Get system stats error:', error);
            throw error;
        }
    }

    async getAllUsers(page = 1, search = '') {
        const offset = (page - 1) * this.usersPerPage;
        let url = `/api/admin/users?limit=${this.usersPerPage}&offset=${offset}`;

        if (search) {
            url = `/api/admin/users/search?q=${encodeURIComponent(search)}&limit=${this.usersPerPage}`;
        }

        try {
            const response = await fetchWithAuth(url);
            const users = await response.json();

            // In a real app, the API would return pagination info
            // For now, we'll estimate total users
            if (page === 1) {
                this.totalUsers = users.length === this.usersPerPage ? page * this.usersPerPage + 1 : users.length;
            }

            return users;
        } catch (error) {
            console.error('Get users error:', error);
            throw error;
        }
    }

    async updateUserRole(userId, newRole) {
        try {
            const response = await fetchWithAuth(`/api/admin/users/${userId}/role`, {
                method: 'PUT',
                body: JSON.stringify({ role: newRole })
            });

            if (!response.ok) {
                const error = await response.json();
                throw new Error(error.message || 'Failed to update user role');
            }

            showNotification('User role updated successfully', 'success');
            return true;
        } catch (error) {
            console.error('Update user role error:', error);
            throw error;
        }
    }

    async toggleUserStatus(userId, action) {
        try {
            const response = await fetchWithAuth(`/api/admin/users/${userId}/status/${action}`, {
                method: 'PUT'
            });

            if (!response.ok) {
                const error = await response.json();
                throw new Error(error.message || `Failed to ${action} user`);
            }

            showNotification(`User ${action}d successfully`, 'success');
            return true;
        } catch (error) {
            console.error('Toggle user status error:', error);
            throw error;
        }
    }

    async deleteUser(userId) {
        try {
            const response = await fetchWithAuth(`/api/admin/users/${userId}`, {
                method: 'DELETE'
            });

            if (!response.ok) {
                const error = await response.json();
                throw new Error(error.message || 'Failed to delete user');
            }

            showNotification('User deleted successfully', 'success');
            return true;
        } catch (error) {
            console.error('Delete user error:', error);
            throw error;
        }
    }

    async searchUsers(query) {
        try {
            const response = await fetchWithAuth(`/api/admin/users/search?q=${encodeURIComponent(query)}`);
            return response.json();
        } catch (error) {
            console.error('Search users error:', error);
            throw error;
        }
    }

    async getCategoriesStats() {
        try {
            const response = await fetchWithAuth('/api/posts/categories/stats');
            return response.json();
        } catch (error) {
            console.error('Get categories stats error:', error);
            throw error;
        }
    }

    renderUserRow(user) {
        return `
            <tr data-user-id="${user.id}">
                <td>
                    <div class="user-info">
                        <img src="${user.profile_image || '/assets/default-avatar.svg'}" 
                             alt="${user.display_name}" class="user-avatar">
                        <div>
                            <strong>${user.display_name}</strong>
                            <div class="user-email">${user.email}</div>
                        </div>
                    </div>
                </td>
                <td>
                    <span class="role-badge role-${user.role}">${getRoleLabel(user.role)}</span>
                </td>
                <td>
                    <span class="status-badge ${user.is_active ? 'status-active' : 'status-inactive'}">
                        ${user.is_active ? 'Active' : 'Inactive'}
                    </span>
                </td>
                <td>${user.post_count || 0}</td>
                <td>
                    <div class="table-actions">
                        <button class="action-btn edit-btn" data-user-id="${user.id}">
                            <i class="fas fa-edit"></i>
                        </button>
                        ${user.is_active ?
            `<button class="action-btn toggle-btn deactivate-btn" data-user-id="${user.id}">
                                <i class="fas fa-user-slash"></i>
                            </button>` :
            `<button class="action-btn toggle-btn activate-btn" data-user-id="${user.id}">
                                <i class="fas fa-user-check"></i>
                            </button>`
        }
                        <button class="action-btn delete-btn" data-user-id="${user.id}">
                            <i class="fas fa-trash"></i>
                        </button>
                    </div>
                </td>
            </tr>
        `;
    }

    renderPagination() {
        const totalPages = Math.ceil(this.totalUsers / this.usersPerPage);
        if (totalPages <= 1) return '';

        let pagination = '<div class="pagination">';

        // Previous button
        if (this.currentPage > 1) {
            pagination += `<button class="page-btn" data-page="${this.currentPage - 1}">«</button>`;
        }

        // Page numbers
        for (let i = 1; i <= totalPages; i++) {
            if (i === 1 || i === totalPages || (i >= this.currentPage - 1 && i <= this.currentPage + 1)) {
                pagination += `<button class="page-btn ${i === this.currentPage ? 'active' : ''}" data-page="${i}">${i}</button>`;
            } else if (i === this.currentPage - 2 || i === this.currentPage + 2) {
                pagination += '<span class="page-dots">...</span>';
            }
        }

        // Next button
        if (this.currentPage < totalPages) {
            pagination += `<button class="page-btn" data-page="${this.currentPage + 1}">»</button>`;
        }

        pagination += '</div>';
        return pagination;
    }

    renderCategoryChart(stats) {
        if (!stats || !stats.categories) return '<p>No category data available</p>';

        const maxCount = Math.max(...Object.values(stats.categories));

        return Object.entries(stats.categories)
            .map(([category, count]) => {
                const percentage = (count / maxCount) * 100;
                return `
                    <div class="chart-bar">
                        <div class="chart-label">${getCategoryLabel(category)}</div>
                        <div class="chart-value">
                            <div class="chart-fill" style="width: ${percentage}%"></div>
                        </div>
                        <div class="chart-count">${count}</div>
                    </div>
                `;
            })
            .join('');
    }

    renderRoleChart(users) {
        const roleCounts = {
            student: 0,
            alumni: 0,
            moderator: 0,
            admin: 0
        };

        users.forEach(user => {
            roleCounts[user.role] = (roleCounts[user.role] || 0) + 1;
        });

        const maxCount = Math.max(...Object.values(roleCounts));

        return Object.entries(roleCounts)
            .filter(([role, count]) => count > 0)
            .map(([role, count]) => {
                const percentage = (count / maxCount) * 100;
                return `
                    <div class="chart-bar">
                        <div class="chart-label">${getRoleLabel(role)}</div>
                        <div class="chart-value">
                            <div class="chart-fill" style="width: ${percentage}%; background-color: ${getRoleColor(role)}"></div>
                        </div>
                        <div class="chart-count">${count}</div>
                    </div>
                `;
            })
            .join('');
    }

    renderRecentActivity(activities) {
        // Mock activities
        const mockActivities = [
            {
                type: 'post',
                user: 'John Doe',
                action: 'created a new post',
                title: 'Welcome to AITU 2024!',
                time: '2 hours ago'
            },
            {
                type: 'comment',
                user: 'Jane Smith',
                action: 'commented on',
                title: 'Campus Events This Week',
                time: '4 hours ago'
            },
            {
                type: 'like',
                user: 'Alex Johnson',
                action: 'liked',
                title: 'Study Group Forming',
                time: '6 hours ago'
            },
            {
                type: 'user',
                user: 'System',
                action: 'New user registered',
                title: 'Mike Wilson',
                time: '1 day ago'
            }
        ];

        return mockActivities.map(activity => `
            <li class="activity-item">
                <div class="activity-time">${activity.time}</div>
                <div class="activity-content">
                    <div class="activity-icon activity-${activity.type}">
                        <i class="fas fa-${activity.type === 'post' ? 'newspaper' :
            activity.type === 'comment' ? 'comment' :
                activity.type === 'like' ? 'heart' : 'user'}"></i>
                    </div>
                    <div>
                        <strong>${activity.user}</strong> ${activity.action}
                        ${activity.title ? `<em>"${activity.title}"</em>` : ''}
                    </div>
                </div>
            </li>
        `).join('');
    }
}

const adminManager = new AdminManager();

// Admin panel
document.addEventListener('DOMContentLoaded', async () => {
    await checkAuthStatus();

    if (!authManager.isAuthenticated() || !authManager.hasRole('admin')) {
        showNotification('Access denied. Admin privileges required.', 'error');
        window.location.href = 'index.html';
        return;
    }

    await loadAdminData();
    initAdminEventListeners();
});

async function loadAdminData() {
    try {
        // Load system stats
        const stats = await adminManager.getSystemStats();
        updateStatsDisplay(stats);

        // Load users
        await loadUsers();

        // Load category stats
        const categoryStats = await adminManager.getCategoriesStats();
        renderCategoryChart(categoryStats);

        // Load recent activity
        renderRecentActivity();

    } catch (error) {
        console.error('Error loading admin data:', error);
        showNotification('Failed to load admin data', 'error');
    }
}

function updateStatsDisplay(stats) {
    document.getElementById('total-users').textContent = stats.total_users || 0;
    document.getElementById('total-posts').textContent = stats.total_posts || 0;
    document.getElementById('total-comments').textContent = stats.total_comments || 0;
    document.getElementById('total-likes').textContent = stats.total_likes || 0;
}

async function loadUsers(search = '') {
    try {
        const users = await adminManager.getAllUsers(adminManager.currentPage, search);
        const usersList = document.getElementById('users-list');

        if (users.length === 0) {
            usersList.innerHTML = `
                <tr>
                    <td colspan="5" class="text-center">
                        No users found
                    </td>
                </tr>
            `;
            return;
        }

        usersList.innerHTML = users.map(user => adminManager.renderUserRow(user)).join('');

        // Pagination
        const pagination = document.getElementById('users-pagination');
        pagination.innerHTML = adminManager.renderPagination();

        document.querySelectorAll('.page-btn').forEach(btn => {
            btn.addEventListener('click', async () => {
                const page = parseInt(btn.dataset.page);
                adminManager.currentPage = page;
                await loadUsers(search);
            });
        });

    } catch (error) {
        console.error('Error loading users:', error);
        document.getElementById('users-list').innerHTML = `
            <tr>
                <td colspan="5" class="text-center error">
                    Failed to load users
                </td>
            </tr>
        `;
    }
}

function renderCategoryChart(categoryStats) {
    const chartContainer = document.getElementById('category-chart');
    if (chartContainer) {
        chartContainer.innerHTML = adminManager.renderCategoryChart(categoryStats);
    }
}

function renderRecentActivity() {
    const activityContainer = document.getElementById('recent-activity');
    if (activityContainer) {
        activityContainer.innerHTML = adminManager.renderRecentActivity();
    }
}

function initAdminEventListeners() {
    // User search
    const searchInput = document.getElementById('user-search');
    const searchBtn = document.getElementById('search-users');

    const performSearch = debounce(async () => {
        const query = searchInput.value.trim();
        adminManager.currentPage = 1;
        await loadUsers(query);
    }, 500);

    searchInput.addEventListener('input', performSearch);
    searchBtn.addEventListener('click', performSearch);

    // User action handlers
    document.addEventListener('click', async (e) => {
        // Edit user
        if (e.target.closest('.edit-btn')) {
            const userId = e.target.closest('.edit-btn').dataset.userId;
            await openEditUserModal(userId);
        }

        // Deactivate user
        if (e.target.closest('.deactivate-btn')) {
            const userId = e.target.closest('.deactivate-btn').dataset.userId;
            if (confirm('Are you sure you want to deactivate this user?')) {
                try {
                    await adminManager.toggleUserStatus(userId, 'deactivate');
                    await loadUsers();
                } catch (error) {
                    showNotification(error.message, 'error');
                }
            }
        }

        // Activate user
        if (e.target.closest('.activate-btn')) {
            const userId = e.target.closest('.activate-btn').dataset.userId;
            if (confirm('Are you sure you want to activate this user?')) {
                try {
                    await adminManager.toggleUserStatus(userId, 'activate');
                    await loadUsers();
                } catch (error) {
                    showNotification(error.message, 'error');
                }
            }
        }

        // Delete user
        if (e.target.closest('.delete-btn')) {
            const userId = e.target.closest('.delete-btn').dataset.userId;
            if (confirm('Are you sure you want to delete this user? This action cannot be undone.')) {
                try {
                    await adminManager.deleteUser(userId);
                    await loadUsers();
                } catch (error) {
                    showNotification(error.message, 'error');
                }
            }
        }
    });

    // Modal handling
    const modal = document.getElementById('user-modal');
    const modalClose = document.getElementById('modal-close');
    const modalCancel = document.getElementById('modal-cancel');
    const modalSave = document.getElementById('modal-save');
    const modalDelete = document.getElementById('modal-delete');

    modalClose.addEventListener('click', closeModal);
    modalCancel.addEventListener('click', closeModal);

    modalSave.addEventListener('click', async () => {
        const userId = document.getElementById('edit-user-id').value;
        const role = document.getElementById('edit-role').value;
        const status = document.getElementById('edit-status').value;

        try {
            // Update role
            await adminManager.updateUserRole(userId, role);

            // Update status if changed
            if (status === 'active') {
                await adminManager.toggleUserStatus(userId, 'activate');
            } else {
                await adminManager.toggleUserStatus(userId, 'deactivate');
            }

            closeModal();
            await loadUsers();

        } catch (error) {
            showNotification(error.message, 'error');
        }
    });

    modalDelete.addEventListener('click', async () => {
        const userId = document.getElementById('edit-user-id').value;

        if (confirm('Are you sure you want to delete this user? This action cannot be undone.')) {
            try {
                await adminManager.deleteUser(userId);
                closeModal();
                await loadUsers();
            } catch (error) {
                showNotification(error.message, 'error');
            }
        }
    });

    // Close modal on background click
    modal.addEventListener('click', (e) => {
        if (e.target === modal) {
            closeModal();
        }
    });

    // Save settings
    document.getElementById('save-settings').addEventListener('click', () => {
        const maintenanceMode = document.getElementById('site-maintenance').checked;
        const allowRegistration = document.getElementById('new-user-registration').checked;
        const requireApproval = document.getElementById('post-approval').checked;

        // In a real app, this would save to the backend
        showNotification('Settings saved successfully', 'success');
    });
}

async function openEditUserModal(userId) {
    try {
        const row = document.querySelector(`[data-user-id="${userId}"]`);
        const userEmail = row.querySelector('.user-email').textContent;
        const userName = row.querySelector('strong').textContent;
        const userRole = row.querySelector('.role-badge').textContent.toLowerCase();
        const userStatus = row.querySelector('.status-badge').textContent === 'Active' ? 'active' : 'inactive';

        document.getElementById('edit-user-id').value = userId;
        document.getElementById('edit-display-name').value = userName;
        document.getElementById('edit-email').value = userEmail;
        document.getElementById('edit-role').value = userRole;
        document.getElementById('edit-status').value = userStatus;

        document.getElementById('user-modal').style.display = 'flex';

    } catch (error) {
        console.error('Error opening edit modal:', error);
        showNotification('Failed to load user details', 'error');
    }
}

function closeModal() {
    document.getElementById('user-modal').style.display = 'none';
}