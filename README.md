# AdGuard VPN Tray Indicator

[English](./README.md) | [Русский](./README_ru.md)
***

> ⚠️ **Warning!**
>
> This project was vibe coded aka written by an AI. I don't know Go. The code was
> created to solve a specific problem and might not follow all best
> practices.

A simple tray indicator app for Linux that lets you manage
`adguardvpn-cli`. It provides a basic GUI to connect, disconnect, and
view the status, since there is no official GUI for Linux.

## Features

*   Displays an icon in the system tray.
*   Shows the current connection status and location.
*   A "Connect" button if the VPN is disconnected, and a "Disconnect"
    button if it's connected.

It connects to the fastest available location, equivalent to running
`adguardvpn-cli connect -f`.

## Requirements

1.  **AdGuard VPN CLI**: The `adguardvpn-cli` utility must be
    installed, and you must be logged into your account.
2.  **Go**: Version 1.18 or newer.
3.  **Terminal**: `x-terminal-emulator` is used to enter the sudo
    password when connecting.
4.  **System Libraries** (for Debian/Ubuntu):
    ```bash
    sudo apt-get install build-essential libgtk-3-dev libappindicator3-dev
    ```
5.  You are expected to have already logged into your AdGuard account
    using `adguardvpn-cli login`.

## Building

1.  Clone the repository:
    ```bash
    git clone https://github.com/FroggerHH/adguard-indicator.git
    cd adguard-indicator
    ```

2.  Build the binary:
    ```bash
    go build
    ```

## Usage

1.  Run the compiled application:
    ```bash
    ./adguard-vpn-indicator-indicator
    ```

2.  For debugging, you can run it with the `-v` flag for verbose
    logging in the console:
    ```bash
    ./adguard-indicator -v
    ```

After launching, an icon will appear in the system tray. Right-click on
it to see the status and control menu.
