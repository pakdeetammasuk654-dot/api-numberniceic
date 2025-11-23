// public/js/modules/scoreAnimator.js

export function initializeScoreAnimator(analysisData, sunNameDisplay) {
    function animateValue(id) {
        const obj = document.getElementById(id);
        if (!obj) return;

        const end = parseInt(obj.dataset.target, 10) || 0;
        const duration = 1000;
        const start = 0;
        let startTimestamp = null;

        const step = (timestamp) => {
            if (!startTimestamp) startTimestamp = timestamp;
            const progress = Math.min((timestamp - startTimestamp) / duration, 1);
            obj.innerHTML = Math.floor(progress * (end - start) + start);
            if (progress < 1) {
                window.requestAnimationFrame(step);
            }
        };
        window.requestAnimationFrame(step);
    }

    if (analysisData) {
        const name = analysisData.dataset.name;
        const kakis = analysisData.dataset.kakis ? analysisData.dataset.kakis.split(',') : [];

        if (sunNameDisplay) {
            let nameHtml = '';
            for (const char of name) {
                if (kakis.includes(char)) {
                    nameHtml += `<span class="bad-char">${char}</span>`;
                } else {
                    nameHtml += char;
                }
            }
            sunNameDisplay.innerHTML = nameHtml;
        }

        // Simplified and corrected animation calls
        animateValue("totalScore");
        animateValue("goodScore");
        animateValue("badScore");
    }
}
