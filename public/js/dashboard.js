document.addEventListener('DOMContentLoaded', function() {
    // Handle clickable rows
    document.querySelectorAll('.clickable-row').forEach(row => {
        row.addEventListener('click', function() {
            window.location.href = this.dataset.href;
        });
    });

    // Handle delete button confirmation
    document.querySelectorAll('.btn-delete').forEach(button => {
        button.addEventListener('click', function(event) {
            event.stopPropagation(); // Prevent row click event from firing
            if (!confirm('ต้องการลบชื่อนี้ใช่หรือไม่?')) {
                event.preventDefault();
            }
        });
    });
});
