"use strict";

var conn;

var ProtoBuf = dcodeIO.ProtoBuf;
var ChatMessage = ProtoBuf.loadProtoFile("/protobuf/chat.proto").build('protobuf.ChatMessage');
var GameMessage = ProtoBuf.loadProtoFile("/protobuf/game.proto").build('protobuf.GameMessage');

var Game = {
    init: function () {
        var that = this;
        this.msg = document.getElementById("log");
        var form = document.getElementById("form");

        this.players = {};

        this.connect();

        me.event.subscribe("mainPlayerMovement", function(x, y) {
            var buffer = new GameMessage({
                id: that.player.name,
                action: GameMessage.Action.MOVING,
                position: {x: parseInt(that.player.pos.x), y: parseInt(that.player.pos.y)}
            }).encode();

            that.conn.send(buffer.toArrayBuffer());
        })
    },

    connect: function() {
        var that = this;
        this.conn = new WebSocket("ws://" + location.host + "/ws/game");
        this.conn.binaryType = "arraybuffer";

        this.conn.onopen = function(evt) {
            that.msg.innerHTML += "here comes a new challenger!\n";
        };

        this.conn.onmessage = function(evt) {
            var message = GameMessage.decode(evt.data);
            if (that.player && message.id === that.player.name) return;

            switch (message.action) {
                case GameMessage.Action.NEW:
                    that.player = me.game.world.getChildByName('mainplayer')[0];
                    that.player.name = message.id;
                    break;
                case GameMessage.Action.SPAWN:
                    var player = new game.PlayerEntity(message.id, message.position.x, message.position.y);
                    me.game.world.addChild(player, 4);
                    break;
                case GameMessage.Action.MOVING:
                    console.log('me.game.world.getChildByName(' + message.id + ')[0]');
                    var player = me.game.world.getChildByName(message.id)[0];

                    // TODO OMGWTFBBQFIXTHIS
                    // Move eternamente para um lado ou para o outro, para preservar as animações
                    // Talvez deva ser interceptado pelo server (unmarshall), e envie apenas: "moving" e "stop" quando for para parar
                    if (player.pos.x - message.position.x > 0) {
                        player.moveLeft();
                        // player.standStill();
                    }
                    else {
                        player.moveRight();
                        // player.standStill();
                    }

                    // player.pos.x = message.position.x;
                    // player.pos.y = message.position.y;

                    // player.body.vel.x = 0;
                    // player.body.vel.y = 0;
                    me.game.world.sort(true);
                    break;
                case GameMessage.Action.DEAD:
                    var player = me.game.world.getChildByName(message.id)[0];
                    me.game.world.removeChild(player);
                    break;
            }
        };

        this.conn.onclose = function(evt) {
            setTimeout(function() { that.connect(); }, 1000);
        };
    }
};

var Chat = {

    init: function() {
        var that = this;
        this.msg = document.getElementById("log");
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
            that.msg.innerHTML = "<h1>Connection opened</h1>";
        };

        this.conn.onmessage = function(evt) {
            var message = ChatMessage.decode(evt.data);
            that.msg.innerHTML += message.name + ": " + message.text + "\n";
        };

        this.conn.onclose = function(evt) {
            setTimeout(that.connect, 1000);
            that.msg.innerHTML = "<h1>Connection closed.</h1>";
        };
    }
};
