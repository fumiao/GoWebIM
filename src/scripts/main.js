var loc = window.location;
var uri = 'ws:';

if (loc.protocol === 'https:') {
    uri = 'wss:';
}
uri += '//' + loc.host;
//uri += loc.pathname + 'ws';
uri += '/ws?uuid=' + document.getElementById("uuid").value;

var Osend = document.getElementById("send")
var Omsg = document.getElementById("msg")
var Ost = document.getElementById("st")
var Oxiaoxi = document.getElementById("xiaoxi")
ws = new WebSocket(uri)
ws.onopen = function() {
    console.log('Connected')
}

ws.onmessage = function(evt) {
    Oxiaoxi.innerHTML += "<p>" + evt.data + "</p>"
}

Osend.onclick = function() {
    var msg = "{\"msg\": \"" + Omsg.value + "\",\"st\": \"" + Ost.value + "\"}"
    //ws.send(Omsg.value);
    ws.send(msg);
    Omsg.value = ""
    Omsg.focus()
}
