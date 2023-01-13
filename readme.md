# Compiler API

![Logo Image](https://github.com/Ankan002/compiler-api/blob/main/assets/images/readme-logo.png)

This project is a standalone api that can be used to compile code on web on the fly.

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)
![GitHub](https://img.shields.io/badge/github-%23121011.svg?style=for-the-badge&logo=github&logoColor=white)
![DigitalOcean](https://img.shields.io/badge/DigitalOcean-%230167ff.svg?style=for-the-badge&logo=digitalOcean&logoColor=white)

## Installation Guide 🛠

### Via Docker

- Run the following command

```bash
docker run -d -p 8000:8000 -e PORT={PORT} ankan02/compile-me-api:0.0.2
```

- Replace `PORT` with the port of your choice.
- Also replace `8000:8000` with the `CONTAINER_PORT:PORT`.

### Via Github

- First install the following compilers:

    - MinGW (C/C++)
    - JDK (Java)
    - NodeJS v18 (Javascript)
    - Python v3.10 (Python)
    - TSNode (Typescript)
    - Golang
    - Rust
    - Kotlin
    - Mono MCS (C#)
  
- Now to clone the GitHub Repository run the following command:

```bash
git clone https://github.com/Ankan002/compile-me-api.git
```

- The go into the directory using the following command:

```bash
cd compile-me-api
```

- Now install the all the required dependencies using the following command:

```bash
go mod tidy
```

- Then create two files for the environment variables
  
  - .env
  - .env.production

- Now fill the `.env` file with the following variables:

| Variable | Description                                                             | Value                                                                                                      |
|----------|-------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------|
| PORT     | Port on which app will <br> be running.                                 | Valid `PORT` Number                                                                                        |
| GO_ENV   | Used define the environment <br> in which you want <br> to run the code | `development` (You might <br> use `production` but that <br> will not load the <br> environment variables) |

- Also, if you want to run the software using docker, then fill the `.env.production` with the following variables:

| Variable | Description                             | Value               |
|----------|-----------------------------------------|---------------------|
| PORT     | Port on which app will <br> be running. | Valid `PORT` Number |

- Now to run the project simply use the following command:

```bash
go run main.go
```

- Optionally, you might prefer to run it via Docker, then simply run:

```bash
docker compose up
```

- Congratulations... you have a running app. 🎉🎉... Enjoy...

### Contributors

- [Ankan Bhattacharya](https://github.com/Ankan002)

- [Dilpreet Grover](https://github.com/dfordp)

### Support My Work

[![BuyMeACoffee](https://img.shields.io/badge/Buy%20Me%20a%20Coffee-ffdd00?style=for-the-badge&logo=buy-me-a-coffee&logoColor=black)](https://buymeacoffee.com/ankan002)
