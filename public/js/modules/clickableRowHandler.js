// public/js/modules/clickableRowHandler.js

export function initializeClickableRows() {
    document.querySelectorAll('.clickable-row').forEach(row => {
        row.addEventListener('click', function() {
            if (this.dataset.href) {
                window.location.href = this.dataset.href;
            }
        });
    });
}
