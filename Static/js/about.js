const appStore = document.getElementById('app-store');
const taptap = document.getElementById('taptap');
const googlePlay = document.getElementById('google-play');
const steam = document.getElementById('steam');
const wegame = document.getElementById('wegame');
const nintendoSwitch = document.getElementById('nintendo-switch');

const github = document.getElementById('github');
const afdian = document.getElementById('afdian');
const donation = document.getElementById('donation');

const qqGroup = document.getElementById('qq-group');

appStore.addEventListener('click', () => {
    window.open('https://itunes.apple.com/cn/app/muse-dash/id1361473095?mt=8');
})

taptap.addEventListener('click', () => {
    window.open('https://www.taptap.com/app/60809');
})

googlePlay.addEventListener('click', () => {
    window.open('https://play.google.com/store/apps/details?id=com.prpr.musedash');
})

steam.addEventListener('click', () => {
    window.open('https://store.steampowered.com/app/774171/Muse_Dash/');
})

wegame.addEventListener('click', () => {
    window.open('https://www.wegame.com.cn/store/2000902/MuseDash');
})

nintendoSwitch.addEventListener('click', () => {
    window.open('https://ec.nintendo.com/JP/ja/titles/70010000014537');
})

github.addEventListener('click', () => {
    window.open('https://github.com/shallow-zzt/MuseDashInfo');
})

afdian.addEventListener('click', () => {
    window.open('https://afdian.com/a/musedashinfotools');
})

donation.addEventListener('click', () => {
    const donationPic = document.getElementById('donation-pic');
    if(donationPic.style.display == 'none') {
        donationPic.style.display = 'block';
    } else {
        donationPic.style.display = 'none';
    }
})

qqGroup.addEventListener('click', () => {
    window.open('https://qm.qq.com/cgi-bin/qm/qr?k=dbWvJ8g51rIY0z3vaQh6AK7GB2KKnFnf&jump_from=webapi&authKey=uSmzZZAmsCqcLOdoiRs+WxwChnyfslzcmJBUG42GzHSOMqfZuL++/M1BrNBfmyHk');
})
