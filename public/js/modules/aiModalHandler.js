// public/js/modules/aiModalHandler.js

export function initializeAiModalHandler(btnLinguistics, aiModal, closeModal, aiContent, nameInput) {
    if (btnLinguistics && aiModal && closeModal && aiContent && nameInput) {
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
}
