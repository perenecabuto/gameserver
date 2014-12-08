"use strict";

var ProtoBuf = dcodeIO.ProtoBuf;

ProtoBuf.loadProtoFile("protobuf/object.proto", function(err, builder) {
    var Obj = builder.build('protobuf.Object');
    var Position = Obj.Position;
    var Status = Obj.Status;

    var player = new Obj(1, new Position(0, 0), Status.IDLE);
    console.log(player, player.encode(), player.toArrayBuffer());

    var board = document.getElementById("board");

    function updatePlayer() {
        player.position.x = parseInt(Math.random() * 1000);
        player.position.y = parseInt(Math.random() * 1000);

        var decodedPlayer = Obj.decode(player.toArrayBuffer());
        //console.log(decodedPlayer, decodedPlayer.encode(), decodedPlayer.toArrayBuffer());

        board.innerHTML = JSON.stringify(player.toRaw());

        setTimeout(updatePlayer, 1000);
    }

    updatePlayer();
});
