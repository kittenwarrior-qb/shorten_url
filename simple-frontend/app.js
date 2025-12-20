// API Configuration
const API_BASE_URL = 'https://shorten.quocbui.dev/api/v1';

// Get token from localStorage
function getToken() {
    return localStorage.getItem('guest_token');
}

// Save token to localStorage
function saveToken(token) {
    if (token) {
        localStorage.setItem('guest_token', token);
    }
}

// Show/hide elements
function showElement(id) {
    document.getElementById(id).style.display = 'block';
}

function hideElement(id) {
    document.getElementById(id).style.display = 'none';
}

// Show error message
function showError(message) {
    const errorEl = document.getElementById('error');
    errorEl.textContent = message;
    showElement('error');
    hideElement('result');
}

// Show success result
function showResult(data) {
    hideElement('error');
    
    const shortUrlInput = document.getElementById('shortUrl');
    shortUrlInput.value = data.link.short_url;
    
    // Show QR code if available
    const qrCodeEl = document.getElementById('qrCode');
    if (data.link.qr_code) {
        qrCodeEl.innerHTML = `<img src="${data.link.qr_code}" alt="QR Code">`;
    }
    
    showElement('result');
    
    // Save token if provided (guest user)
    if (data.token) {
        saveToken(data.token);
    }
    
    // Reload links list
    loadMyLinks();
}

// Copy to clipboard
function copyToClipboard() {
    const shortUrlInput = document.getElementById('shortUrl');
    shortUrlInput.select();
    document.execCommand('copy');
    
    const btn = event.target;
    const originalText = btn.textContent;
    btn.textContent = 'Copied!';
    setTimeout(() => {
        btn.textContent = originalText;
    }, 2000);
}

// Shorten URL
async function shortenURL(url, alias) {
    const token = getToken();
    
    const headers = {
        'Content-Type': 'application/json'
    };
    
    if (token) {
        headers['Authorization'] = `Bearer ${token}`;
    }
    
    const body = { url };
    if (alias) {
        body.alias = alias;
    }
    
    const response = await fetch(`${API_BASE_URL}/shorten`, {
        method: 'POST',
        headers: headers,
        body: JSON.stringify(body)
    });
    
    if (!response.ok) {
        const error = await response.json();
        throw new Error(error.message || 'Failed to shorten URL');
    }
    
    return await response.json();
}

// Load my links
async function loadMyLinks() {
    const token = getToken();
    if (!token) {
        return;
    }
    
    try {
        const response = await fetch(`${API_BASE_URL}/me/links?page=1&per_page=10`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        
        if (!response.ok) {
            throw new Error('Failed to load links');
        }
        
        const data = await response.json();
        displayLinks(data.data.links || []);
    } catch (error) {
        console.error('Error loading links:', error);
    }
}

// Display links
function displayLinks(links) {
    const linksListEl = document.getElementById('linksList');
    
    if (!links || links.length === 0) {
        linksListEl.innerHTML = '<p class="empty-state">No links yet</p>';
        return;
    }
    
    linksListEl.innerHTML = links.map(link => `
        <div class="link-item" onclick="showLinkDetail('${link.short_code}')">
            <div class="link-item-header">
                <a href="${link.short_url}" target="_blank" class="link-item-url" onclick="event.stopPropagation()">
                    ${link.short_url}
                </a>
                <div class="link-item-stats">
                    <span>${link.click_count} clicks</span>
                    <span>${new Date(link.created_at).toLocaleDateString()}</span>
                </div>
            </div>
            <div class="link-item-original">${link.original_url}</div>
        </div>
    `).join('');
}

// Show link detail with QR code
async function showLinkDetail(shortCode) {
    const token = getToken();
    if (!token) return;
    
    try {
        const response = await fetch(`${API_BASE_URL}/me/links/${shortCode}`, {
            headers: { 'Authorization': `Bearer ${token}` }
        });
        
        if (!response.ok) throw new Error('Failed to load link');
        
        const data = await response.json();
        const link = data.data.link;
        
        document.getElementById('shortUrl').value = link.short_url;
        
        const qrCodeEl = document.getElementById('qrCode');
        if (link.qr_code) {
            qrCodeEl.innerHTML = `<img src="${link.qr_code}" alt="QR Code">`;
        }
        
        showElement('result');
        hideElement('error');
        
        // Scroll to result
        document.getElementById('result').scrollIntoView({ behavior: 'smooth' });
    } catch (error) {
        showError('Failed to load link details');
    }
}

// Form submit handler
document.getElementById('shortenForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    
    const url = document.getElementById('url').value.trim();
    const alias = document.getElementById('alias').value.trim();
    
    const submitBtn = e.target.querySelector('button[type="submit"]');
    const btnText = submitBtn.querySelector('.btn-text');
    const btnLoading = submitBtn.querySelector('.btn-loading');
    
    // Show loading state
    submitBtn.disabled = true;
    btnText.style.display = 'none';
    btnLoading.style.display = 'inline';
    hideElement('error');
    hideElement('result');
    
    try {
        const data = await shortenURL(url, alias);
        showResult(data.data);
        
        // Reset form
        document.getElementById('url').value = '';
        document.getElementById('alias').value = '';
    } catch (error) {
        showError(error.message);
    } finally {
        // Reset button state
        submitBtn.disabled = false;
        btnText.style.display = 'inline';
        btnLoading.style.display = 'none';
    }
});

// Load links on page load
window.addEventListener('DOMContentLoaded', () => {
    loadMyLinks();
});
