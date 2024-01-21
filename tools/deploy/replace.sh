#if [[ $EUID -ne 0 ]]; then
#   echo "[ERROR] This script must be run as root"
#   exit 1
#fi

if [ $# -eq 0 ]; then
    echo "[ERROR] Image version not specified"
    exit 1
fi

$APP_NAME="bwapi"

docker stop $APP_NAME && docker rm $APP_NAME
docker run -d --restart=unless-stopped --env-file ./.env -p 8083:8000 --name $APP_NAME loupeznik/better-wapi:$1
