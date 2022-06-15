# DashD
> Lightweight daemon for Raspberry Pi driven kiosks

## Features
  - software control for screen buttons (power, source, menu, plus, minus)
  - lighting integration with Artemis-RGB (UDP and HTTP support)
  - website display using WebKitGTK
  - control via HTTP server

### Display
On startup DashD will create a single WebKitGTK window with the specified size, which can be resized later. If the window is closed, DashD will open a new window, however any previously shown website will be lost.

| Command Line Flag | Default Value | Unit  |
| :---------------: | :-----------: | :---: |
|  `display_width`  |     1024      | Pixel |
| `display_height`  |      768      | Pixel |

### Lighting
DashD can drive a single addressable LED strip with WS281X LEDs (any GRB LEDs should work). On most Raspberry Pi GPIO 12, 18, 40, and 52 can be used for connecting the data line of the LED strip. However if you are using Model B+, 2B or 3B only GPIO 12 and 18 are supported due to a hardware limitation.

| Command Line Flag | Default Value | Unit  |
| :---------------: | :-----------: | :---: |
|  `ledstrip_pin`   |      18       | GPIO  |
| `ledstrip_length` |      62       | Count |

### Screen
DashD virtualizes the buttons of the screen, so they can be remote controlled via the HTTP server. The GPIOs should be connected to transistors in a way that a high output corresponds to pressing down the button.

| Command Line Flag | Default Value | Unit  |
| :---------------: | :-----------: | :---: |
|    `power_pin`    |      17       | GPIO  |
|   `source_pin`    |      24       | GPIO  |
|    `menu_pin`     |      27       | GPIO  |
|    `plus_pin`     |      22       | GPIO  |
|    `minus_pin`    |      23       | GPIO  |

### Server
On startup DashD will start a HTTP server on the specified port, which serves a basic web interface. It allows the user to change the website shown, resize the window, configure the lighting and press the virtual screen buttons.

| Command Line Flag | Default Value | Unit  |
| :---------------: | :-----------: | :---: |
|    `http_port`    |      80       | Port  |
|    `udp_port`     |     1872      | Port  |

## Installation
Download the latest release from [GitHub Releases](https://github.com/EliasStar/DashD/releases/latest). Then run the following commands in the directory where you downloaded DashD:
``` shell
sudo apt-get update
sudo apt-get install gtk-3 webkit2gtk-4.0
sudo mv DashD /usr/local/sbin/dashd
sudo chmod +x /usr/local/sbin/dashd
```

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
