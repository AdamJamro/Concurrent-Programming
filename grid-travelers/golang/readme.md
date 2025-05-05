## <b>Grid Travelers</b>
### <b>GOLANG implementation</b>
Go Version 1.24.0

to compile a project navigate to the folder with desired <b>main.go </b> file

- <em>Compile the project:</em>
```bash
go build ./main.go
```


- <em>Run the project:</em>

```bash
go run ./main.go
```

- <em>Download dependencies (if needed)</em>:
```bash
    go mod download
```


## Visualisation
This project obtained two visualisation methods

1. [React client with websocket IPC](/grid-travelers/visualisation/react_visaulisation)
2. [Bash script - animation in terminal](/grid-travelers/visualisation/display-travel.bash)

### Usage

```bash
    [run simulation] > tmp_out.txt
    display-travel.bash tmp_out.txt
```