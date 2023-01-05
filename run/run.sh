# !/bin/bash

# load docker image
docker load gdrive-telegram-bot-app-latest.tar

# run docker container
docker run -d --restart=always --env-file $HOME/.env --name gdrive-telegram-bot-app gdrive-telegram-bot-app-latest