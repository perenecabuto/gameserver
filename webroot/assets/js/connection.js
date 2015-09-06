var GameConnection = {
    init: function (game) {
        var that = this;
        this.game = game;
        this.players = {};

        $(document).on("player-action", function(e, action, pos) {
            if (!that.game.player) {
                console.log("No player found to send position: ", action, pos);
                return;
            };

            var buffer = new GameMessage({
                id: that.game.player.id,
                action: action,
                position: {x: parseInt(pos.x), y: parseInt(pos.y)}
            }).encode();

            console.log("Sending position: ", action, pos);
            that.conn.send(buffer.toArrayBuffer());
        });

        this.connect();
    },

    connect: function() {
        var that = this;
        this.conn = new WebSocket("ws://" + location.host + "/ws/game");
        this.conn.binaryType = "arraybuffer";

        this.conn.onopen = function(evt) {
            Chat.send("here comes a new challenger!");
        };

        this.conn.onmessage = function(evt) {
            var message = GameMessage.decode(evt.data);
            if (that.game.player && message.id === that.game.player.name) return;

            switch (message.action) {
                case GameMessage.Action.CREATE:
                    that.game.player.id = message.id;
                    break;
                case GameMessage.Action.SPAWN:
                    // CREATE REMOTE PLAYER HERE
                    //var player = new game.PlayerEntity(message.id, message.position.x, message.position.y);
                    //me.game.world.addChild(player, 4);

                    break;
                default:
                    // UPDATE PLAYER POSITION HERE
                    //var player = me.game.world.getChildByName(message.id)[0];
                    //console.log('player', player, 'action', message.action, 'position', message.position);
                    //player.updateAction(message.action, message.position);

                    break;
            }
        };

        this.conn.onclose = function(evt) {
            Chat.send("! lost game connection");
            setTimeout(function() { that.connect(); }, 1000);
        };
    }
};

var Chat = {
    init: function() {
        var that = this;

        var form = document.getElementById("form");
        form.addEventListener("submit", function(e) {
            e.preventDefault();

            if (that.conn && form.msg.value) {
                var player_message = new ChatMessage('user', form.msg.value);
                var buffer = player_message.encode();

                that.conn.send(buffer.toArrayBuffer());
                form.msg.value = "";
            } else {
                console.log("Could not send message connection or message is in invalid state");
            }

            return false;
        });

        this.connect();
    },

    connect: function() {
        var that = this;
        this.conn = new WebSocket("ws://" + location.host + "/ws/chat");
        this.conn.binaryType = "arraybuffer";

        this.conn.onopen = function() {
            that.send("<h1>Connection opened</h1>");
        };

        this.conn.onmessage = function(evt) {
            var message = ChatMessage.decode(evt.data);
            that.send(message.name + ": " + message.text);
        };

        this.conn.onclose = function(evt) {
            setTimeout(that.connect, 1000);
            that.send("! Connection closed.");
        };
    },

    send: function(message) {
        this.msg = this.msg || document.getElementById("messages");
        this.msg.innerHTML += message + "\n";
    }
};
