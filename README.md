# Whapp-Deltachat
Whapp-Deltachat is a bridge between WhatsApp and [Deltachat](https://delta.chat). It uses
[go-whatsapp](https://github.com/Rhymen/go-whatsapp) and
[go-deltachat](https://github.com/hugot/go-deltachat) respectively to map verified
deltachat groups to whatsapp conversations. It is currently not full-featured and most
likely is ridden with bugs, but pretty usable. Use at your own risk 🤓.

## Requirements
- __A system that go-deltachat will run on.__ It should work on all systems that are supported
  by golang and libdeltachat, but the static object file that is distributed with
  go-deltachat has been compiled on an amd64 linux system, so for other architectures you
  will need to compile the lib yourself.
- __An Android or iOS device that is capable of running the WhatsApp mobile app.__
  go-whatsapp uses the WhatsApp Web API, so when you use this bridge you will still need
  to have a device that is running WhatsApp somewhere.

## Install
There are no stable releases as of yet, so the only installation method right now is just
installing from source.

Installation example:
```bash
# clone
git clone git@github.com:hugot/whapp-deltachat.git

# build
cd ./whapp-deltachat
go build .

cp ./config.yml.example ./config.yml

# Now edit the config.yml file
"$EDITOR" ./config.yml

# Finally, run the program.
./whapp-deltachat ./config.yml

```

## First run
During the first time the program runs, there are two QR codes that need to be scanned:
First, the program will print a QR code to stdout. Scan this QR code from your deltachat
client to become a verified contact of the bridge. The bridge will then create a verified
chat with you and send you a QR code to scan from your WhatsApp mobile app. When you have
done this, the bridge will start doing it's thing.

### Known bug during first run
The first time that the bridge runs it will receive a lot of messages at the same time and
not in chronological order. This means that pre-existing WhatsApp messages will not be in
chronological order when they appear in your deltachat client. New messages received while
the bridge is running __will__ of course arrive in chronological order.

## Feature checklist
This is a list of features that are either implemented already or will be implemented in
the future.

- [x] Receive text messages
- [x] Receive image messages
- [x] Receive video messages
- [x] Receive audio messages
- [x] Receive document messages
- [ ] Receive contact (vcf) messages
- [x] Send text messages
- [ ] Send image messages
- [ ] Send video messages
- [ ] Send audio messages
- [ ] Send document messages
- [ ] Send contact (vcf) messages
- [ ] See chat's profile pic
- [ ] Sync user's deltachat profile pic with whatsapp profile pic
- [ ] Change group chat name


## Legal
This code is in no way affiliated with, authorized, maintained, sponsored or endorsed by
WhatsApp or any of its affiliates or subsidiaries. This is an independent and unofficial
software. Use at your own risk.
