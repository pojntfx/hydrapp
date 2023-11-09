import { useEffect, useState } from "react";
import { linkWebSocket } from "@pojntfx/dudirekta";

const App = () => {
  const [remote, setRemote] = useState({
    ExamplePrintString: async (msg: string) => {},
    ExamplePrintStruct: async (input: any) => {},
    ExampleReturnError: async () => {},
    ExampleReturnString: async () => "",
    ExampleReturnStruct: async (): Promise<any> => {},
    ExampleReturnStringAndError: async () => "",
    ExampleNotification: async () => {},
  });

  useEffect(() => {
    (async () => {
      const socket = new WebSocket(
        new URLSearchParams(window.location.search).get("socketURL") ||
          "ws://localhost:1337"
      );

      socket.addEventListener("close", (e) => {
        console.error("Disconnected with error:", e.reason);
      });

      await new Promise<void>((res, rej) => {
        socket.addEventListener("open", () => res());
        socket.addEventListener("error", rej);
      });

      setRemote(
        linkWebSocket(
          socket,
          {
            ExampleNotification: async (msg: string) => {
              if (
                "Notification" in window &&
                Notification.permission !== "granted"
              ) {
                await Notification.requestPermission();
              }

              if ("Notification" in window) {
                new Notification(msg);
              } else {
                alert(msg);
              }
            },
          },
          {
            ExamplePrintString: async (msg: string) => {},
            ExamplePrintStruct: async (input: any) => {},
            ExampleReturnError: async () => {},
            ExampleReturnString: async () => "",
            ExampleReturnStruct: async (): Promise<any> => {},
            ExampleReturnStringAndError: async () => "",
            ExampleNotification: async () => {},
          },

          1000 * 10,

          JSON.stringify,
          JSON.parse,

          (v) => v,
          (v) => v
        )
      );
    })();
  }, []);

  return (
    <main>
      <h1>Hydrapp Dudirekta-Parcel Example</h1>

      <div>
        <button
          onClick={async () => {
            await remote.ExamplePrintString(prompt("String to print")!);
          }}
        >
          Print string
        </button>

        <button
          onClick={async () => {
            await remote.ExamplePrintStruct({
              name: prompt("Name to print")!,
            });
          }}
        >
          Print struct
        </button>

        <button
          onClick={async () => {
            try {
              await remote.ExampleReturnError();
            } catch (e) {
              alert(JSON.stringify((e as Error).message));
            }
          }}
        >
          Return error
        </button>

        <button
          onClick={async () => {
            const res = await remote.ExampleReturnString();

            alert(JSON.stringify(res));
          }}
        >
          Return string
        </button>

        <button
          onClick={async () => {
            const res = await remote.ExampleReturnStruct();

            alert(JSON.stringify(res));
          }}
        >
          Return struct
        </button>

        <button
          onClick={async () => {
            try {
              await remote.ExampleReturnStringAndError();
            } catch (e) {
              alert(JSON.stringify((e as Error).message));
            }
          }}
        >
          Return string and error
        </button>

        <button
          onClick={async () => {
            await remote.ExampleNotification();
          }}
        >
          Get three notifications
        </button>
      </div>
    </main>
  );
};

export default App;
