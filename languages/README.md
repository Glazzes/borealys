## Borealys docker image

### About
This image is based on `ubuntu:20.04` with the following changes:
- users
    - Adds 100 new unprivileged users whose will run the code
    - Add limits to memory, open files and open processes per user
    
- Adds binaries for each one of the programming language named folders within
this directory