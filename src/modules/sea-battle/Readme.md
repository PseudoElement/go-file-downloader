# Scenario

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

`player1_id: K1*+,K2.,O3.___player2_id: E7.,E4.___`
`K1+ - has ship`
`K1+* - striked ship`
`K1. - empty cell`


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


