# DashD
> Lightweight daemon for Raspberry Pi driven kiosks

## Features
- software control for screen buttons (power, source, menu, plus, minus)
- lighting integration with Artemis-RGB (UDP and HTTP support)
- website display using Chromium based browser
- control via HTTP server

### Display
On startup DashD will create a browser window with the specified position and size, which can be changed later over HTTP. If the window is closed, DashD will exit.

|     CLI Flag      |           Default           |     Unit      |
| :---------------: | :-------------------------: | :-----------: |
| `display_enabled` |            true             |     Bool      |
|  `browser_path`   | "/usr/bin/chromium-browser" | String (Path) |
|   `default_url`   |                             | String (URL)  |
|    `window_x`     |              0              | UInt (Pixel)  |
|    `window_y`     |              0              | UInt (Pixel)  |
|  `window_width`   |            1920             | UInt (Pixel)  |
|  `window_height`  |            1080             | UInt (Pixel)  |

### Lighting
DashD can drive a single addressable LED strip with WS281X LEDs (any GRB LEDs should work). On most Raspberry Pi GPIO 12, 18, 40, and 52 can be used for connecting the data line of the LED strip. However if you are using Model B+, 2B or 3B only GPIO 12 and 18 are supported due to a hardware limitation.

|      CLI Flag      | Default |     Unit     |
| :----------------: | :-----: | :----------: |
| `lighting_enabled` |  true   |     Bool     |
|   `ledstrip_pin`   |   18    | UInt (GPIO)  |
| `ledstrip_length`  |   100   | UInt (Count) |

### Screen
DashD virtualizes the buttons of the screen, so they can be remote controlled via the HTTP server. The GPIOs should be connected to transistors in a way that a high output corresponds to pressing down the button.

|     CLI Flag     | Default |    Unit     |
| :--------------: | :-----: | :---------: |
| `screen_enabled` |  true   |    Bool     |
|   `power_pin`    |   17    | UInt (GPIO) |
|   `source_pin`   |   24    | UInt (GPIO) |
|    `menu_pin`    |   27    | UInt (GPIO) |
|    `plus_pin`    |   22    | UInt (GPIO) |
|   `minus_pin`    |   23    | UInt (GPIO) |

### Server
On startup DashD will start a HTTP server on the specified port, which serves a basic web interface. It allows the user to change the website shown, move and resize the window, configure the lighting and press the virtual screen buttons, if those modules are enabled.

|     CLI Flag     | Default |    Unit     |
| :--------------: | :-----: | :---------: |
| `server_enabled` |  true   |    Bool     |
|  `server_port`   |   80    | UInt (Port) |

### Socket
DashD will listen for UDP packets on the specified port, which can be used to control the lighting with [Artemis-RGB](https://github.com/Artemis-RGB/Artemis) or any compatible program that supports the Artemis Lighting Protocol. If lighting is disabled, this module will be disabled.

|     CLI Flag     | Default |    Unit     |
| :--------------: | :-----: | :---------: |
| `socket_enabled` |  true   |    Bool     |
|  `socket_port`   |  1872   | UInt (Port) |


## Installation
Download the latest release from [GitHub Releases](https://github.com/EliasStar/DashD/releases/latest) for your platform. Then run the following commands in the directory where you downloaded DashD:
``` sh
sudo apt-get update
sudo apt-get install chromium-browser
sudo mv DashD.* /usr/local/sbin/dashd
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

xhost +

/usr/local/sbin/dashd &> /var/log/dashd.log
```

Make a new systemd unit for xinit at `/etc/systemd/system/xinit.service` with this content:
``` systemd-unit
[Unit]
Description=Xinit display manager using startx

[Service]
Type=simple
ExecStart=/usr/bin/startx -- -nocursor
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
To build the docker image run `docker build -t dashd_builder:arm32 -f Dockerfile.arm32 .` in the root directory of this repository. Then run `docker run --rm --volume $(pwd):/dashd/app dashd_builder:arm32` to build DashD for 32-bit. The binary is located at `build/DashD.arm32`. Replace `arm32` with `arm64` in the previous commands to build for 64-bit.


## HTTP API
The endpoints ignore the HTTP method used, however POST is recommended for all, except `/` and `/config`, which should be used with GET.

|   Endpoint | Description                             |          Parameters | Parameter Description                                                                                          |
| ---------: | :-------------------------------------- | ------------------: | :------------------------------------------------------------------------------------------------------------- |
|        `/` | Get the web interface                   |                     |                                                                                                                |
| `/display` | Set the website to be displayed         |               `url` | the URL of the website as a percent-encoded string<br>data URLs are also supported                             |
|    `/move` | Set the window position                 |    `posX`<br>`posY` | the x position in pixel<br>the y position in pixel                                                             |
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
DashD - Lightweight daemon for Raspberry Pi driven kiosks<br>
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
