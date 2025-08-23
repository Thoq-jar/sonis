# Sonis

Open-source lossless streaming audio server

> [!NOTE]
> Currently Sonis only works for `.flac` files.
> wav, mp3, and the popular formats will be supported in the future and are planned.

## Setup
> [!IMPORTANT]
> [Go](https://go.dev/) must be installed on your system.
> 
> [Docker]("https://docker.com") is recommended but not required if you build using go toolchain

- Begin by creating a `music` directory in the root of the project.

- With Docker (recommended):
    - Start Sonis using `docker-compose up`

- With Go toolchain
    - Run `go build .` to build the server.
    - Run `./sonis` to start the server.

- Open `http://localhost:8080` in your browser.
- Select a track and enjoy!

## License
> © 2025 [Thoq](https://thoq.dev)

This project uses the MIT license, see the
[LICENSE](LICENSE) for more details!
