import { useEffect, useState } from "react";
import { bind } from "@pojntfx/dudirekta";

export default () => {
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
    bind(
      () =>
        new WebSocket(
          new URLSearchParams(window.location.search).get("socketURL") ||
            "ws://localhost:1337"
        ),
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
      remote,
      setRemote
    );
  }, []);

  return (
    <main>
      <h1>Hydrapp Dudirekta Example</h1>

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
