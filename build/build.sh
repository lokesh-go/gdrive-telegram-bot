# !/bin/bash

# remove old build image
docker rmi -f gdrive-telegram-bot-app:latest

# remove existing docker local file
rm -rf gdrive-telegram-bot-app-latest.tar

# build new app image
docker build --tag "gdrive-telegram-bot-app:latest" --force-rm=true --no-cache=true --file ../docker/Dockerfile ../

# save docker file locally
docker save gdrive-telegram-bot-app:latest -o gdrive-telegram-bot-app-latest.tar
