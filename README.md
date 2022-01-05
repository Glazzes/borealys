## Borealys

### About
Borealys is a small and simple remote code execution engine, the idea
behind this project is quite simple, I have been curious of websites like
leetcode and language specific playgrounds running code that has been 
written in a web browser, so I decided to try my luck and create my own.

This project runs in a `ubuntu:20.04` docker container based image with
the following changes:
- An extra group called runners
- A root level folder called binaries containing all the binaries for each
  supported language
  
- 100 unpriveleged users whose will the code through the `runuser` command

### Goals
I decided to build this project in GoLang as I really wanted to learn it
and this project required a fast and lightweight programming language, as
a Spring developer I'm aware of Java/Kotlin's memory consumption, so Spring
was not an option.

My knowledge of the Linux operating system is still quite limited out of
the average user needs, so this project brings me many of the
following challenges:
- Create new system users
- Get a better understanding of group and permissions
- Limit user resources, open processes, open files and more
- Learn the basics of Bash scripting

### Behaviour
- Users will send their code along with the selected language name and version,
  if available, user's code will be run against the selected language, if
  the code succeeds the stdout of the executed code will be sent in the response,
  if your code can not be executed lets say by a missing `;` you will get
  the stderr back

### Supported languages
Hopefully in their newest versions!
- Bash
- Go
- Java
- Javascript 
- Kotlin (work in progress)
- Python

### Api endpoints
- GET /api/languages => returns an array of all the available languages,
  example:
```json
[
  {
    "language": "java",
    "timeout": 4,
    "version": "17"
  },
  {
    "language": "python",
    "timeout": 3,
    "version": "3.10.1"
  }
]
```

- POST /api/run => Given a valid piece of code will give you back the stdout,
if the code can not be executed, you will receive the stderr instead

```json
{
  "language": "Javascript",
  "code": [
    "const words = ['hola', 'hi']",
    "\n",
    "for(let word of words){",
    "   console.log(`current word is ${word}`)",
    "}"
  ]
}
```