// public/js/modules/solarSystemAnimator.js

export function initializeSolarSystemAnimator() {
    const summaryBox = document.querySelector('.summary-box');
    if (!summaryBox) return;

    // Function to trigger the animation
    const triggerAnimation = () => {
        summaryBox.classList.remove('animate-fade-in');
        // Void-read to force reflow, ensuring the class removal is processed before adding it back
        void summaryBox.offsetWidth;
        summaryBox.classList.add('animate-fade-in');
    };

    // Initial animation on page load
    triggerAnimation();

    // Listen for a custom event that signals a recalculation
    document.addEventListener('recalculation', triggerAnimation);
}
