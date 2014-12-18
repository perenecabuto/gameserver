"use strict";

var conn;

var ProtoBuf = dcodeIO.ProtoBuf;
var ChatMessage = ProtoBuf.loadProtoFile("/protobuf/chat.proto").build('protobuf.ChatMessage');
var GameMessage = ProtoBuf.loadProtoFile("/protobuf/game.proto").build('protobuf.GameMessage');

var Game = {
    init: function () {
        this.connect();
    },

    connect: function() {
        var that = this;
        this.conn = new WebSocket("ws://" + location.host + "/ws/game");
        this.conn.binaryType = "arraybuffer";

        this.conn.onopen = function(evt) {
            //that.msg.innerHTML = "<h1>New User in game</h1>";

            //setInterval(function() { npc.pos.x -= 5; }, 100);
            var buffer= new ChatMessage('user', "A new player has connected").encode();
            that.conn.send(buffer.toArrayBuffer());
        };

        this.conn.onmessage = function(evt) {
            //var message = GameMessage.decode(evt.data);
            //console.log(message);
            var player = new game.PlayerEntity(parseInt(Math.random() * 1000), 10);
            me.game.world.addChild(player, 4);
            me.game.world.sort(true);
            that.player = player;
        };

        this.conn.onclose = function(evt) {
            setTimeout(that.connect, 1000);
            me.game.world.removeChild(that.player);
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

Chat.init();
Game.init();
