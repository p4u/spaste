# spaste

Simple Paste Bin like HTTP service.

### API

The following HTTP end point methods are available


`/add/<paste_name>` adds a new paste identified by paste\_name. The POST "data" value is saved as content.

`/add` adds a new paste and returns the automatic identifcation name (unix timestamp). The POST "data" value is saved as content.

`/list` list all pastes

`/get/<paste_name>` returns the content of paste identified by paste\_name

### Example using curl

Start server at port 8888

```
./spaste 8888
2019/01/27 22:31:45 Starting server ar port 8888
2019/01/27 22:37:21 Received [add new_paste] from 127.0.0.1:59754
2019/01/27 22:37:21 Adding 20 bytes of data to new_paste
```

Add a new paste

```
curl 127.0.0.1:8888/add/new_paste -d "data=This is my new paste"
key=new_paste
```

Get the paste content

```
curl 127.0.0.1:8888/get/new_paste
This is my new paste
```

### Example using spaste.sh

`spaste.sh` is a helper script which works on OpenWRT and other embedded devices (with busybox).

Usage: `pipe_input | sh spaste.sh [id_name]`

Default parameters execution.

```
echo "HELLO WORLD!" | sh spaste.sh hello
```

Or sepcify server and port by setting variable content

```
echo "HELLO WORLD!" | SERVER=127.0.0.1 PORT=8888 sh spaste.sh hello
```

