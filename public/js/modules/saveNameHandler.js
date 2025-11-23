// public/js/modules/saveNameHandler.js

export function initializeSaveNameHandler(btnSaveName, nameInput, birthDayInput) {
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
}
