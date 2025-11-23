// public/js/modules/inputHandler.js

export function initializeInputHandler(nameInput, birthDayInput, clearInputBtn, typingStatus) {
    let typingTimer;
    const doneTypingInterval = 1000; // 1 second

    const submitForm = () => {
        // Assuming the form is the parent of one of the inputs or can be found easily
        // For a more robust solution, pass the form element or a submit callback
        const form = nameInput ? nameInput.closest('form') : birthDayInput ? birthDayInput.closest('form') : null;
        if (form) {
            form.submit();
        } else {
            console.error('Form not found for submission.');
        }
    };

    const doneTyping = () => {
        if (typingStatus) typingStatus.style.display = 'none';
        submitForm();
    };

    if (nameInput) {
        nameInput.addEventListener('keyup', () => {
            clearTimeout(typingTimer);
            if (typingStatus) typingStatus.style.display = 'block';
            if (clearInputBtn) clearInputBtn.style.display = nameInput.value ? 'block' : 'none';
            typingTimer = setTimeout(doneTyping, doneTypingInterval);
        });

        nameInput.addEventListener('keydown', () => {
            clearTimeout(typingTimer);
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
            submitForm();
        });
    }
}
