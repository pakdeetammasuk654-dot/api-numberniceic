document.addEventListener('DOMContentLoaded', function () {
    // --- Desktop View ---
    const desktopRows = document.querySelectorAll('.dashboard-desktop-view .clickable-row');
    desktopRows.forEach(row => {
        const deleteButton = row.querySelector('.btn-delete');
        const analysisLink = row.dataset.href;

        row.addEventListener('click', (e) => {
            if (deleteButton && deleteButton.contains(e.target)) {
                return;
            }
            if (analysisLink) {
                window.location.href = analysisLink;
            }
        });

        if (deleteButton) {
            deleteButton.addEventListener('click', (e) => {
                e.stopPropagation();
                e.preventDefault();
                const confirmation = window.confirm('คุณแน่ใจหรือไม่ว่าต้องการลบชื่อนี้?');
                if (confirmation) {
                    window.location.href = deleteButton.href;
                }
            });
        }
    });

    // --- Mobile View ---
    const mobileCards = document.querySelectorAll('.dashboard-mobile-view .clickable-row');
    mobileCards.forEach(card => {
        const deleteButton = card.querySelector('.btn-delete');
        const analysisLink = card.dataset.href;

        // Handle card click for analysis
        card.addEventListener('click', (e) => {
            if (deleteButton && deleteButton.contains(e.target)) {
                return;
            }
            if (analysisLink) {
                window.location.href = analysisLink;
            }
        });

        // Handle delete confirmation
        if (deleteButton) {
            deleteButton.addEventListener('click', (e) => {
                e.stopPropagation();
                e.preventDefault();
                const confirmation = window.confirm('คุณแน่ใจหรือไม่ว่าต้องการลบชื่อนี้?');
                if (confirmation) {
                    window.location.href = deleteButton.href;
                }
            });
        }
    });
});
