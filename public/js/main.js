document.addEventListener('DOMContentLoaded', function() {
    const hamburgerButton = document.getElementById('hamburger-button');
    const navMenuWrapper = document.getElementById('nav-menu-wrapper');

    if (hamburgerButton && navMenuWrapper) {
        hamburgerButton.addEventListener('click', function() {
            navMenuWrapper.classList.toggle('active');
        });
    }
});
