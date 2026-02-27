// Extra JavaScript for Straw documentation

// Add copy button functionality for code blocks without it
document.addEventListener('DOMContentLoaded', function() {
  // Ensure all code blocks have copy buttons
  const codeBlocks = document.querySelectorAll('pre > code');
  codeBlocks.forEach(function(codeBlock) {
    const pre = codeBlock.parentElement;
    if (!pre.querySelector('.md-clipboard')) {
      // Create copy button if Material's clipboard isn't present
      const button = document.createElement('button');
      button.className = 'md-clipboard md-icon';
      button.title = 'Copy to clipboard';
      button.innerHTML = '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><path d="M19 21H8V7h11m0-2H8a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h11a2 2 0 0 0 2-2V7a2 2 0 0 0-2-2m-3-4H4a2 2 0 0 0-2 2v14h2V3h12V1z"/></svg>';
      
      button.addEventListener('click', function() {
        navigator.clipboard.writeText(codeBlock.textContent).then(function() {
          button.classList.add('md-clipboard--active');
          setTimeout(function() {
            button.classList.remove('md-clipboard--active');
          }, 1000);
        });
      });
      
      pre.appendChild(button);
    }
  });
});

// Add keyboard navigation enhancements
document.addEventListener('keydown', function(e) {
  // Press '?' to open search
  if (e.key === '?' && !e.ctrlKey && !e.metaKey && !e.altKey) {
    const searchInput = document.querySelector('.md-search__input');
    if (searchInput) {
      e.preventDefault();
      searchInput.focus();
    }
  }
  
  // Press 'g' then 'h' to go home
  if (e.key === 'g') {
    const keyBuffer = window.keyBuffer || [];
    keyBuffer.push('g');
    window.keyBuffer = keyBuffer;
    
    setTimeout(function() {
      window.keyBuffer = [];
    }, 500);
    
    if (keyBuffer.length === 2 && keyBuffer[1] === 'h') {
      e.preventDefault();
      window.location.href = '/';
    }
  }
});

// Add smooth scrolling for anchor links
document.querySelectorAll('a[href^="#"]').forEach(anchor => {
  anchor.addEventListener('click', function (e) {
    e.preventDefault();
    const target = document.querySelector(this.getAttribute('href'));
    if (target) {
      target.scrollIntoView({
        behavior: 'smooth',
        block: 'start'
      });
    }
  });
});

// Track page views for analytics (if implemented)
if (typeof gtag !== 'undefined') {
  gtag('config', 'GA_TRACKING_ID', {
    page_path: window.location.pathname
  });
}
