// public/js/admin.js

document.addEventListener('DOMContentLoaded', function() {
    // --- Admin Panel Links ---
    const manageUsersLink = document.getElementById('manage-users-link');
    if (manageUsersLink) {
        manageUsersLink.addEventListener('click', function(event) {
            event.preventDefault();
            alert('ระบบจัดการผู้ใช้จะเปิดให้ใช้งานเร็วๆ นี้!');
        });
    }

    const viewStatsLink = document.getElementById('view-stats-link');
    if (viewStatsLink) {
        viewStatsLink.addEventListener('click', function(event) {
            event.preventDefault();
            alert('ระบบสถิติจะเปิดให้ใช้งานเร็วๆ นี้!');
        });
    }

    // --- Types Page Delete Confirmation ---
    const deleteButtons = document.querySelectorAll('.delete-type-btn');
    deleteButtons.forEach(function(button) {
        button.addEventListener('click', function(event) {
            if (!confirm('คุณแน่ใจหรือไม่ที่จะลบประเภทนี้?')) {
                event.preventDefault();
            }
        });
    });
});
