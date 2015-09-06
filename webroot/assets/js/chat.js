var Chat = {
    init: function() {
        var that = this;
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

    send: function(message) {
        this.conn.send(message);
    },

    doBinds: function() {
        this.msg =  document.getElementById("messages");
        this.form = document.getElementById("form");
        this.form.addEventListener("submit", function(e) {
            e.preventDefault();

            if (!form.msg.value.match(/^\s*$/)) {
                var player_message = new ChatMessage('user', form.msg.value);
                var buffer = player_message.encode();
                Chat.send(buffer.toArrayBuffer());
            }

            form.msg.value = "";
            return false;
        });
    }
};
