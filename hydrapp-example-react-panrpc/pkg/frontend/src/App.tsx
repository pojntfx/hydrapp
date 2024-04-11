import { ILocalContext, IRemoteContext, Registry } from "@pojntfx/panrpc";
import { JSONParser } from "@streamparser/json-whatwg";
import { useEffect, useState } from "react";
import useAsyncEffect from "use-async";

class Local {
  async ExampleNotification(ctx: ILocalContext, msg: string) {
    if ("Notification" in window && Notification.permission !== "granted") {
      await Notification.requestPermission();
    }

    if ("Notification" in window) {
      new Notification(msg);
    } else {
      alert(msg);
    }
  }
}

class Remote {
  async ExamplePrintString(ctx: IRemoteContext, msg: string) {
    return;
  }

  async ExamplePrintStruct(ctx: IRemoteContext, input: any) {
    return;
  }

  async ExampleReturnError(ctx: IRemoteContext) {
    return;
  }

  async ExampleReturnString(ctx: IRemoteContext): Promise<string> {
    return "";
  }

  async ExampleReturnStruct(ctx: IRemoteContext): Promise<any> {
    return {};
  }

  async ExampleReturnStringAndError(ctx: IRemoteContext): Promise<string> {
    return "";
  }

  async ExampleCallback(ctx: IRemoteContext) {
    return;
  }

  async ExampleClosure(
    ctx: IRemoteContext,
    length: number,
    onIteration: (ctx: ILocalContext, i: number, b: string) => Promise<string>
  ): Promise<number> {
    return 0;
  }
}

const App = () => {
  const [clients, setClients] = useState(0);
  useEffect(() => console.log(clients, "clients connected"), [clients]);

  const [reconnect, setReconnect] = useState(false);
  const [registry] = useState(
    new Registry(
      new Local(),
      new Remote(),

      {
        onClientConnect: () => setClients((v) => v + 1),
        onClientDisconnect: () =>
          setClients((v) => {
            if (v === 1) {
              setReconnect(true);
            }

            return v - 1;
          }),
      }
    )
  );

  useAsyncEffect(async () => {
    if (reconnect) {
      await new Promise((r) => {
        setTimeout(r, 100);
      });

      setReconnect(false);

      return () => {};
    }

    const addr =
      new URLSearchParams(window.location.search).get("socketURL") ||
      "ws://localhost:1337";

    const socket = new WebSocket(addr);

    socket.addEventListener("error", (e) => {
      console.error("Disconnected with error, reconnecting:", e);

      setReconnect(true);
    });

    await new Promise<void>((res, rej) => {
      socket.addEventListener("open", () => res());
      socket.addEventListener("error", rej);
    });

    const encoder = new WritableStream({
      write(chunk) {
        socket.send(JSON.stringify(chunk));
      },
    });

    const parser = new JSONParser({
      paths: ["$"],
      separator: "",
    });
    const parserWriter = parser.writable.getWriter();
    const parserReader = parser.readable.getReader();
    const decoder = new ReadableStream({
      start(controller) {
        parserReader
          .read()
          .then(async function process({ done, value }) {
            if (done) {
              controller.close();

              return;
            }

            controller.enqueue(value?.value);

            parserReader
              .read()
              .then(process)
              .catch((e) => controller.error(e));
          })
          .catch((e) => controller.error(e));
      },
    });
    socket.addEventListener("message", (m) =>
      parserWriter.write(m.data as string)
    );
    socket.addEventListener("close", () => {
      parserReader.cancel();
      parserWriter.abort();
    });

    registry.linkStream(
      encoder,
      decoder,

      (v) => v,
      (v) => v
    );

    console.log("Connected to", addr);

    return () => socket.close();
  }, [reconnect]);

  return clients > 0 ? (
    <main>
      <h1>Hydrapp React and panrpc Example</h1>

      <div>
        <button
          onClick={() =>
            registry.forRemotes(async (_, remote) => {
              try {
                await remote.ExamplePrintString(
                  undefined,
                  prompt("String to print")!
                );
              } catch (e) {
                alert(JSON.stringify((e as Error).message));
              }
            })
          }
        >
          Print string
        </button>

        <button
          onClick={() =>
            registry.forRemotes(async (_, remote) => {
              try {
                await remote.ExamplePrintStruct(undefined, {
                  name: prompt("Name to print")!,
                });
              } catch (e) {
                alert(JSON.stringify((e as Error).message));
              }
            })
          }
        >
          Print struct
        </button>

        <button
          onClick={() =>
            registry.forRemotes(async (_, remote) => {
              try {
                await remote.ExampleReturnError(undefined);
              } catch (e) {
                alert(JSON.stringify((e as Error).message));
              }
            })
          }
        >
          Return error
        </button>

        <button
          onClick={() =>
            registry.forRemotes(async (_, remote) => {
              try {
                const res = await remote.ExampleReturnString(undefined);

                alert(JSON.stringify(res));
              } catch (e) {
                alert(JSON.stringify((e as Error).message));
              }
            })
          }
        >
          Return string
        </button>

        <button
          onClick={() =>
            registry.forRemotes(async (_, remote) => {
              try {
                const res = await remote.ExampleReturnStruct(undefined);

                alert(JSON.stringify(res));
              } catch (e) {
                alert(JSON.stringify((e as Error).message));
              }
            })
          }
        >
          Return struct
        </button>

        <button
          onClick={() =>
            registry.forRemotes(async (_, remote) => {
              try {
                await remote.ExampleReturnStringAndError(undefined);
              } catch (e) {
                alert(JSON.stringify((e as Error).message));
              }
            })
          }
        >
          Return string and error
        </button>

        <button
          onClick={() =>
            registry.forRemotes(async (_, remote) => {
              try {
                await remote.ExampleCallback(undefined);
              } catch (e) {
                alert(JSON.stringify((e as Error).message));
              }
            })
          }
        >
          Get three time notifications (with callback)
        </button>

        <button
          onClick={() =>
            registry.forRemotes(async (_, remote) => {
              try {
                await remote.ExampleClosure(undefined, 3, async (_, i, b) => {
                  if (
                    "Notification" in window &&
                    Notification.permission !== "granted"
                  ) {
                    await Notification.requestPermission();
                  }

                  if ("Notification" in window) {
                    new Notification(`In iteration ${i} ${b}`);
                  } else {
                    alert(`In iteration ${i} ${b}`);
                  }

                  return "This is from the frontend";
                });
              } catch (e) {
                alert(JSON.stringify((e as Error).message));
              }
            })
          }
        >
          Get three iteration notifications (with closure)
        </button>
      </div>
    </main>
  ) : (
    "Connecting ..."
  );
};

export default App;
