// public/analysis.js

document.addEventListener("DOMContentLoaded", () => {

    // 1. Animation: Score Counters
    const counters = [
        document.getElementById('totalScore'),
        document.getElementById('goodScore'),
        document.getElementById('badScore')
    ];

    counters.forEach(counter => {
        if (!counter) return;
        const target = parseInt(counter.getAttribute('data-target'), 10) || 0;
        let start = 0;
        const duration = 2000;
        const startTime = performance.now();

        function update(t) {
            const p = Math.min((t - startTime) / duration, 1);
            const e = (p === 1) ? 1 : 1 - Math.pow(2, -10 * p); // Ease out expo
            const val = Math.floor(e * (target - start) + start);
            counter.innerText = val;
            if (p < 1) requestAnimationFrame(update);
            else counter.innerText = target;
        }
        requestAnimationFrame(update);
    });

    // 2. Form Logic: Debounce Auto Submit & Clear Button
    const nameInput = document.getElementById('nameInput');
    const clearBtn = document.getElementById('clearInputBtn');
    const typingStatus = document.getElementById('typingStatus');
    const form = document.querySelector('.analysis-form');
    let typingTimer;
    const doneTypingInterval = 1500; // 1.5 วินาที

    if (nameInput && clearBtn) {
        const updateBtnState = () => {
            clearBtn.style.display = nameInput.value.length > 0 ? 'block' : 'none';
        };

        // Keyup: เริ่มนับถอยหลังเมื่อพิมพ์เสร็จ
        nameInput.addEventListener('keyup', () => {
            clearTimeout(typingTimer);
            if (nameInput.value.trim()) {
                typingStatus.style.display = 'block';
                typingTimer = setTimeout(() => {
                    typingStatus.innerText = '🚀 กำลังวิเคราะห์...';
                    form.submit();
                }, doneTypingInterval);
            } else {
                typingStatus.style.display = 'none';
            }
        });

        // Keydown: หยุดนับถ้ากำลังพิมพ์อยู่
        nameInput.addEventListener('keydown', () => {
            clearTimeout(typingTimer);
            typingStatus.style.display = 'none';
        });

        // Input: Validation & Icon Toggle
        nameInput.addEventListener('input', () => {
            // ให้พิมพ์ได้เฉพาะ ไทย, อังกฤษ, ช่องว่าง
            let val = nameInput.value;
            let cleanVal = val.replace(/[^a-zA-Z\u0E00-\u0E7F\s]/g, '');
            cleanVal = cleanVal.replace(/\s{2,}/g, ' '); // ลดช่องว่างซ้ำ

            if (val !== cleanVal) {
                nameInput.value = cleanVal;
            }
            updateBtnState();
        });

        // Clear Button
        clearBtn.addEventListener('click', () => {
            nameInput.value = '';
            updateBtnState();
            nameInput.focus();
            clearTimeout(typingTimer);
            typingStatus.style.display = 'none';
        });

        // Init State
        updateBtnState();
    }

    // 3. Kakis Highlight Logic (อักษรสีแดง)
    // อ่านข้อมูลจาก Hidden Div ที่เราฝังไว้ใน HTML
    const dataContainer = document.getElementById('analysisData');

    if (dataContainer) {
        const fullName = dataContainer.getAttribute('data-name');
        // แปลง string "ก,ข,ค" เป็น array ["ก", "ข", "ค"]
        const kakisString = dataContainer.getAttribute('data-kakis');
        const badChars = kakisString ? kakisString.split(',') : [];

        const sunEl = document.getElementById('sunNameDisplay');
        const similarEl = document.getElementById('similarNameDisplay');

        if (fullName) {
            const coloredHtml = renderColoredName(fullName, badChars);
            if (sunEl) sunEl.innerHTML = coloredHtml;
            if (similarEl) similarEl.innerHTML = coloredHtml;
        }
    }

    function renderColoredName(name, badChars) {
        if (!name) return "";
        let html = "";
        for (let c of name) {
            if (badChars.includes(c)) {
                html += `<span class="bad-char">${c}</span>`;
            } else {
                html += `<span class="good-char">${c}</span>`;
            }
        }
        return html;
    }

    // 4. Sample Click Handler (คลิกชื่อตัวอย่าง)
    document.querySelectorAll('.sample-item').forEach(item => {
        item.addEventListener('click', () => {
            const name = item.getAttribute('data-name');
            if (nameInput && name) {
                nameInput.value = name;
                document.querySelector('.card').scrollIntoView({ behavior: 'smooth' });
                form.submit();
            }
        });
    });

    // 5. AI Linguistics Modal
    const btnLang = document.getElementById('btnLinguistics');
    const modal = document.getElementById('aiModal');
    const closeModal = document.querySelector('.close-modal');
    const aiContent = document.getElementById('aiContent');

    if (btnLang && dataContainer) {
        btnLang.addEventListener('click', async () => {
            const currentName = dataContainer.getAttribute('data-name');
            if (!currentName) return;

            modal.style.display = "block";
            aiContent.innerHTML = '<div class="ai-loading">⏳ กำลังสอบถาม Gemini AI...<br><small>โปรดรอสักครู่</small></div>';

            try {
                const response = await fetch(`/api/linguistics?name=${encodeURIComponent(currentName)}`);
                const data = await response.json();
                if (response.ok) {
                    aiContent.innerHTML = data.text.replace(/\n/g, '<br>');
                } else {
                    aiContent.innerHTML = `<div style="color:red;">⚠️ เกิดข้อผิดพลาด: ${data.error || 'Unknown error'}</div>`;
                }
            } catch (error) {
                aiContent.innerHTML = `<div style="color:red;">⚠️ ไม่สามารถเชื่อมต่อกับเซิร์ฟเวอร์ได้</div>`;
            }
        });
    }

    if (closeModal) closeModal.onclick = () => modal.style.display = "none";
    window.onclick = (event) => { if (event.target == modal) modal.style.display = "none"; }

    // 6. Detail Section Toggle
    const btnShowDetails = document.getElementById('btnShowDetails');
    const detailSection = document.getElementById('detailSection');
    const btnCloseDetails = document.getElementById('btnCloseDetails');

    if (btnShowDetails && detailSection) {
        btnShowDetails.addEventListener('click', () => {
            detailSection.style.display = 'block';
            detailSection.scrollIntoView({ behavior: 'smooth', block: 'start' });
        });
    }
    if (btnCloseDetails && detailSection) {
        btnCloseDetails.addEventListener('click', () => {
            detailSection.style.display = 'none';
        });
    }
});