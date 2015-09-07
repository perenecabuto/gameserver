var Chat = {
    init: function() {
        this.playerName = window.prompt("Player name:", localStorage.getItem('playerName'));
        localStorage.setItem('playerName', this.playerName);
        this.connect();
    },

    connect: function() {
        var that = this;
        this.conn = new WebSocket("ws://" + location.host + "/ws/chat");
        this.conn.binaryType = "arraybuffer";

        this.conn.onopen = function() {
            that.doBinds();
            that.write("<h1>Connection opened</h1>");
        };

        this.conn.onmessage = function(evt) {
            var message = ChatMessage.decode(evt.data);
            that.write(message.name + ": " + message.text);
        };

        this.conn.onclose = function(evt) {
            setTimeout(that.connect, 1000);
            that.write("! Connection closed.");
        };
    },

    write: function(message) {
        this.msg.innerHTML += message + "\n";
    },

    send: function(label, value) {
        var message = new ChatMessage(label, value);
        var buffer = message.encode();
        this.conn.send(buffer.toArrayBuffer());
    },

    doBinds: function() {
        var that = this;
        this.msg =  document.getElementById("messages");
        this.form = document.getElementById("form");
        this.form.addEventListener("submit", function(e) {
            e.preventDefault();

            if (!form.msg.value.match(/^\s*$/)) {
                Chat.send(that.playerName, form.msg.value);
            }

            form.msg.value = "";
            return false;
        });
    }
};
