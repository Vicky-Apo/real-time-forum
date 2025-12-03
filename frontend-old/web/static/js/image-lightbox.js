	function openLightbox(img) {
    const lightbox = document.getElementById('image-lightbox');
    const lightboxImg = lightbox.querySelector('img');
    lightboxImg.src = img.src;
    lightbox.classList.add('active');
}

function closeLightbox() {
    const lightbox = document.getElementById('image-lightbox');
    lightbox.classList.remove('active');
    lightbox.querySelector('img').src = '';
}

// Allow closing by clicking outside the image
document.addEventListener('DOMContentLoaded', function() {
    const lightbox = document.getElementById('image-lightbox');
    if (lightbox) {
        lightbox.addEventListener('click', function(e) {
            if (e.target === lightbox) {
                closeLightbox();
            }
        });
    }
    // Optional: Escape key closes lightbox
    document.addEventListener('keydown', function(e) {
        if (e.key === "Escape") {
            closeLightbox();
        }
    });
});