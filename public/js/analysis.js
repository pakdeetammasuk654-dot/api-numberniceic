document.addEventListener('DOMContentLoaded', function () {
    const nameInput = document.getElementById('nameInput');
    const birthDayInput = document.getElementById('birthDayInput');
    const clearInputBtn = document.getElementById('clearInputBtn');
    const typingStatus = document.getElementById('typingStatus');
    const analysisData = document.getElementById('analysisData');
    const sunNameDisplay = document.getElementById('sunNameDisplay');
    const btnShowDetails = document.getElementById('btnShowDetails');
    const btnLinguistics = document.getElementById('btnLinguistics');
    const btnCloseDetails = document.getElementById('btnCloseDetails');
    const detailSection = document.getElementById('detailSection');
    const aiModal = document.getElementById('aiModal');
    const closeModal = document.querySelector('.close-modal');
    const aiContent = document.getElementById('aiContent');
    const btnSaveName = document.getElementById('btnSaveName');

    let typingTimer;
    const doneTypingInterval = 1000; // 1 second

    // --- Form and Input Handling ---
    if (nameInput) {
        nameInput.addEventListener('keyup', () => {
            clearTimeout(typingTimer);
            typingStatus.style.display = 'block';
            clearInputBtn.style.display = nameInput.value ? 'block' : 'none';
            typingTimer = setTimeout(doneTyping, doneTypingInterval);
        });

        nameInput.addEventListener('keydown', () => {
            clearTimeout(typingTimer);
        });

        if (nameInput.value) {
            clearInputBtn.style.display = 'block';
        }
    }

    if (clearInputBtn) {
        clearInputBtn.addEventListener('click', () => {
            nameInput.value = '';
            clearInputBtn.style.display = 'none';
            nameInput.focus();
        });
    }

    if (birthDayInput) {
        birthDayInput.addEventListener('change', () => {
            document.querySelector('form').submit();
        });
    }

    function doneTyping() {
        typingStatus.style.display = 'none';
        document.querySelector('form').submit();
    }

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

    window.analyzeName = function(name) {
        const birthDay = birthDayInput.value;
        window.location.href = `/analysis?name=${encodeURIComponent(name)}&birth_day=${birthDay}`;
    }

    // --- Animation and Display ---
    if (analysisData) {
        const name = analysisData.dataset.name;
        const kakis = analysisData.dataset.kakis.split(',');

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

    // --- Detail Section Toggle ---
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

    // --- AI Linguistics Modal ---
    if (btnLinguistics && aiModal && closeModal && aiContent) {
        btnLinguistics.addEventListener('click', () => {
            const name = nameInput.value;
            aiContent.innerHTML = '<p>กำลังวิเคราะห์ด้วย AI... 🤖</p>';
            aiModal.style.display = 'block';
            fetch(`/api/linguistics?name=${encodeURIComponent(name)}`)
                .then(response => response.json())
                .then(data => {
                    if (data.text) {
                        aiContent.innerHTML = data.text.replace(/\n/g, '<br>');
                    } else {
                        aiContent.textContent = 'ขออภัย, ไม่สามารถวิเคราะห์ได้ในขณะนี้';
                    }
                })
                .catch(() => {
                    aiContent.textContent = 'เกิดข้อผิดพลาดในการเชื่อมต่อ';
                });
        });

        closeModal.addEventListener('click', () => {
            aiModal.style.display = 'none';
        });

        window.addEventListener('click', (event) => {
            if (event.target === aiModal) {
                aiModal.style.display = 'none';
            }
        });
    }

    // --- Save Name ---
    if (btnSaveName) {
        btnSaveName.addEventListener('click', () => {
            const name = nameInput.value;
            const birthDay = birthDayInput.value;
            fetch('/api/save-name', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ name: name, birth_day: birthDay })
            })
            .then(response => response.json())
            .then(data => {
                if (data.success) {
                    alert('บันทึกชื่อสำเร็จ!');
                } else if (data.redirect) {
                    window.location.href = data.redirect;
                } else {
                    alert('เกิดข้อผิดพลาด: ' + data.error);
                }
            })
            .catch(() => alert('เกิดข้อผิดพลาดในการเชื่อมต่อ'));
        });
    }
});
