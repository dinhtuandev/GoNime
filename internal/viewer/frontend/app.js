let images = [];
let currentIndex = 0;
let zoomLevel = 100;
let viewerToken = '';

let chapters = [];
let currentChapterIndex = 0;
let currentChapterTitle = '';

const loading = document.getElementById('loading');
const pageInfo = document.getElementById('page-info');
const pageInput = document.getElementById('page-input');
const zoomLevelDisplay = document.getElementById('zoom-level');
const prevBtn = document.getElementById('prev-btn');
const nextBtn = document.getElementById('next-btn');
const chapterInfo = document.getElementById('chapter-info');
const prevChapterBtn = document.getElementById('prev-chapter-btn');
const nextChapterBtn = document.getElementById('next-chapter-btn');
const container = document.getElementById('viewer-container');
const column = document.getElementById('manga-column');

let syncInterval = null;
let lastSignature = '';

function buildSignature(data) {
    const imgs = data.images || [];
    return [
        imgs.length,
        data.start || 0,
        data.chapterIndex || 0,
        imgs[0] || '',
        imgs[imgs.length - 1] || ''
    ].join('|');
}

async function fetchState() {
    const response = await fetch('/api/images', { cache: 'no-store' });
    const data = await response.json();
    if (data.token) {
        viewerToken = data.token;
    }
    return data;
}

function updateFromState(data, forceRender = false) {
    const signature = buildSignature(data);
    if (!forceRender && signature === lastSignature) {
        return;
    }

    lastSignature = signature;
    images = data.images || [];
    currentIndex = Math.max(0, Math.min(data.start || 0, images.length - 1));
    chapters = data.chapters || [];
    currentChapterIndex = data.chapterIndex || 0;
    currentChapterTitle = data.chapterTitle || '';

    renderColumn();
    updateChapterInfo();
    updatePageInfo();
    scrollToIndex(currentIndex, false);
}

async function initViewer(imageList, startIndex = 0) {
    if (Array.isArray(imageList) && imageList.length > 0) {
        updateFromState({ images: imageList, start: startIndex, chapters: [], chapterIndex: 0, chapterTitle: '' }, true);
    }

    try {
        const data = await fetchState();
        updateFromState(data, true);
    } catch (e) {
        console.error('Failed to load images:', e);
        loading.textContent = 'Khong the tai du lieu';
    }

    if (!syncInterval) {
        syncInterval = setInterval(async () => {
            try {
                const data = await fetchState();
                updateFromState(data);
            } catch (_) {
                // no-op
            }
        }, 800);
    }
}

function renderColumn() {
    column.innerHTML = '';

    if (!images || images.length === 0) {
        loading.style.display = 'block';
        loading.textContent = 'Khong co hinh anh';
        return;
    }

    loading.style.display = 'none';

    images.forEach((src, idx) => {
        const wrap = document.createElement('div');
        wrap.className = 'manga-page-wrap';

        const img = document.createElement('img');
        img.className = 'manga-page';
        img.src = src;
        img.alt = `Manga page ${idx + 1}`;
        img.loading = 'lazy';
        img.dataset.index = String(idx);

        if (idx === currentIndex) {
            img.classList.add('active');
        }

        img.addEventListener('click', () => {
            currentIndex = idx;
            highlightActive();
            updatePageInfo();
        });

        wrap.appendChild(img);
        column.appendChild(wrap);
    });

    applyZoom();
    preloadAroundCurrent();
}

function highlightActive() {
    document.querySelectorAll('.manga-page.active').forEach(el => el.classList.remove('active'));
    const active = column.querySelector(`.manga-page[data-index="${currentIndex}"]`);
    if (active) {
        active.classList.add('active');
    }
    preloadAroundCurrent();
}

function preloadAroundCurrent() {
    if (!images || images.length === 0) return;
    const preloadSpan = 3;
    for (let i = 1; i <= preloadSpan; i++) {
        const next = currentIndex + i;
        const prev = currentIndex - i;
        if (next < images.length) {
            const img = new Image();
            img.src = images[next];
        }
        if (prev >= 0) {
            const img = new Image();
            img.src = images[prev];
        }
    }
}

function updateCurrentIndexByViewport() {
    const pages = column.querySelectorAll('.manga-page');
    if (!pages.length) return;

    const viewportTop = container.scrollTop;
    const viewportCenter = viewportTop + container.clientHeight * 0.35;

    let nearestIdx = currentIndex;
    let nearestDistance = Number.POSITIVE_INFINITY;

    pages.forEach((img) => {
        const rectTop = img.offsetTop;
        const distance = Math.abs(rectTop - viewportCenter);
        if (distance < nearestDistance) {
            nearestDistance = distance;
            nearestIdx = Number(img.dataset.index || '0');
        }
    });

    if (nearestIdx !== currentIndex) {
        currentIndex = nearestIdx;
        highlightActive();
        updatePageInfo();
    }
}

function scrollToIndex(index, smooth = true) {
    if (index < 0 || index >= images.length) return;
    const target = column.querySelector(`.manga-page[data-index="${index}"]`);
    if (!target) return;
    target.scrollIntoView({ block: 'start', behavior: smooth ? 'smooth' : 'auto' });
}

function updatePageInfo() {
    if (images.length === 0) {
        pageInfo.textContent = '0 / 0';
        pageInput.value = 1;
        prevBtn.disabled = true;
        nextBtn.disabled = true;
        return;
    }

    pageInfo.textContent = `${currentIndex + 1} / ${images.length}`;
    pageInput.value = currentIndex + 1;
    prevBtn.disabled = currentIndex === 0;
    nextBtn.disabled = currentIndex >= images.length - 1;
}

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

function prevPage() {
    if (currentIndex > 0) {
        currentIndex -= 1;
        highlightActive();
        updatePageInfo();
        scrollToIndex(currentIndex);
    }
}

function nextPage() {
    if (currentIndex < images.length - 1) {
        currentIndex += 1;
        highlightActive();
        updatePageInfo();
        scrollToIndex(currentIndex);
    }
}

function firstPage() {
    currentIndex = 0;
    highlightActive();
    updatePageInfo();
    scrollToIndex(0);
}

function lastPage() {
    currentIndex = Math.max(images.length - 1, 0);
    highlightActive();
    updatePageInfo();
    scrollToIndex(currentIndex);
}

function goToPage(pageNum) {
    const idx = Number.parseInt(pageNum, 10) - 1;
    if (Number.isNaN(idx) || idx < 0 || idx >= images.length) {
        pageInput.value = currentIndex + 1;
        return;
    }
    currentIndex = idx;
    highlightActive();
    updatePageInfo();
    scrollToIndex(currentIndex);
}

function prevChapter() {
    if (currentChapterIndex > 0) {
        notifyChapterChange(currentChapterIndex - 1);
    }
}

function nextChapter() {
    if (currentChapterIndex < chapters.length - 1) {
        notifyChapterChange(currentChapterIndex + 1);
    }
}

async function notifyChapterChange(newIndex) {
    try {
        const response = await fetch('/api/chapter', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'X-Viewer-Token': viewerToken
            },
            body: JSON.stringify({ index: newIndex })
        });
        const data = await response.json();
        if (data.images) {
            updateFromState({
                images: data.images,
                start: 0,
                chapters,
                chapterIndex: newIndex,
                chapterTitle: data.chapterTitle || ''
            }, true);
        }
    } catch (e) {
        console.error('Failed to load chapter:', e);
    }
}

function zoomIn() {
    if (zoomLevel < 180) {
        zoomLevel += 10;
        applyZoom();
    }
}

function zoomOut() {
    if (zoomLevel > 60) {
        zoomLevel -= 10;
        applyZoom();
    }
}

function fitToWindow() {
    zoomLevel = 100;
    applyZoom();
}

function applyZoom() {
    const scale = zoomLevel / 100;
    column.style.setProperty('--zoom', String(scale));
    zoomLevelDisplay.textContent = `${zoomLevel}%`;
}

document.addEventListener('keydown', (e) => {
    switch (e.key) {
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

container.addEventListener('scroll', () => {
    updateCurrentIndexByViewport();
});

container.addEventListener('wheel', (e) => {
    if (!e.ctrlKey) return;
    e.preventDefault();
    if (e.deltaY < 0) {
        zoomIn();
    } else {
        zoomOut();
    }
}, { passive: false });

window.addEventListener('load', () => {
    initViewer();
});

window.addEventListener('beforeunload', () => {
    if (syncInterval) {
        clearInterval(syncInterval);
        syncInterval = null;
    }
});

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
