# Scenario

`ON SERVER INIT LOAD ALL EXISTING ROOMS IF EXISTS`

1. Create room with unique id and name
2. Connect to room by name
3. Send room name to friend
4. Friend connects to room by name

=====

# On reload page

1. Find all rooms of player by player name
2. Connect to selected room by room id
3. Listen socket

=====

# SeaBattle ships posiitions example:

`player1_id: K1+*,K2.,O3.___player2_id: E7.,E4.___`
`K1+ - has ship`
`K1+* - striked ship`
`K1. - missed shot`


Front: 
User_1 - created ships on map.
User_2 - created ships on map.
Enemy field is empty.
User_1 sent his positions to back.
User_2 sent his positions to back.
Random definer who begins game.
In room are contained info about both players.
Then Users select cell one after another.
Server checks if choice stroke enemy ship.
Then sends info to both client NewStepMsgResp{}.
Frontend fields show updated info.


==========

Create new room - call _createRoomController, save room in map and db, add player in room.players, connect user to socket.
Connect to existing room - call _connectToRoomWsController, add player in room.players, connect user to socket, send ConnectPlayerResp to client.

Set player READY - send msg to socket, set player.positions, save all new positions in db, send PlayerReadyResp msg to client.
Make step in game - send msg to socket, update player position in room.players[player.id], save all new positions in db, send PlayerStepResp msg to client.

===========

Requirements ro ROOM name: minLen(1), required, unique.


===========
GetRoomsList resp :
resp: {
	rooms: {
		room_id_1: {room_name: "room_1", room_id: "room_id_1", created_at: "12 Dec 2025", players: [Player, Player]},
		room_id_2: {room_name: "room_2", room_id: "room_id_2", created_at: "13 Dec 2025", players: [Player, Player]}
	}
}
