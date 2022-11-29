<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta http-equiv="X-UA-Compatible" content="IE=edge" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>{{ .AppName }}</title>
  </head>
  <body>
    <h1>{{ .AppName }}</h1>

    <h2>Create a new task</h2>

    <form action="/create" method="post">
      <label for="title">Title</label>
      <br />
      <input type="text" name="title" id="title" required />
      <br />

      <label for="body">Body</label>
      <br />
      <textarea type="text" name="body" id="body" required></textarea>
      <br />

      <input type="submit" value="Create task" />
    </form>

    <h2>All tasks</h2>

    {{"{{"}} if gt (len .Todos) 0 {{"}}"}}
    <ul>
      {{"{{"}} range $id, $todo := .Todos {{"}}"}}
      <li>
        <h3>{{"{{"}} $todo.Title {{"}}"}}</h3>
        <p>{{"{{"}} $todo.Body {{"}}"}}</p>

        <form action="/delete?id={{"{{"}} $id {{"}}"}}" method="post">
          <input type="submit" value="Delete task" />
        </form>
      </li>
      {{"{{"}}
        end
      {{"}}"}}
    </ul>
    {{"{{"}} else {{"}}"}}
    <span>No tasks yet.</span>
    {{"{{"}} end {{"}}"}}

    <footer>
      <p>
        <code>
          Rendered by Go {{"{{"}} .GoVersion {{"}}"}} {{"{{"}} .GoOS {{"}}"}}/{{"{{"}} .GoArch {{"}}"}} at
          {{"{{"}}
          .RenderTime {{"}}"}}
        </code>
      </p>
    </footer>
  </body>
</html>
