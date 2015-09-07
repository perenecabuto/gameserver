var GameConnection = {
    init: function (game) {
        var that = this;
        this.game = game;
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
            var remotePlayer;

            if (message.id === that.playerId) return;

            if (that.game.playerPool) {
                remotePlayer = that.game.playerPool.filter(function(player) { return player.id == message.id }, true).first;
            }

            switch (message.action) {
                case GameMessage.Action.CREATE:
                    that.playerId = message.id;
                    player = that.game.playerPool.getFirstExists(false);
                    player.id = that.playerId;
                    player.reset(message.position.x, message.position.y);
                    player.action = GameMessage.Action.STOP;

                    that.game.input.keyboard.onDownCallback = function(e) {
                        if (that.game.input.keyboard.isDown(Phaser.Keyboard.SPACEBAR)) {
                            player.action = GameMessage.Action.JUMP;
                        } else if (that.game.input.keyboard.isDown(Phaser.Keyboard.LEFT)) {
                            player.action = GameMessage.Action.MOVE_LEFT;
                        } else if (that.game.input.keyboard.isDown(Phaser.Keyboard.RIGHT)) {
                            player.action = GameMessage.Action.MOVE_RIGHT;
                        } else if (!player.jumping) {
                            player.action = GameMessage.Action.STOP;
                        }


                        if (!that.sendingPosition && player.action != GameMessage.Action.STOP) {
                            console.log('MOUSE DOWN');
                            that.sendingPosition = true;
                            that.send(player.id, player.action, player.body.x, player.body.y);
                        }
                    };
                    that.game.input.keyboard.onUpCallback = function() {
                        console.log('MOUSE UP');
                        that.sendingPosition = false;
                        player.action = GameMessage.Action.STOP;
                        that.send(player.id, player.action, player.x, player.y);
                    };
                    break;
                case GameMessage.Action.SPAWN:
                    if (!remotePlayer) {
                        remotePlayer = that.game.playerPool.getFirstExists(false);
                        remotePlayer.id = message.id;
                    }
                    remotePlayer.action = GameMessage.Action.STOP;
                    remotePlayer.reset(message.position.x, message.position.y);
                    break;
                case GameMessage.Action.DIE:
                    if (remotePlayer) {
                        remotePlayer.kill();
                    }
                    break;
                case GameMessage.Action.STOP:
                    if (remotePlayer) {
                        remotePlayer.action = message.action;
                        remotePlayer.x = message.position.x;
                        remotePlayer.y = message.position.y;
                    }
                    break;
                default:
                    if (remotePlayer) remotePlayer.action = message.action;
            }
        };

        this.conn.onclose = function(evt) {
            that.send(player.id, GameMessage.Action.DIE, player.x, player.y);
            setTimeout(function() { that.connect(); }, 1000);
        };
    },

    send: function(id, action, x, y) {
        var buffer = new GameMessage({
            id: id,
            action: action,
            position: {x: x, y: y}
        }).encode();

        console.log("Sengind ID: ", player.id, "action", player.action);
        this.conn.send(buffer.toArrayBuffer());
    }
};
