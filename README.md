<p align="center"><img src="https://user-images.githubusercontent.com/1092882/60883564-20142380-a268-11e9-988a-d98fb639adc6.png" alt="webgo gopher" width="256px"/></p>

# Gdrive Telegram Bot

Gdrive telegram bot implemented the file searching feature on Google drive.

---

**Commands** I have added in the bot:
  
  - **/srch** search_file_name - Searches the file from your google drive.
  - **/info** - Information about the server

    - CPU
    - RAM used by program
    - No of Go routines are created by program
    - Server internet download speed
    - Server internet upload speed
    - Server internet latency


---

Process to Run this project

  - Adds these into environment variable
    - Adds path of Google credentials.json file.
    - Adds Google token path where token will be created
    - Adds telegram bot token
    - Adds upload folder path


- Build the docker image ( build.sh scripts in build directory )
- Run the project ( run.sh scripts in the run directory )