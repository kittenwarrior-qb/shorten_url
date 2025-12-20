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

// Display analytics
function displayAnalytics(analytics) {
    const analyticsEl = document.getElementById('analytics');
    
    if (!analytics || analytics.total_clicks === 0) {
        analyticsEl.style.display = 'none';
        return;
    }
    
    let html = '<h4>Analytics</h4>';
    html += `<div class="analytics-stat"><strong>Total Clicks:</strong> ${analytics.total_clicks}</div>`;
    
    // Countries
    if (analytics.countries && Object.keys(analytics.countries).length > 0) {
        html += '<div class="analytics-section"><strong>Countries:</strong><ul>';
        Object.entries(analytics.countries)
            .sort((a, b) => b[1] - a[1])
            .forEach(([country, count]) => {
                html += `<li>${country}: ${count}</li>`;
            });
        html += '</ul></div>';
    }
    
    // Devices
    if (analytics.devices && Object.keys(analytics.devices).length > 0) {
        html += '<div class="analytics-section"><strong>Devices:</strong><ul>';
        Object.entries(analytics.devices)
            .sort((a, b) => b[1] - a[1])
            .forEach(([device, count]) => {
                html += `<li>${device}: ${count}</li>`;
            });
        html += '</ul></div>';
    }
    
    // Browsers
    if (analytics.browsers && Object.keys(analytics.browsers).length > 0) {
        html += '<div class="analytics-section"><strong>Browsers:</strong><ul>';
        Object.entries(analytics.browsers)
            .sort((a, b) => b[1] - a[1])
            .forEach(([browser, count]) => {
                html += `<li>${browser}: ${count}</li>`;
            });
        html += '</ul></div>';
    }
    
    // OS
    if (analytics.os && Object.keys(analytics.os).length > 0) {
        html += '<div class="analytics-section"><strong>Operating Systems:</strong><ul>';
        Object.entries(analytics.os)
            .sort((a, b) => b[1] - a[1])
            .forEach(([os, count]) => {
                html += `<li>${os}: ${count}</li>`;
            });
        html += '</ul></div>';
    }
    
    // Referers
    if (analytics.referers && Object.keys(analytics.referers).length > 0) {
        html += '<div class="analytics-section"><strong>Referers:</strong><ul>';
        Object.entries(analytics.referers)
            .sort((a, b) => b[1] - a[1])
            .slice(0, 5) // Top 5
            .forEach(([referer, count]) => {
                html += `<li>${referer || 'Direct'}: ${count}</li>`;
            });
        html += '</ul></div>';
    }
    
    analyticsEl.innerHTML = html;
    analyticsEl.style.display = 'block';
}

// Show link detail with QR code and analytics
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
        const analytics = data.data.analytics;
        
        document.getElementById('shortUrl').value = link.short_url;
        
        const qrCodeEl = document.getElementById('qrCode');
        if (link.qr_code) {
            qrCodeEl.innerHTML = `<img src="${link.qr_code}" alt="QR Code">`;
        }
        
        // Display analytics
        displayAnalytics(analytics);
        
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
