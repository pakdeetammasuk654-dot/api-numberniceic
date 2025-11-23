import { initializeClickableRows } from './modules/clickableRowHandler.js';
import { initializeDeleteConfirmation } from './modules/deleteConfirmationHandler.js';

document.addEventListener('DOMContentLoaded', function() {
    initializeClickableRows();
    initializeDeleteConfirmation();
});
