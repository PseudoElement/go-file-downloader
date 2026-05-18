async function apiReq() {
  return fetch("https://davidovich.online/api/v1/health/test-json", {
    method: "GET",
    headers: {
      Origin: "https://davidovich.online",
    },
    // body: JSON.stringify({
    //   room_name: "room_258",
    //   nickname: "pablus",
    // }),
  })
    .then((res) => res.json())
    .then((data) => console.log("DATA:", data))
    .catch((err) => console.log("ERROR:", err));
}

async function callTimesN(n) {
  for (let i = 0; i < n; i++) await apiReq();
}

(async () => {
  while (true) {
    callTimesN(1_000);
    await new Promise((res) => setTimeout(res, 2_000));
  }
})();
