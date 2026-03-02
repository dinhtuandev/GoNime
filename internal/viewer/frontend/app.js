// Global state
let images = [];
let currentIndex = 0;
let zoomLevel = 100;
let apiBase = '';

// Chapter navigation state
let chapters = [];
let currentChapterIndex = 0;
let currentChapterTitle = '';

// DOM elements
const mangaImage = document.getElementById('manga-image');
const loading = document.getElementById('loading');
const pageInfo = document.getElementById('page-info');
const pageInput = document.getElementById('page-input');
const zoomLevelDisplay = document.getElementById('zoom-level');
const prevBtn = document.getElementById('prev-btn');
const nextBtn = document.getElementById('next-btn');
const chapterInfo = document.getElementById('chapter-info');
const prevChapterBtn = document.getElementById('prev-chapter-btn');
const nextChapterBtn = document.getElementById('next-chapter-btn');

// Initialize viewer - called from Wails or manually
async function initViewer(imageList, startIndex = 0) {
    images = imageList || [];
    currentIndex = startIndex || 0;
    
    if (images.length === 0) {
        // Try to load from API
        try {
            const response = await fetch('/api/images');
            const data = await response.json();
            images = data.images || [];
            currentIndex = data.start || 0;
            chapters = data.chapters || [];
            currentChapterIndex = data.chapterIndex || 0;
            currentChapterTitle = data.chapterTitle || '';
        } catch (e) {
            console.error('Failed to load images:', e);
        }
    }
    
    if (images.length === 0) {
        loading.textContent = 'Không có hình ảnh';
        return;
    }
    
    updatePageInfo();
    updateChapterInfo();
    loadImage(currentIndex);
}

// Load image at specific index
function loadImage(index) {
    if (index < 0 || index >= images.length) return;
    
    currentIndex = index;
    loading.style.display = 'block';
    mangaImage.style.display = 'none';
    
    mangaImage.src = images[currentIndex];
    updatePageInfo();
}

// Update page info display
function updatePageInfo() {
    pageInfo.textContent = `${currentIndex + 1} / ${images.length}`;
    pageInput.value = currentIndex + 1;
    
    // Update button states
    prevBtn.disabled = currentIndex === 0;
    nextBtn.disabled = currentIndex >= images.length - 1;
}

// Update chapter info display
function updateChapterInfo() {
    if (chapters.length > 0 && currentChapterIndex >= 0) {
        chapterInfo.textContent = chapters[currentChapterIndex]?.title || `Chapter ${currentChapterIndex + 1}`;
        prevChapterBtn.disabled = currentChapterIndex <= 0;
        nextChapterBtn.disabled = currentChapterIndex >= chapters.length - 1;
    } else {
        chapterInfo.textContent = currentChapterTitle || 'Chapter';
        prevChapterBtn.disabled = true;
        nextChapterBtn.disabled = true;
    }
}

// Navigation functions
function prevPage() {
    if (currentIndex > 0) {
        loadImage(currentIndex - 1);
    }
}

function nextPage() {
    if (currentIndex < images.length - 1) {
        loadImage(currentIndex + 1);
    }
}

function firstPage() {
    loadImage(0);
}

function lastPage() {
    loadImage(images.length - 1);
}

function goToPage(pageNum) {
    const index = parseInt(pageNum) - 1;
    if (index >= 0 && index < images.length) {
        loadImage(index);
    } else {
        pageInput.value = currentIndex + 1;
    }
}

// Chapter navigation functions
function prevChapter() {
    if (currentChapterIndex > 0) {
        currentChapterIndex--;
        // Notify backend to load previous chapter
        notifyChapterChange(currentChapterIndex);
    }
}

function nextChapter() {
    if (currentChapterIndex < chapters.length - 1) {
        currentChapterIndex++;
        // Notify backend to load next chapter
        notifyChapterChange(currentChapterIndex);
    }
}

// Notify backend about chapter change
async function notifyChapterChange(newIndex) {
    try {
        const response = await fetch('/api/chapter', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ index: newIndex })
        });
        const data = await response.json();
        
        if (data.images) {
            images = data.images;
            currentIndex = 0;
            currentChapterTitle = data.chapterTitle || '';
            updatePageInfo();
            updateChapterInfo();
            loadImage(0);
        }
    } catch (e) {
        console.error('Failed to load chapter:', e);
    }
}

// Zoom functions
function zoomIn() {
    if (zoomLevel < 300) {
        zoomLevel += 25;
        applyZoom();
    }
}

function zoomOut() {
    if (zoomLevel > 25) {
        zoomLevel -= 25;
        applyZoom();
    }
}

function fitToWindow() {
    zoomLevel = 100;
    applyZoom();
}

function applyZoom() {
    mangaImage.style.transform = `scale(${zoomLevel / 100})`;
    zoomLevelDisplay.textContent = `${zoomLevel}%`;
}

// Event callbacks
function imageLoaded() {
    loading.style.display = 'none';
    mangaImage.style.display = 'block';
}

function imageError() {
    loading.textContent = 'Lỗi tải hình ảnh';
}

// Keyboard navigation
document.addEventListener('keydown', (e) => {
    switch(e.key) {
        case 'ArrowLeft':
        case 'a':
            prevPage();
            break;
        case 'ArrowRight':
        case 'd':
            nextPage();
            break;
        case 'ArrowUp':
        case 'w':
            prevChapter();
            break;
        case 'ArrowDown':
        case 's':
            nextChapter();
            break;
        case 'Home':
            firstPage();
            break;
        case 'End':
            lastPage();
            break;
        case '+':
        case '=':
            zoomIn();
            break;
        case '-':
            zoomOut();
            break;
    }
});

// Mouse wheel zoom
document.getElementById('viewer-container').addEventListener('wheel', (e) => {
    if (e.ctrlKey) {
        e.preventDefault();
        if (e.deltaY < 0) {
            zoomIn();
        } else {
            zoomOut();
        }
    }
});

// Drag to pan
let isDragging = false;
let startX, startY, scrollLeft, scrollTop;

const container = document.getElementById('viewer-container');

container.addEventListener('mousedown', (e) => {
    isDragging = true;
    startX = e.pageX - container.offsetLeft;
    startY = e.pageY - container.offsetTop;
    scrollLeft = container.scrollLeft;
    scrollTop = container.scrollTop;
    container.style.cursor = 'grabbing';
});

container.addEventListener('mouseleave', () => {
    isDragging = false;
    container.style.cursor = 'grab';
});

container.addEventListener('mouseup', () => {
    isDragging = false;
    container.style.cursor = 'grab';
});

container.addEventListener('mousemove', (e) => {
    if (!isDragging) return;
    e.preventDefault();
    const x = e.pageX - container.offsetLeft;
    const y = e.pageY - container.offsetTop;
    const walkX = (x - startX) * 1.5;
    const walkY = (y - startY) * 1.5;
    container.scrollLeft = scrollLeft - walkX;
    container.scrollTop = scrollTop - walkY;
});

// Auto-initialize when page loads
window.addEventListener('load', () => {
    initViewer();
});

// Expose functions globally for Wails
window.initViewer = initViewer;
window.prevPage = prevPage;
window.nextPage = nextPage;
window.firstPage = firstPage;
window.lastPage = lastPage;
window.goToPage = goToPage;
window.prevChapter = prevChapter;
window.nextChapter = nextChapter;
window.zoomIn = zoomIn;
window.zoomOut = zoomOut;
window.fitToWindow = fitToWindow;