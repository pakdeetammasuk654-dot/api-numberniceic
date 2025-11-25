// public/js/modules/detailToggler.js

export function initializeDetailToggler(btnShowDetails, detailSection, btnCloseDetails) {
    if (btnShowDetails && detailSection && btnCloseDetails) {
        btnShowDetails.addEventListener('click', () => {
            detailSection.style.display = 'block';
            // Scroll to the detail section with a smooth animation
            setTimeout(() => {
                detailSection.scrollIntoView({ behavior: 'smooth', block: 'start' });
            }, 100); // A small delay ensures the element is fully rendered before scrolling
            btnShowDetails.style.display = 'none';
        });
        btnCloseDetails.addEventListener('click', () => {
            detailSection.style.display = 'none';
            btnShowDetails.style.display = 'block';
        });
    }
}
