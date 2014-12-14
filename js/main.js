"use strict";

var conn;
var msg = document.getElementById("log");
var form = document.getElementById("form");

var ProtoBuf = dcodeIO.ProtoBuf;
var Chat;

ProtoBuf.loadProtoFile("protobuf/chat.proto", function(err, builder) {
    Chat = builder.build('protobuf.Chat');
});

form.addEventListener("submit", function(e) {
    e.preventDefault();

    if (conn && form.msg.value) {
        var player_message = new Chat('user', form.msg.value);
        var buffer = player_message.encode();

        conn.send(buffer.toArrayBuffer());

        form.msg.value = "";
    } else {
        console.log("Could not send message connection or message is in invalid state");
    }

    return false;
});

function connect() {
    conn = new WebSocket("ws://localhost:4000/ws");
    conn.binaryType = "arraybuffer";

    conn.onopen = function() {
        msg.innerHTML = "<h1>Connection opened</h1>";
    };

    conn.onmessage = function(evt) {
        var message = Chat.decode(evt.data);
        msg.innerHTML += message.name + ": " + message.text + "\n";
    };

    conn.onclose = function(evt) {
        setTimeout(connect, 1000);
        msg.innerHTML = "<h1>Connection closed.</h1>";
    };
}

connect();
