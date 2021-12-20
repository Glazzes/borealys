## Borealys

### About
Borealys is a small and simple remote code execution engine, the idea
behind this project is quite simple, the user sends their code in one
of the available languages, and the api will run the given
code and return the output back to the user

### Goals
I decided to build this project in GoLang as I really wanted to learn it
and this project required a fast and lightweight programming language, as
a Spring developer I'm aware of Java/Kotlin's memory consumption, so Spring
was not an option.

My knowledge of the Linux operating system is still quite limited out of
the average user needs, so this project brings me many of the
following challenges:
- Create new system users
- Understanding of group and permissions
- Limit user resources, open processes, open files and more
- Learn the basics of Bash and Make

### Behaviour
- Users will send their code along with the selected language name and version,
  if available, user's code will be run against the selected language, either if the
  code succeeds or fails, the output of this process will be sent in the
  http response along a status code representing the final state of
  the operation.

### Api endpoints
- GET /api/languages => returns an array of all the available languages,
  example:
```json
[
  {
    "language": "java",
    "versions": ["17", "11", "8"]
  },
  {
    "language": "python",
    "versions": ["3.10.1"]
  }
]
```