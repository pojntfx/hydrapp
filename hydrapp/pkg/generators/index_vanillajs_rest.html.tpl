<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta http-equiv="X-UA-Compatible" content="IE=edge" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <meta name="color-scheme" content="light dark">
    <title>{{ .AppName }}</title>
  </head>
  <body>
    <h1>{{ .AppName }}</h1>

    <ul>
      <li>
        <a id="server-time-link">Get server time</a>
      </li>

      <li>
        <a id="ifconfig-io-link"
          >Get results of ifconfig.io/all.json proxied through backend</a
        >
      </li>

      <li>
        <a id="envs-link">List all system environment variables</a>
      </li>
    </ul>

    <script>
      function getUrlParameter(location, name) {
        name = name.replace(/[\[]/, "\\[").replace(/[\]]/, "\\]");
        var regex = new RegExp("[\\?&]" + name + "=([^&#]*)");
        var results = regex.exec(location);
        return results === null
          ? ""
          : decodeURIComponent(results[1].replace(/\+/g, " "));
      }

      document.getElementById("server-time-link").href =
        getUrlParameter(window.location.href, "socketURL") + "/servertime";

      document.getElementById("ifconfig-io-link").href =
        getUrlParameter(window.location.href, "socketURL") + "/ifconfigio";

      document.getElementById("envs-link").href =
        getUrlParameter(window.location.href, "socketURL") + "/envs";
    </script>
  </body>
</html>
