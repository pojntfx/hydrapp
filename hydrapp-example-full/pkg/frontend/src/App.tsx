const socket = new WebSocket("ws://localhost:1337");

const local = {
  ExampleNotification: async (msg: string) => {
    if ("Notification" in window && Notification.permission !== "granted") {
      await Notification.requestPermission();
    }

    if ("Notification" in window) {
      new Notification(msg);
    } else {
      alert(msg);
    }
  },
};
const remote = {
  ExamplePrintString: async (msg: string) => {},
  ExamplePrintStruct: async (input: any) => {},
  ExampleReturnError: async () => {},
  ExampleReturnString: async () => "",
  ExampleReturnStruct: async (): Promise<any> => {},
  ExampleReturnStringAndError: async () => "",
  ExampleNotification: async () => {},
};

for (let functionName in remote) {
  (remote as any)[functionName] = async (...args: any[]) => {
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

export default () => {
  return (
    <main>
      <h1>Hydrapp Full Example</h1>

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
              alert(JSON.stringify(e));
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
              alert(JSON.stringify(e));
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
