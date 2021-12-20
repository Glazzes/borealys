## Borealys docker image

### About
This image is based on `ubuntu:20.04` with the following changes:
- users
    - Adds 100 new un privileged users whose will run the code
    - Add limits to memory, max open files, number or processes per user
    
- Adds binaries for each one of the programming language named folders in
this directory