// InfraPilot Enterprise - Landing Page Script

// Smooth scrolling
document.querySelectorAll('a[href^="#"]').forEach(anchor => {
    anchor.addEventListener('click', function (e) {
        e.preventDefault();
        const target = document.querySelector(this.getAttribute('href'));
        if (target) target.scrollIntoView({ behavior: 'smooth', block: 'start' });
    });
});

// Copy to clipboard
function copyToClipboard(text) {
    navigator.clipboard.writeText(text).then(() => {
        const btn = event.target;
        btn.textContent = 'Copied!';
        btn.style.background = '#10b981';
        setTimeout(() => { btn.textContent = 'Copy'; btn.style.background = ''; }, 2000);
    });
}

// Navbar scroll effect
window.addEventListener('scroll', () => {
    const navbar = document.querySelector('.navbar');
    navbar.style.boxShadow = window.pageYOffset > 100
        ? '0 4px 12px rgba(0,0,0,0.15)'
        : '0 2px 4px rgba(0,0,0,0.1)';
});

// Intersection Observer for animations
const observer = new IntersectionObserver((entries) => {
    entries.forEach(entry => {
        if (entry.isIntersecting) {
            entry.target.style.opacity = '1';
            entry.target.style.transform = 'translateY(0)';
        }
    });
}, { threshold: 0.1 });

document.querySelectorAll('.feature-card, .pricing-card, .doc-card, .download-card').forEach(card => {
    card.style.opacity = '0';
    card.style.transform = 'translateY(20px)';
    card.style.transition = 'opacity 0.6s ease, transform 0.6s ease';
    observer.observe(card);
});

// Dashboard preview parallax
const preview = document.querySelector('.dashboard-preview');
if (preview) {
    let mx = 0, my = 0, cx = 0, cy = 0;
    document.addEventListener('mousemove', (e) => {
        const r = preview.getBoundingClientRect();
        mx = (e.clientX - r.left - r.width/2) / 50;
        my = (e.clientY - r.top - r.height/2) / 50;
    });
    function animate() {
        cx += (mx - cx) * 0.1;
        cy += (my - cy) * 0.1;
        preview.style.transform = `perspective(1000px) rotateY(${-5 + cx}deg) rotateX(${5 - cy}deg)`;
        requestAnimationFrame(animate);
    }
    animate();
}

// Trigger initial animations
window.addEventListener('DOMContentLoaded', () => {
    setTimeout(() => {
        document.querySelectorAll('.feature-card, .pricing-card, .doc-card, .download-card')
            .forEach((card, i) => setTimeout(() => {
                card.style.opacity = '1';
                card.style.transform = 'translateY(0)';
            }, i * 100));
    }, 300);
});

console.log('%c InfraPilot Enterprise v1.0 ', 'background: #2563eb; color: white; padding: 10px; font-size: 16px; font-weight: bold;');