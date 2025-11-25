// public/js/modules/inputHandler.js

export function initializeInputHandler(nameInput, birthDayInput, clearInputBtn, typingStatus) {
    let typingTimer;
    const doneTypingInterval = 2000; // 2 seconds

    const fetchAnalysis = () => {
        const form = nameInput.closest('form');
        if (!form) {
            console.error('Form not found');
            return;
        }

        const formData = new FormData(form);
        const params = new URLSearchParams(formData);

        if (typingStatus) {
            typingStatus.textContent = '⏳ กำลังคำนวณ...';
            typingStatus.style.display = 'block';
        }

        fetch('/analysis', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
            },
            body: params.toString(),
        })
        .then(response => response.text())
        .then(html => {
            const parser = new DOMParser();
            const doc = parser.parseFromString(html, 'text/html');
            const newResultWrapper = doc.getElementById('result-wrapper');
            const currentResultWrapper = document.getElementById('result-wrapper');

            if (newResultWrapper && currentResultWrapper) {
                currentResultWrapper.innerHTML = newResultWrapper.innerHTML;
            }
            
            // Dispatch event to re-initialize scripts on the new content
            document.dispatchEvent(new CustomEvent('analysisUpdated'));

            if (typingStatus) {
                typingStatus.style.display = 'none';
                typingStatus.textContent = '⏳ กำลังรอพิมพ์เสร็จ...';
            }
        })
        .catch(error => {
            console.error('Error fetching analysis:', error);
            if (typingStatus) {
                typingStatus.textContent = 'เกิดข้อผิดพลาด!';
            }
        });
    };

    const doneTyping = () => {
        if (typingStatus) typingStatus.style.display = 'none';
        fetchAnalysis();
    };

    if (nameInput) {
        nameInput.addEventListener('keyup', (e) => {
            // Ignore keys that don't change the value or are used for submission
            if (['Shift', 'Control', 'Alt', 'Meta', 'ArrowUp', 'ArrowDown', 'ArrowLeft', 'ArrowRight', 'Enter'].includes(e.key)) {
                return;
            }
            clearTimeout(typingTimer);
            if (typingStatus) typingStatus.style.display = 'block';
            if (clearInputBtn) clearInputBtn.style.display = nameInput.value ? 'block' : 'none';
            typingTimer = setTimeout(doneTyping, doneTypingInterval);
        });

        nameInput.addEventListener('keydown', (e) => {
            // On Enter key, prevent default form submission and trigger analysis immediately
            if (e.key === 'Enter') {
                e.preventDefault(); 
                clearTimeout(typingTimer); // Stop the typing timer
                fetchAnalysis();
            } else {
                // For other keys, just clear the timer
                clearTimeout(typingTimer);
            }
        });

        if (nameInput.value && clearInputBtn) {
            clearInputBtn.style.display = 'block';
        }
    }

    if (clearInputBtn) {
        clearInputBtn.addEventListener('click', () => {
            if (nameInput) nameInput.value = '';
            if (clearInputBtn) clearInputBtn.style.display = 'none';
            if (nameInput) nameInput.focus();
        });
    }

    if (birthDayInput) {
        birthDayInput.addEventListener('change', () => {
            fetchAnalysis();
        });
    }
}
