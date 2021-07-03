#/bin/sh

sonar-scanner.bat -D"sonar.projectKey=grab-gitlab-group" -D"sonar.sources=." -D"sonar.host.url=http://192.168.1.6:9000" -D"sonar.login=8e7d427242f9f23ea0d915f7b9da9bf929e1b38c"
