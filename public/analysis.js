// public/analysis.js

document.addEventListener("DOMContentLoaded", () => {

    // --- 1. Initial Setup ---
    const nameInput = document.getElementById('nameInput');
    const birthDayInput = document.getElementById('birthDayInput');
    const clearBtn = document.getElementById('clearInputBtn');
    const typingStatus = document.getElementById('typingStatus');
    let typingTimer;
    const doneTypingInterval = 1500;

    // เรียกใช้ครั้งแรกตอนโหลดหน้าเว็บ
    initResultFeatures();
    initFormLogic();

    // --- 2. Core Function: Perform Analysis (AJAX) ---
    async function performAnalysis(overrideName) {
        const name = overrideName || nameInput.value.trim();
        const birthDay = birthDayInput.value;

        if (!name) return;

        typingStatus.style.display = 'block';
        typingStatus.innerText = '🚀 กำลังวิเคราะห์...';

        try {
            const formData = new FormData();
            formData.append('name', name);
            formData.append('birth_day', birthDay);

            const response = await fetch('/analysis', {
                method: 'POST',
                body: formData
            });

            if (response.ok) {
                const html = await response.text();
                const parser = new DOMParser();
                const doc = parser.parseFromString(html, 'text/html');

                // ดึงเนื้อหาใหม่มาแทนที่
                const newResult = doc.getElementById('result-wrapper');
                const currentResult = document.getElementById('result-wrapper');

                if (newResult && currentResult) {
                    currentResult.innerHTML = newResult.innerHTML;

                    // 🔥 สำคัญมาก: ต้องเรียกฟังก์ชันนี้ใหม่ทุกครั้งหลังโหลด AJAX เสร็จ
                    // เพื่อให้ปุ่ม "บันทึกชื่อ" (และกราฟอื่นๆ) กลับมาทำงาน
                    initResultFeatures();
                }
            }
        } catch (error) {
            console.error("Analysis Error:", error);
        } finally {
            typingStatus.style.display = 'none';
        }
    }

    window.analyzeName = function(name) {
        if(nameInput) nameInput.value = name;
        document.querySelector('.card').scrollIntoView({ behavior: 'smooth' });
        performAnalysis(name);
    };

    // --- 3. Form Logic Setup ---
    function initFormLogic() {
        if (nameInput && clearBtn) {
            const updateBtnState = () => {
                clearBtn.style.display = nameInput.value.length > 0 ? 'block' : 'none';
            };
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
            nameInput.addEventListener('keydown', () => { clearTimeout(typingTimer); });
            nameInput.addEventListener('input', () => {
                // กรองให้พิมพ์ได้เฉพาะภาษาไทย อังกฤษ และช่องว่าง
                let val = nameInput.value;
                let cleanVal = val.replace(/[^a-zA-Z\u0E00-\u0E7F\s]/g, '').replace(/\s{2,}/g, ' ');
                if (val !== cleanVal) nameInput.value = cleanVal;
                updateBtnState();
            });
            clearBtn.addEventListener('click', () => {
                nameInput.value = '';
                updateBtnState();
                nameInput.focus();
                clearTimeout(typingTimer);
                typingStatus.style.display = 'none';
            });
            updateBtnState();
        }

        if(birthDayInput) {
            birthDayInput.addEventListener('change', () => {
                if (nameInput.value.trim()) { performAnalysis(); }
            });
        }

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

    // --- 4. Result Features ---
    // ฟังก์ชันนี้รวม Logic ของปุ่มและกราฟิกทั้งหมด
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

        // 4.2 Highlight Text
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

        // 4.3 Detail Toggle
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

        // 4.4 AI Modal
        const btnLang = document.getElementById('btnLinguistics');
        const modal = document.getElementById('aiModal');
        const closeModal = document.querySelector('.close-modal');
        const aiContent = document.getElementById('aiContent');

        if (btnLang && dataContainer) {
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
        if (closeModal) closeModal.onclick = () => modal.style.display = "none";
        window.onclick = (event) => { if (event.target == modal) modal.style.display = "none"; }

        // 🔥 4.5 SAVE BUTTON LOGIC (สำคัญมาก ต้องอยู่ตรงนี้!) 🔥
        // ใช้ getElementById ตามที่เราแก้ใน HTML
        const btnSave = document.getElementById('btnSaveName');

        if (btnSave && dataContainer) {
            // เคลียร์ event เก่าก่อน (กันเบิ้ล)
            const newBtnSave = btnSave.cloneNode(true);
            btnSave.parentNode.replaceChild(newBtnSave, btnSave);

            newBtnSave.addEventListener('click', async () => {
                const name = dataContainer.getAttribute('data-name');
                const birthDay = document.getElementById('birthDayInput').value;

                newBtnSave.innerText = "⏳ กำลังบันทึก...";
                newBtnSave.disabled = true;

                try {
                    const response = await fetch('/api/save-name', {
                        method: 'POST',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify({ name: name, birth_day: birthDay })
                    });

                    const resData = await response.json();

                    if (response.status === 401 || resData.redirect) {
                        // กรณีไม่ได้ Login
                        if(confirm("คุณต้องเข้าสู่ระบบก่อนบันทึกชื่อ\nกด 'ตกลง' เพื่อไปหน้าเข้าสู่ระบบ")) {
                            window.location.href = '/login';
                        } else {
                            newBtnSave.innerText = "💾 บันทึกชื่อ";
                            newBtnSave.disabled = false;
                        }
                    } else if (response.ok) {
                        // กรณีสำเร็จ
                        newBtnSave.innerText = "✅ บันทึกแล้ว";
                        newBtnSave.style.color = "green";
                    } else {
                        // กรณี Error อื่นๆ
                        alert("เกิดข้อผิดพลาด: " + resData.error);
                        newBtnSave.innerText = "💾 บันทึกชื่อ";
                        newBtnSave.disabled = false;
                    }
                } catch (err) {
                    console.error(err);
                    alert("เชื่อมต่อเซิร์ฟเวอร์ไม่ได้");
                    newBtnSave.innerText = "💾 บันทึกชื่อ";
                    newBtnSave.disabled = false;
                }
            });
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
});