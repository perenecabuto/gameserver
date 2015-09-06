var ProtoBuf = dcodeIO.ProtoBuf;
var ChatMessage = ProtoBuf.loadProtoFile("/protobuf/chat.proto").build('protobuf.ChatMessage');
var GameMessage = ProtoBuf.loadProtoFile("/protobuf/game.proto").build('protobuf.GameMessage');

(function() {
    'use strict';

    var game = new Game().onComplete(function(game) {
        Chat.init();
        GameConnection.init(game);
    });

    var phaser = new Phaser.Game(800, 400, Phaser.AUTO, '', game);
})();
