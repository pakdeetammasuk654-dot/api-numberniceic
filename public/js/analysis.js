// public/js/analysis.js
import { initializeInputHandler } from './modules/inputHandler.js';
import { initializeNameInteraction } from './modules/nameInteraction.js';
import { initializeScoreAnimator } from './modules/scoreAnimator.js';
import { initializeDetailToggler } from './modules/detailToggler.js';
import { initializeAiModalHandler } from './modules/aiModalHandler.js';
import { initializeSaveNameHandler } from './modules/saveNameHandler.js';

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

    // Initialize Input Handler
    initializeInputHandler(nameInput, birthDayInput, clearInputBtn, typingStatus);

    // Initialize Name Interaction (Sample Names, Similar Names)
    initializeNameInteraction(birthDayInput);

    // Initialize Score Animator
    initializeScoreAnimator(analysisData, sunNameDisplay);

    // Initialize Detail Section Toggle
    initializeDetailToggler(btnShowDetails, detailSection, btnCloseDetails);

    // Initialize AI Linguistics Modal
    initializeAiModalHandler(btnLinguistics, aiModal, closeModal, aiContent, nameInput);

    // Initialize Save Name
    initializeSaveNameHandler(btnSaveName, nameInput, birthDayInput);
});
