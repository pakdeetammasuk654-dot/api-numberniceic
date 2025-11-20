// public/analysis.js

document.addEventListener("DOMContentLoaded", () => {

    // --- 1. Initial Setup ---
    const nameInput = document.getElementById('nameInput');
    const birthDayInput = document.getElementById('birthDayInput');
    const clearBtn = document.getElementById('clearInputBtn');
    const typingStatus = document.getElementById('typingStatus');
    let typingTimer;
    const doneTypingInterval = 1500; // 1.5 วินาที

    // เริ่มต้นฟังก์ชันต่างๆ ครั้งแรก
    initResultFeatures();
    initFormLogic();

    // --- 2. Core Function: Perform Analysis (AJAX) ---
    // ฟังก์ชันนี้จะถูกเรียกเมื่อหยุดพิมพ์ หรือเปลี่ยนวันเกิด หรือคลิกชื่อตัวอย่าง
    async function performAnalysis(overrideName) {
        const name = overrideName || nameInput.value.trim();
        const birthDay = birthDayInput.value;

        if (!name) return;

        typingStatus.style.display = 'block';
        typingStatus.innerText = '🚀 กำลังวิเคราะห์...';

        try {
            // ใช้ Fetch ส่งข้อมูลแบบ POST
            const formData = new FormData();
            formData.append('name', name);
            formData.append('birth_day', birthDay);

            const response = await fetch('/analysis', {
                method: 'POST',
                body: formData
            });

            if (response.ok) {
                const html = await response.text();

                // แปลง HTML String เป็น DOM Element
                const parser = new DOMParser();
                const doc = parser.parseFromString(html, 'text/html');

                // ดึงเฉพาะส่วน #result-wrapper จาก HTML ใหม่
                const newResult = doc.getElementById('result-wrapper');
                const currentResult = document.getElementById('result-wrapper');

                if (newResult && currentResult) {
                    // แทนที่เนื้อหาเดิมด้วยเนื้อหาใหม่ (โดยไม่รีโหลดหน้า)
                    currentResult.innerHTML = newResult.innerHTML;

                    // เรียกฟังก์ชันเพื่อผูก Event Listener ให้กับ Elements ใหม่
                    initResultFeatures();
                }
            }
        } catch (error) {
            console.error("Analysis Error:", error);
        } finally {
            typingStatus.style.display = 'none';
        }
    }

    // ทำให้ฟังก์ชันนี้เรียกใช้ได้จาก HTML (สำหรับ onclick ในตาราง)
    window.analyzeName = function(name) {
        if(nameInput) nameInput.value = name;
        // เลื่อนหน้าจอขึ้นไปที่ฟอร์มนิดนึงเพื่อให้รู้ว่าเปลี่ยนแล้ว
        document.querySelector('.card').scrollIntoView({ behavior: 'smooth' });
        performAnalysis(name);
    };

    // --- 3. Form Logic Setup ---
    function initFormLogic() {
        if (nameInput && clearBtn) {
            const updateBtnState = () => {
                clearBtn.style.display = nameInput.value.length > 0 ? 'block' : 'none';
            };

            // Keyup: นับถอยหลังเมื่อพิมพ์เสร็จ
            nameInput.addEventListener('keyup', () => {
                clearTimeout(typingTimer);
                if (nameInput.value.trim()) {
                    typingStatus.style.display = 'block';
                    typingStatus.innerText = '⏳ กำลังรอพิมพ์เสร็จ...';
                    typingTimer = setTimeout(() => performAnalysis(), doneTypingInterval);
                } else {
                    typingStatus.style.display = 'none';
                }
            });

            // Keydown: หยุดนับถ้ากำลังพิมพ์
            nameInput.addEventListener('keydown', () => {
                clearTimeout(typingTimer);
            });

            // Input Clean & Btn State
            nameInput.addEventListener('input', () => {
                let val = nameInput.value;
                let cleanVal = val.replace(/[^a-zA-Z\u0E00-\u0E7F\s]/g, '').replace(/\s{2,}/g, ' ');
                if (val !== cleanVal) nameInput.value = cleanVal;
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
            updateBtnState();
        }

        // Change Birthday: Trigger Analysis Immediately
        if(birthDayInput) {
            birthDayInput.addEventListener('change', () => {
                if (nameInput.value.trim()) {
                    performAnalysis();
                }
            });
        }

        // Sample Click Handler
        document.querySelectorAll('.sample-item').forEach(item => {
            item.addEventListener('click', () => {
                const name = item.getAttribute('data-name');
                if (nameInput && name) {
                    nameInput.value = name;
                    document.querySelector('.card').scrollIntoView({ behavior: 'smooth' });
                    performAnalysis(name);
                }
            });
        });
    }

    // --- 4. Result Features (ต้องเรียกซ้ำทุกครั้งที่ AJAX โหลดเสร็จ) ---
    function initResultFeatures() {
        // 4.1 Animation Counters
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
                const e = (p === 1) ? 1 : 1 - Math.pow(2, -10 * p);
                counter.innerText = Math.floor(e * (target - start) + start);
                if (p < 1) requestAnimationFrame(update); else counter.innerText = target;
            }
            requestAnimationFrame(update);
        });

        // 4.2 Kakis Highlight Logic
        const dataContainer = document.getElementById('analysisData');
        if (dataContainer) {
            const fullName = dataContainer.getAttribute('data-name');
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

        // 4.3 Detail Section Toggle
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

        // 4.4 AI Linguistics
        const btnLang = document.getElementById('btnLinguistics');
        const modal = document.getElementById('aiModal');
        const closeModal = document.querySelector('.close-modal');
        const aiContent = document.getElementById('aiContent');

        if (btnLang && dataContainer) {
            // ลบ Event Listener เก่า (ถ้ามี) เพื่อกันการกดซ้ำซ้อน (จริงๆ replace innerHTML ก็หายแล้ว แต่กันเหนียว)
            const newBtnLang = btnLang.cloneNode(true);
            btnLang.parentNode.replaceChild(newBtnLang, btnLang);

            newBtnLang.addEventListener('click', async () => {
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

        // Modal Close Logic (Global)
        if (closeModal) closeModal.onclick = () => modal.style.display = "none";
        window.onclick = (event) => { if (event.target == modal) modal.style.display = "none"; }
    }

    // Helper for Color Rendering
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
});