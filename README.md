[![Contributors][contributors-shield]][contributors-url]
[![Language][language-shield]][language-url]
[![Issues][issues-shield]][issues-url]
[![GPL-3.0 License][license-shield]][license-url]
[![Watchers][watchers-shield]][watchers-url]

<!-- PROJECT LOGO -->
<br />
<p align="center">

  <h3 align="center">hanacall-timer</h3>

  <p align="center">
    Logs the time it takes to complete HANA_CALL's to HANA database
    <br />
    <a href="https://github.com/rfparedes/hanacall-timer/issues">Report Bug</a>
    ·
    <a href="https://github.com/rfparedes/hanacall-timer/issues">Request Feature</a>
  </p>
</p>

<!-- TABLE OF CONTENTS -->
<details open="open">
  <summary><h2 style="display: inline-block">Table of Contents</h2></summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#technical-details">Technical Details</a></li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#build-it-yourself">Build It Yourself</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
  </ol>
</details>

<!-- ABOUT THE PROJECT -->
## About The Project

hanacall-timer was developed so that users can see how long HANA_CALL's from the resource agent to HANA take to return. This can be used to confirm the HANA_CALL's are returning in a timely manner and to have data showing the length of time (in ms) for further troubleshooting.
<!-- GETTING STARTED -->
## Getting Started

To get a local copy up and running follow these simple steps.

### Prerequisites

* SAP HANA High Availability Cluster

### Installation

Download the binary from Releases (<https://github.com/rfparedes/hanacall-timer/releases/latest/download/hanacall-timer>) to `/usr/local/bin` on the server and run:

```sh
sudo chmod +x /usr/local/bin/hanacall-timer
```

Start it

```sh
sudo hanacall-timer start --sidadm <SIDADM>
```

Check Status Anytime

```sh
hanacall-timer status
```

## Technical Details

* hanacall-timer will time two HANA interface calls. These calls are made exactly as they are from within the SAPHanaSR resource agents.
   1. systemReplicationStatus.py
   2. landscapeHostConfiguration.py
* When started, a systemd service and timer are created.  The timer is enabled and will run hanacall-timer every 60 seconds
* Only one instance of hanacall-timer will run at any given time. If HANA_CALL takes longer than 60 seconds, no additional calls will be made until the HANA_CALL's return
* When starting, the HANA <SIDADM> user needs to be specified as this is the user making the HANA_CALL
* The log includes command output, command return codes and timings logged to /var/log/hanacall-timer.log
* The output of the timings are also logged in csv format to /var/log/hanacall-timer.csv in the format:
  
  `RFC3339 date-time, systemReplicationStatus.py time (ms), landscapeHostConfiguration.py time (ms)`

As an example, you can see the HANA_CALL landscapeHostConfiguration timings increasing on this server from 293ms to a high of 32.097s:
```sh
hana11:~ # tail -f /var/log/hanacall-timer.csv
2021-03-27T19:57:45-04:00,297,293
2021-03-27T19:58:46-04:00,317,316
2021-03-27T19:59:46-04:00,373,373
2021-03-27T20:00:46-04:00,594,594
2021-03-27T20:01:47-04:00,421,420
2021-03-27T20:02:47-04:00,2375,2346
2021-03-27T20:03:48-04:00,3670,3639
2021-03-27T20:04:49-04:00,3314,3297
2021-03-27T20:05:49-04:00,32538,32097
2021-03-27T20:06:50-04:00,31271,31055
```

## Usage

### To start, run

```sh
sudo hanacall-timer start --sidadm <SIDADM>
```

### To stop, run

```sh
sudo hanacall-timer stop
```

### To see the current start/stop status of hanacall-timer, run

```sh
hanacall-timer status
```

## Build it yourself

* You'll need a go compiler installed

Clone it

```sh
git clone https://github.com/rfparedes/hanacall-timer.git
```

Build it

```sh
cd hanacall-timer
go build -o hanacall-timer
```

Move it

```sh
mv hanacall-timer /usr/local/bin
sudo chmod +x /usr/local/bin/hanacall-timer
```

Start it

```sh
sudo hanacall start --sidadm <SIDADM>
```

<!-- CONTRIBUTING -->
## Contributing

Contributions are what make the open source community such an amazing place to be learn, inspire, and create. Any contributions you make are **greatly appreciated**.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<!-- LICENSE -->
## License

Distributed under the GPL-3.0 License. See `LICENSE` for more information.

<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[contributors-shield]: https://img.shields.io/github/contributors/rfparedes/hanacall-timer?color=%20%2330BA78
[contributors-url]: https://github.com/rfparedes/hanacall-timer/graphs/contributors
[language-shield]: https://img.shields.io/github/languages/top/rfparedes/hanacall-timer?color=%20%2330BA78
[language-url]: https://github.com/rfparedes/hanacall-timer/search?l=go
[watchers-shield]: https://img.shields.io/github/watchers/rfparedes/hanacall-timer?color=%20%2330BA78&style=social
[watchers-url]:https://github.com/rfparedes/hanacall-timer/watchers
[issues-shield]: https://img.shields.io/github/issues/rfparedes/hanacall-timer?color=%20%2330BA78
[issues-url]: https://github.com/rfparedes/hanacall-timer/issues
[license-shield]: https://img.shields.io/github/license/rfparedes/hanacall-timer?color=%20%2330BA78
[license-url]: https://github.com/rfparedes/hanacall-timer/blob/main/LICENSE
