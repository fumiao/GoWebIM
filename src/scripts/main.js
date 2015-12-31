var loc = window.location;
var uri = 'ws:';

if (loc.protocol === 'https:') {
    uri = 'wss:';
}
uri += '//' + loc.host;
//uri += loc.pathname + 'ws';
uri += '/ws';

ws = new WebSocket(uri)
ws.onopen = function() {
    console.log('Connected')
}

ws.onmessage = function(evt) {
    /*var out = document.getElementById('output');
    out.innerHTML += evt.data + '<br>';*/
    console.log(evt.data)
}

/*setInterval(function() {
    ws.send('Hello, Server!');
}, 1000);*/

var Osend = document.getElementById("send")
var Omsg = document.getElementById("msg")
Osend.onclick = function() {
    ws.send(Omsg.value);
    Omsg.value = ""
    Omsg.focus()
}
