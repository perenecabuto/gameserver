var conn;
var msg = document.getElementById("log");
var form = document.getElementById("form");

form.addEventListener("submit", function(e) {
    e.preventDefault();

    if (conn && form.msg.value) {
        conn.send(form.msg.value);
        form.msg.value = "";
    } else {
        console.log("Could not send message connection or message is in invalid state");
    }

    return false;
});

function connect() {
    conn = new WebSocket("ws://localhost:3000/ws");

    conn.onopen = function() {
        msg.innerHTML = "<h1>Connection opened</h1>";
    };
    conn.onmessage = function(evt) {
        msg.innerHTML += evt.data + "\n";
    };
    conn.onclose = function(evt) {
        setTimeout(connect, 1000);
        msg.innerHTML = "<h1>Connection closed.</h1>";
    };
}

connect();
