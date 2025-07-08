#!/bin/sh

# Массив серверов
SERVERS=(
  "https://127.0.0.1:8443"
)

# Определение ОС, менеджера пакетов и имени устройства
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

# Проверка curl и jq
check_and_install() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "$1 not found."
    if [ -n "$PKG_MANAGER" ]; then
      echo "Installing $1 using $PKG_MANAGER..."
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
      echo "No package manager detected. Install $1 manually."
    fi
  fi
}

check_and_install curl
check_and_install jq

# Параметры запроса
REQ_TYPE="confirm"
REQ_MESSAGE="Example operation requires your approval"
REQ_CODE="usb"

# Поиск доступного сервера
for SERVER in "${SERVERS[@]}"; do
  echo "Testing $SERVER..."
  if curl -k -s --connect-timeout 2 "$SERVER" >/dev/null; then
    echo "Found reachable server: $SERVER"
    RESPONSE=$(curl -k -s -X POST "$SERVER" \
      -H "Content-Type: application/json" \
      -d "{\"type\":\"$REQ_TYPE\",\"message\":\"$REQ_MESSAGE\",\"device\":\"$DEVICE_NAME\",\"code\":\"$REQ_CODE\"}")
    echo "Response from server: $RESPONSE"
    exit 0
  fi
done

echo "No reachable servers found."
exit 1