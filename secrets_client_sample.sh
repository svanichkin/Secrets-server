#!/bin/sh

#./secret_client_sample.sh password "Get me password for LUKS" "luks"

# Type: "confirm", "password", "text"
# Message: "Get me password for LUKS"
# Device: "Script get name this device from system" 
# On Server, device name used to automatically search for a password 
# in a folder named 'code' and retrieve the password if a match is found.
# Code: "luks"
# For example, the client may request a password for a specific application,
# such as "luks" (e.g., when a password is needed after a device reboot and disk reattachment).
# Or, for example, request confirmation for a USB flash drive inserted into the device;
# you can specify "usb" as an additional hint.

SERVERS="
https://localhost:8443
https://127.0.0.1:8443
"

OS=""
PKG_MANAGER=""
DEVICE_NAME=""

UNAME=$(uname -s)

if [ "$UNAME" = "Darwin" ]; then
  OS="macOS"
  PKG_MANAGER="brew"
  DEVICE_NAME=$(scutil --get ComputerName)
elif [ -f /etc/openwrt_release ]; then
  OS="OpenWRT"
  PKG_MANAGER="opkg"
  DEVICE_NAME=$(cat /proc/sys/kernel/hostname)
elif [ -f /etc/os-release ]; then
  . /etc/os-release
  OS="$ID"
  DEVICE_NAME=$(hostname)
  case "$ID" in
    ubuntu|debian) PKG_MANAGER="apt-get" ;;
    fedora|centos|rhel) PKG_MANAGER="yum" ;;
    arch) PKG_MANAGER="pacman" ;;
    alpine) PKG_MANAGER="apk" ;;
    *) PKG_MANAGER="" ;;
  esac
else
  OS="$UNAME"
  DEVICE_NAME=$(hostname)
fi

echo "Detected OS: $OS"
echo "Package manager: $PKG_MANAGER"
echo "Device name: $DEVICE_NAME"

check_and_install() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "$1 not found."
    if [ -n "$PKG_MANAGER" ]; then
      echo "Trying to install $1 using $PKG_MANAGER..."
      case "$PKG_MANAGER" in
        apt-get) sudo apt-get update && sudo apt-get install -y "$1" ;;
        yum) sudo yum install -y "$1" ;;
        pacman) sudo pacman -Sy --noconfirm "$1" ;;
        apk) sudo apk add "$1" ;;
        opkg) opkg update && opkg install "$1" ;;
        brew) brew install "$1" ;;
        *) echo "Unknown package manager, install $1 manually." ;;
      esac
    else
      echo "No package manager detected. Please install $1 manually."
    fi
  fi
}

check_and_install curl

REQ_TYPE=${1:-confirm}
REQ_MESSAGE=${2:-"Example operation requires your approval"}
REQ_CODE=${3:-usb}

echo "Request type: $REQ_TYPE"
echo "Request message: $REQ_MESSAGE"
echo "Request code: $REQ_CODE"

[ -f "$STOPFILE" ] && rm -f "$STOPFILE"

PIDS=""

log() {
  echo "[$(date '+%H:%M:%S')] $*"
}

send_request() {
  SERVER=$1
  
  [ -f "$STOPFILE" ] && {
    log "[$SERVER] Stopped before request (stopfile exists)"
    exit 0
  }
  
  log "[$SERVER] Starting request"
  
  RESPONSE=$(curl -k -s --connect-timeout 5 -X POST "$SERVER" \
    -H "Content-Type: application/json" \
    -d "{\"type\":\"$REQ_TYPE\",\"message\":\"$REQ_MESSAGE\",\"device\":\"$DEVICE_NAME\",\"code\":\"$REQ_CODE\"}")
  
  CURL_STATUS=$?
  
  if [ $CURL_STATUS -eq 0 ] && [ -n "$RESPONSE" ]; then
    log "[$SERVER] Got response: $RESPONSE"
    touch "$STOPFILE"
  else
    log "[$SERVER] No valid response or curl error (code $CURL_STATUS)"
  fi
}

log "Starting requests to servers..."

for SERVER in $SERVERS; do
  (
    [ -f "$STOPFILE" ] && exit 0
    
    echo "Trying $SERVER..."
    RESPONSE=$(curl -k -s --connect-timeout 5 -X POST "$SERVER" \
      -H "Content-Type: application/json" \
      -d "{\"type\":\"$REQ_TYPE\",\"message\":\"$REQ_MESSAGE\",\"device\":\"$DEVICE_NAME\",\"code\":\"$REQ_CODE\"}")
    
    if [ $? -eq 0 ] && [ -n "$RESPONSE" ]; then
      echo "Response from $SERVER: $RESPONSE"
      touch "$STOPFILE"
    fi
  ) &
  PIDS="$PIDS $!"
done

while :; do
  if [ -f "$STOPFILE" ]; then
    echo "Request succeeded. Killing others and exiting."
    rm -f "$STOPFILE"
    kill $PIDS 2>/dev/null
    exit 0
  fi
  
  STILL_RUNNING=0
  for PID in $PIDS; do
    if kill -0 $PID 2>/dev/null; then
      STILL_RUNNING=1
      break
    fi
  done
  
  if [ "$STILL_RUNNING" -eq 0 ]; then
    echo "All requests failed."
    exit 1
  fi
  
  sleep 1
done