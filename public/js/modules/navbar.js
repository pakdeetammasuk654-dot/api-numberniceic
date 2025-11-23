// public/js/modules/navbar.js

export function initializeNavbar() {
    const hamburgerButton = document.getElementById('hamburger-button');
    const navMenuWrapper = document.getElementById('nav-menu-wrapper');

    if (hamburgerButton && navMenuWrapper) {
        hamburgerButton.addEventListener('click', function() {
            navMenuWrapper.classList.toggle('active');
        });
    }
}
