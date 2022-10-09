import { v4 } from "uuid";

const bind = (getSocket: () => WebSocket, broker: EventTarget) =>
  new Promise<void>((res, rej) => {
    const getEventName = (id: string) => `rpc:${id}`;

    const socket = getSocket();

    socket.addEventListener("open", () => {
      console.log("Connected to RPC server");

      res();
    });

    socket.addEventListener("error", (e) => {
      console.error("Got error from RPC server:", e);

      rej(e);
    });

    socket.addEventListener("close", async () => {
      console.log("Disconnected from RPC server, reconnecting");

      await bind(getSocket, broker);
    });

    let firstMessage = true;
    socket.addEventListener("message", (e) => {
      const msg = JSON.parse(e.data) as any[];

      if (firstMessage) {
        firstMessage = false;

        (msg as string[]).forEach((name) => {
          (window as any)[name] = async (...args: any[]) => {
            const id = v4();

            socket.send(JSON.stringify([id, name, args]));

            return new Promise<void>((res, rej) => {
              const handleResponse = (e: any) => {
                const rv = (e as CustomEvent).detail;

                if (rv[1] === "") {
                  if (rv[0] === null) {
                    res();
                  } else {
                    res(rv[0]);
                  }
                } else {
                  rej(rv[1]);
                }

                broker.removeEventListener(getEventName(id), handleResponse);
              };

              broker.addEventListener(getEventName(id), handleResponse);
            });
          };
        });

        return;
      }

      broker.dispatchEvent(
        new CustomEvent(getEventName(msg[0]), {
          detail: msg.slice(1),
        })
      );
    });
  });

export default bind;
