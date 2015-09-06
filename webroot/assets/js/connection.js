var GameConnection = {
    init: function (game) {
        var that = this;
        this.game = game;
        this.players = {};

        document.addEventListener("player-action", function(e) {
            var action = e.detail.action;
            var pos = e.detail.pos;
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
            Chat.write("here comes a new challenger!");
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
