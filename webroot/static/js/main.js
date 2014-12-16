"use strict";

var conn;
var msg = document.getElementById("log");
var form = document.getElementById("form");

var ProtoBuf = dcodeIO.ProtoBuf;
var ChatMessage = ProtoBuf.loadProtoFile("protobuf/chat.proto").build("protobuf.ChatMessage");

form.addEventListener("submit", function(e) {
    e.preventDefault();

    if (conn && form.msg.value) {
        var player_message = new ChatMessage('user', form.msg.value);

        conn.send(player_message.toArrayBuffer());

        form.msg.value = "";
    } else {
        console.log("Could not send message connection or message is in invalid state");
    }

    return false;
});

function connect() {
    conn = new WebSocket("wss://" + window.location.host + "/ws");
    conn.binaryType = "arraybuffer";

    conn.onopen = function() {
        msg.innerHTML = "<h1>Connection opened</h1>";
    };

    conn.onmessage = function(evt) {
        var data = ChatMessage.decode(evt.data);
        var message;

        switch (data.getMessageType()) {
            case 1:
                message = '<strong style="color: green">User ' + data.getName() + " connected</strong>\n"
                break;
            case 2:
                message = '<strong style="color: red">User ' + data.getName() + " disconnected</strong>\n"
                break;
            default:
                message = "<strong>&lt;" + data.getName() + "&gt;</strong> " + data.getText() + "\n"
        }

        msg.innerHTML += message;
    };

    conn.onclose = function(evt) {
        setTimeout(connect, 1000);
        msg.innerHTML = "<h1>Connection closed.</h1>";
    };
}

connect();
