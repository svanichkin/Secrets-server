
# Secrets

**Secrets** is a simple server designed to run as a system service on trusted devices (e.g., laptops) to respond to external requests for confirmation or password input.

## Purpose

Secrets is intended for remote approval of actions or password input:

- Confirming the connection of external USB devices.
- Providing passwords for unlocking LUKS-encrypted disks.
- Other operations requiring manual approval from a trusted client.

## How It Works

1. A client device (e.g., Raspberry Pi, USB drive, server) sends a POST request to the Secrets server.
2. The server displays a dialog window asking for confirmation or a password.
3. The user approves or cancels the request.
4. The server returns a simple plain-text response or an HTTP error (in case of denial or cancellation).

## Protocol

- HTTP POST to `https://<IP>:<PORT>`
- JSON request example:
  ```json
  {
    "type": "confirm" | "password" | "text",
    "message": "Example message",
    "code": "optional_code",
    "device": "optional_device"
  }
  ```
  Device name used to automatically search for a password in a folder named 'code' and retrieve the password if a match is found.
  For example, the client may request a password for a specific application, such as "LUKS" (e.g., when a password is needed after a device reboot and disk reattachment). Or, for example, request confirmation for a USB flash drive inserted into the device; you can specify "usb" as an additional hint.

- Response:
  - Success: `1` (for confirm) or the entered text (for password/text).
  - Error: `400 Bad Request` if the user cancels or denies.

## Example Run

Secrets runs as an HTTPS server accepting requests only from trusted IPs (specified in the `trusted` file).

```bash
go run main.go
```
## Configuration Requirements

The configuration requires the following files and folders:

- **Passwords folder**  
  Stores cached passwords for automatic reuse and filling.

- **Server**  
  Contains the server address in the format `IP:PORT`.

- **Trusted file**  
  Recursively finds all **/filename (pattern from config) files with allowed IPs (one per line, space, comma, or newline separated).
  
## License

MIT
