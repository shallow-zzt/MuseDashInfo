function toggleMenu() {
    document.querySelector('.navbar').classList.toggle('menu-open');
    document.querySelector('.overlay').classList.toggle('show');
}

function closeMenu() {
    document.querySelector('.navbar').classList.remove('menu-open');
    document.querySelector('.overlay').classList.remove('show');
}