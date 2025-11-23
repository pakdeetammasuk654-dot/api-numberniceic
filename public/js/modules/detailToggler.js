// public/js/modules/detailToggler.js

export function initializeDetailToggler(btnShowDetails, detailSection, btnCloseDetails) {
    if (btnShowDetails && detailSection && btnCloseDetails) {
        btnShowDetails.addEventListener('click', () => {
            detailSection.style.display = 'block';
            btnShowDetails.style.display = 'none';
        });
        btnCloseDetails.addEventListener('click', () => {
            detailSection.style.display = 'none';
            btnShowDetails.style.display = 'block';
        });
    }
}
