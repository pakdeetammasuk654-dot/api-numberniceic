// public/js/modules/deleteConfirmationHandler.js

export function initializeDeleteConfirmation() {
    document.querySelectorAll('.btn-delete').forEach(button => {
        button.addEventListener('click', function(event) {
            event.stopPropagation(); // Prevent row click event from firing
            if (!confirm('ต้องการลบชื่อนี้ใช่หรือไม่?')) {
                event.preventDefault();
            }
        });
    });
}
