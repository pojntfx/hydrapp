export const bind = (
  getSocket: () => WebSocket,
  local: any,
  remote: any,
  setRemote: (remote: any) => void
) => {
  const socket = getSocket();

  socket.addEventListener("open", () => {
    console.log("Connected to RPC server");
  });

  socket.addEventListener("error", (e) => {
    console.error("Got error from RPC server:", e);
  });

  socket.addEventListener("close", async () => {
    console.log("Disconnected from RPC server, reconnecting in 1s");

    await new Promise((res) => setTimeout(res, 1000));

    bind(getSocket, local, remote, setRemote);
  });

  const r = Object.assign({}, remote);
  for (let functionName in r) {
    (r as any)[functionName] = async (...args: any[]) => {
      return new Promise(async (res, rej) => {
        const id = Math.random().toString(16).slice(2);

        const handleReturn = ({ detail }: any) => {
          const [rv, err] = detail;

          if (err) {
            rej(new Error(err));
          } else {
            res(rv);
          }

          window.removeEventListener(`rpc:${id}`, handleReturn);
        };

        window.addEventListener(`rpc:${id}`, handleReturn);

        socket.send(JSON.stringify([true, id, functionName, args]));
      });
    };
  }
  setRemote(r);

  socket.addEventListener("message", async (e) => {
    const msg = JSON.parse(e.data);

    if (msg[0]) {
      const [_, id, functionName, args] = msg;

      try {
        const res = await (local as any)[functionName](...args);

        socket.send(JSON.stringify([false, id, res, ""]));
      } catch (e) {
        socket.send(JSON.stringify([false, id, "", (e as Error).message]));
      }
    } else {
      window.dispatchEvent(
        new CustomEvent(`rpc:${msg[1]}`, {
          detail: msg.slice(2),
        })
      );
    }
  });
};
