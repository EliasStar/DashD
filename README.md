# DashD
> Lightweight daemon for Raspberry Pi driven kiosks

## Features
  - software control for screen buttons (power, source, menu, plus, minus)
  - lighting integration with Artemis-RGB (UDP and HTTP support)
  - website display using WebKitGTK
  - control via HTTP server

### Display
On startup DashD will create a single WebKitGTK window with the specified size, which can be resized later. If the window is closed, DashD will open a new window, however any previously shown website will be lost.

|     CLI Flag     | Default | Unit  |
| :--------------: | :-----: | :---: |
| `display_width`  |  1024   | Pixel |
| `display_height` |   768   | Pixel |

### Lighting
DashD can drive a single addressable LED strip with WS281X LEDs (any GRB LEDs should work). On most Raspberry Pi GPIO 12, 18, 40, and 52 can be used for connecting the data line of the LED strip. However if you are using Model B+, 2B or 3B only GPIO 12 and 18 are supported due to a hardware limitation.

|     CLI Flag      | Default | Unit  |
| :---------------: | :-----: | :---: |
|  `ledstrip_pin`   |   18    | GPIO  |
| `ledstrip_length` |   62    | Count |

### Screen
DashD virtualizes the buttons of the screen, so they can be remote controlled via the HTTP server. The GPIOs should be connected to transistors in a way that a high output corresponds to pressing down the button.

|   CLI Flag   | Default | Unit  |
| :----------: | :-----: | :---: |
| `power_pin`  |   17    | GPIO  |
| `source_pin` |   24    | GPIO  |
|  `menu_pin`  |   27    | GPIO  |
|  `plus_pin`  |   22    | GPIO  |
| `minus_pin`  |   23    | GPIO  |

### Server
On startup DashD will start a HTTP server on the specified port, which serves a basic web interface. It allows the user to change the website shown, resize the window, configure the lighting and press the virtual screen buttons. DashD will also listen for UDP packets on the specified port, which can be used to control the lighting with [Artemis-RGB](https://github.com/Artemis-RGB/Artemis) or any compatible program that supports the Artemis Lighting Protocol.

|  CLI Flag   | Default | Unit  |
| :---------: | :-----: | :---: |
| `http_port` |   80    | Port  |
| `udp_port`  |  1872   | Port  |


## Installation
Download the latest release from [GitHub Releases](https://github.com/EliasStar/DashD/releases/latest). Then run the following commands in the directory where you downloaded DashD:
``` sh
sudo apt-get update
sudo apt-get install gtk-3.0 webkit2gtk-4.0
sudo mv DashD /usr/local/sbin/dashd
sudo chmod +x /usr/local/sbin/dashd
```


## Usage
One possible way to use DashD is to run it on startup with systemd and xinit as basic display manager. To use this configuration install xinit like so:
``` sh
sudo apt-get update
sudo apt-get install xserver-xorg xfonts-base xinit
```

Then put the following in `/etc/X11/xinit/xinitrc`:
``` sh
#!/bin/bash

/usr/local/sbin/dashd &> /var/log/DashD.log
```

Make a new systemd unit for xinit at `/etc/systemd/system/xinit.service` with this content:
``` systemd-unit
[Unit]
Description=Xinit display manager using startx

[Service]
Type=simple
ExecStart=/usr/bin/startx
Restart=always

[Install]
Alias=display-manager.service
```

Finally enable the service, switch to the `graphical` target and reboot. The logs will be put in `/var/log/DashD.log`.
``` sh
sudo systemctl enable xinit.service
sudo systemctl set-default graphical.target
sudo reboot
```

## Building
DashD utilizes Docker to cross-compile for the Raspberry Pi. The Dockerfile in this repository creates an image that can be used to build DashD.
To build the docker image run `docker build -t dashd_builder .` in the root directory of this repository. Then run `docker run --rm --volume $(pwd):/app dashd_builder` to build DashD. The binary is located at `build/DashD`.


## HTTP API
The endpoints ignore the HTTP method used, however POST is recommended for all, except `/` and `/config`, which should be used with GET.

|   Endpoint | Description                             |          Parameters | Parameter Description                                                                                          |
| ---------: | :-------------------------------------- | ------------------: | :------------------------------------------------------------------------------------------------------------- |
|        `/` | Get the web interface                   |                     |                                                                                                                |
| `/display` | Set the website to be displayed         |               `url` | the URL of the website as a percent-encoded string<br>data URLs are also supported                             |
|  `/resize` | Set the window size                     | `width`<br>`height` | the width in pixel<br>the height in pixel                                                                      |
|  `/config` | Get the LED count in JSON (Artemis API) |                     |                                                                                                                |
|  `/update` | Set colors of LED strip (Artemis API)   |            `base64` | a UDP packet (described below) as a base64 encoded string<br>name is not enforced; must be the first parameter |
|   `/reset` | Set all LEDs to black (Artemis API)     |                     |                                                                                                                |
|   `/power` | Toggle the power button                 |                     |                                                                                                                |
|  `/source` | Toggle the source button                |                     |                                                                                                                |
|    `/menu` | Toggle the menu button                  |                     |                                                                                                                |
|    `/plus` | Toggle the plus button                  |                     |                                                                                                                |
|   `/minus` | Toggle the minus button                 |                     |                                                                                                                |


## UDP Lighting Protocol
Each UDP packet contains a single frame of lighting data. The first byte is the sequence number which is used for basic packet ordering and must be incremented with each request. The second byte is the channel number which is ignored by DashD because it supports only a single channel. The remaining packet consists of byte triplets which represent the RGB value in that order. The number of byte triplets does not need to be equal to the number of LEDs. Any extra bytes are ignored and LEDs not contained in the packet retain their color. To calculate the length of a packet that sets all LEDs use: `2 + (leds * 3)`.


## License
DashD - Lightweight daemon for Raspberry Pi driven kiosks <br>
Copyright (C) 2022 Elias*

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.
