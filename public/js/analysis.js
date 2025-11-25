// public/js/analysis.js
import { initializeInputHandler } from './modules/inputHandler.js';
import { initializeNameInteraction } from './modules/nameInteraction.js';
import { initializeScoreAnimator } from './modules/scoreAnimator.js';
import { initializeDetailToggler } from './modules/detailToggler.js';
import { initializeAiModalHandler } from './modules/aiModalHandler.js';
import { initializeSaveNameHandler } from './modules/saveNameHandler.js';
import { initializeSolarSystemAnimator } from './modules/solarSystemAnimator.js';

function reinitializeScripts() {
    // Re-find elements that might have been replaced
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
    const nameInput = document.getElementById('nameInput');
    const birthDayInput = document.getElementById('birthDayInput');

    // Only re-initialize if the main result container exists
    if (document.querySelector('.summary-box')) {
        initializeScoreAnimator(analysisData, sunNameDisplay);
        initializeDetailToggler(btnShowDetails, detailSection, btnCloseDetails);
        initializeAiModalHandler(btnLinguistics, aiModal, closeModal, aiContent, nameInput);
        initializeSaveNameHandler(btnSaveName, nameInput, birthDayInput);
        initializeSolarSystemAnimator();
        
        // This needs to be re-initialized to handle clicks on new similar names
        initializeNameInteraction(birthDayInput); 
    }
}

document.addEventListener('DOMContentLoaded', function () {
    const nameInput = document.getElementById('nameInput');
    const birthDayInput = document.getElementById('birthDayInput');
    const clearInputBtn = document.getElementById('clearInputBtn');
    const typingStatus = document.getElementById('typingStatus');

    // Initialize Input Handler (this should only be initialized once)
    initializeInputHandler(nameInput, birthDayInput, clearInputBtn, typingStatus);

    // Initial load of scripts
    reinitializeScripts();

    // Listen for the custom event to re-initialize scripts on new content
    document.addEventListener('analysisUpdated', reinitializeScripts);
});
