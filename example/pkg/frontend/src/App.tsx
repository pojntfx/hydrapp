import bind from "./bind";

await bind(
  () =>
    new WebSocket(
      new URLSearchParams(window.location.search).get("socketURL") ||
        "ws://localhost:1337"
    ),
  window
);

export default () => {
  return (
    <main>
      <h1>Hydrapp Example</h1>

      <div>
        <button
          onClick={async () => {
            await examplePrintString(prompt("String to print")!);
          }}
        >
          Print string
        </button>

        <button
          onClick={async () => {
            await examplePrintStruct({
              name: prompt("Name to print")!,
            });
          }}
        >
          Print struct
        </button>

        <button
          onClick={async () => {
            try {
              await exampleReturnError();
            } catch (e) {
              alert(JSON.stringify(e));
            }
          }}
        >
          Return error
        </button>

        <button
          onClick={async () => {
            const res = await exampleReturnString();

            alert(JSON.stringify(res));
          }}
        >
          Return string
        </button>

        <button
          onClick={async () => {
            const res = await exampleReturnStruct();

            alert(JSON.stringify(res));
          }}
        >
          Return struct
        </button>

        <button
          onClick={async () => {
            try {
              await exampleReturnStringAndError();
            } catch (e) {
              alert(JSON.stringify(e));
            }
          }}
        >
          Return string and error
        </button>

        <button
          onClick={async () => {
            const res = await exampleReturnStringAndNil();

            alert(JSON.stringify(res));
          }}
        >
          Return string and nil
        </button>

        <button
          onClick={async () => {
            if (Notification.permission !== "granted") {
              await Notification.requestPermission();
            }

            while (true) {
              const res = await exampleNotification();

              if (res === "") {
                break;
              }

              new Notification(res);
            }
          }}
        >
          Get three notifications
        </button>
      </div>
    </main>
  );
};
