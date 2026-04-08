const STORAGE_KEY_CONFIG = 'txt2ig_last_config';
const STORAGE_KEY_TEXT = 'txt2ig_last_text';
const DEBOUNCE_DELAY = 1000;

function debounce(func, wait) {
    let timeout;
    return function executedFunction(...args) {
        const later = () => {
            clearTimeout(timeout);
            func(...args);
        };
        clearTimeout(timeout);
        timeout = setTimeout(later, wait);
    };
}

function saveToStorage() {
    const textArea = document.getElementById('text');
    const configArea = document.getElementById('config');
    
    if (textArea && textArea.value) {
        localStorage.setItem(STORAGE_KEY_TEXT, textArea.value);
    }
    
    if (configArea && configArea.value) {
        localStorage.setItem(STORAGE_KEY_CONFIG, configArea.value);
    }
}

function loadFromStorage() {
    const textArea = document.getElementById('text');
    const configArea = document.getElementById('config');
    
    if (textArea) {
        const savedText = localStorage.getItem(STORAGE_KEY_TEXT);
        if (savedText && !textArea.value) {
            textArea.value = savedText;
        }
    }
    
    if (configArea) {
        const savedConfig = localStorage.getItem(STORAGE_KEY_CONFIG);
        if (savedConfig && !configArea.value) {
            configArea.value = savedConfig;
        }
    }
}

const debouncedSave = debounce(saveToStorage, DEBOUNCE_DELAY);

function setupAutoSave() {
    const textArea = document.getElementById('text');
    const configArea = document.getElementById('config');
    
    if (textArea) {
        textArea.addEventListener('input', debouncedSave);
    }
    
    if (configArea) {
        configArea.addEventListener('input', debouncedSave);
    }
    
    loadFromStorage();
}

if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', setupAutoSave);
} else {
    setupAutoSave();
}

// Download image from hidden input data
function downloadImageFromData() {
    const base64 = document.getElementById('image-base64').value;
    const format = document.getElementById('image-format').value;
    
    const link = document.createElement('a');
    link.href = 'data:image/' + format + ';base64,' + base64;
    link.download = 'txt2ig-image.' + format;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
}