# fileTun

After reading https://xeiaso.net/blog/anything-message-queue/ this entertaining post, i wanted to implement approximatly what they did in their setup.

essentially, connect two machines using a shared file.

Difference with my imlementation: actual file is needed and there is only default heartbeat.

## Build
```bash
go mod tidy
go build main.go
```

## Run

you need two linux machines with root access and two files. You want to have a `left` and a `right` machine
For the left machine, you need to read from a `right_output.gob` and you need to write to a `left_output.gob`
For the right machine, its the opposite.
Both machines need to have access to the same file.
You can use e.g. sshfs or ntfs or probably s3 or a gcp bucket

On the left machine:
```bash
mkdir testing
sudo ./main --own_cidr 10.0.9.0/24 --own_name left --input ./testing/input.gob --output ./testing/output.gob --peer_cidr='10.0.8.0/24'
```

On the right machine:
```bash
mkdir testing
sudo ./main --own_cidr 10.0.8.0/24 --own_name right --input ./testing/input.gob --output ./testing/output.gob --peer_cidr='10.0.9.0/24'
```

On both machines you should have a tun device, named `left` or `right` depending on the machine.

<!-- TODO: -->
You could now e.g. spawn a webserver on the right machine and have it listen to :8080 with something like
```bash
echo "it works" > index.html
python3 -m http.server 8080 --bind 10.0.8.1
```
and from the left machine, you should be abled to
```bash
curl 10.0.8.1:8080
```
You should see "it works" output as result of your curl command.
```bash
tail -f testing/input.gob testing/output.gob
```
You should see some content in both input.gob and output.gob
You should also be abled to ping or ssh anything on the remote subnet



## epic development roadmap

- [x] add tests
- [] add better tests
- [] clean things up
- [] make the cli useful and intuitive
- [x] add different file names - paths from cli
- [] move file reading / writing into seperate file, interface
- - [] add more backends ＼(^o^)／
- - - [] add pastebin backend
- - - [] add slack or whatsapp or facebook or something like that, instagram chat, tiktok comments
- [] add more settings, config yaml and env vars
- [] figure out how to do things as non-root
- [] use context
