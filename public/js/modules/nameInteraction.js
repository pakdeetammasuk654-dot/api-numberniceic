// public/js/modules/nameInteraction.js

export function initializeNameInteraction(birthDayInput) {
    const analyzeName = function(name) {
        const birthDay = birthDayInput ? birthDayInput.value : '';
        window.location.href = `/analysis?name=${encodeURIComponent(name)}&birth_day=${birthDay}`;
    };

    // Expose analyzeName globally for now if it's called from HTML directly
    // Consider refactoring HTML to call through a module function if possible
    window.analyzeName = analyzeName;

    // --- Sample Names ---
    document.querySelectorAll('.sample-item').forEach(item => {
        item.addEventListener('click', () => {
            analyzeName(item.dataset.name);
        });
    });

    // --- Similar Names Table ---
    document.querySelectorAll('.similar-name-row').forEach(row => {
        row.addEventListener('click', () => {
            analyzeName(row.dataset.name);
        });
        row.style.cursor = 'pointer'; // Add pointer cursor to indicate it's clickable
    });
}
