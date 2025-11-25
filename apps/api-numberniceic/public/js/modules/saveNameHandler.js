// public/js/modules/saveNameHandler.js

export function initializeSaveNameHandler(btnSaveName, nameInput, birthDayInput) {
    const notificationModal = document.getElementById('notificationModal');
    const notificationIcon = document.getElementById('notificationIcon');
    const notificationMessage = document.getElementById('notificationMessage');
    const notificationCloseBtn = document.getElementById('notificationCloseBtn');

    const showNotification = (isSuccess, message) => {
        notificationIcon.innerHTML = isSuccess ? '<i class="fa-solid fa-circle-check"></i>' : '<i class="fa-solid fa-circle-xmark"></i>';
        notificationIcon.className = ''; // Clear previous classes
        notificationIcon.classList.add(isSuccess ? 'success' : 'error');
        notificationMessage.textContent = message;
        notificationModal.classList.add('show');
    };

    const hideNotification = () => {
        notificationModal.classList.remove('show');
    };

    if (notificationCloseBtn) {
        notificationCloseBtn.addEventListener('click', hideNotification);
        notificationModal.addEventListener('click', (e) => {
            if (e.target === notificationModal) {
                hideNotification();
            }
        });
    }

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
                    showNotification(true, 'บันทึกชื่อสำเร็จแล้ว!');
                } else if (data.redirect) {
                    window.location.href = data.redirect;
                } else {
                    showNotification(false, 'เกิดข้อผิดพลาด: ' + data.error);
                }
            })
            .catch(() => showNotification(false, 'เกิดข้อผิดพลาดในการเชื่อมต่อ'));
        });
    }
}
