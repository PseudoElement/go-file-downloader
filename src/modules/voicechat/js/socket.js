function testSocket(roomId) {
  const ws = new WebSocket(
    `http://localhost:8080/api/v1/voicechat/connect?room_id=${roomId}&peer_name=sintol`,
  );

  ws.onopen = () => {
    console.log("onopen called");
    ws.send(
      JSON.stringify({
        action: "CONNECT",
        data: {
          peer_name: "sintol",
          peer_descriptor: "descriptor_RTC",
        },
      }),
    );
  };
  ws.onmessage = (event) => {
    console.log("Message from server: ", event.data);
  };

  return ws;
}

const ws = testSocket("");

ws.send(
  JSON.stringify({
    action: "CONNECT",
    data: {
      peer_name: "sintol",
      peer_descriptor: "descp_SINTOL",
    },
  }),
);
