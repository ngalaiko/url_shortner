// js for https://shortner.galako.rocks 
function getShortUrl() {
    var payload = {
        url: document.getElementById('inputUrl').value
    };
    fetch("/link",
        {
            method: "POST",
            credentials: 'include',
            body: JSON.stringify(payload)
        })
        .then(function(res){ return res.json(); })
        .then(showShortUrl)
}
function showShortUrl(data) {
    if (!data.ok) {
        alert(data.err);
        return;
    }

    var link = "//" + location.hostname + '/' + data.data.short_url;
    document.getElementById('shortUrl').value = "https:" + link ;
}
function logout() {
    fetch("/logout",
        {
            method: "GET",
            credentials: 'include'
        })
        .then(function(res){ window.location.reload();  })
}

