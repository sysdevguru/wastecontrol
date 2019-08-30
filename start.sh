export PORT=5000
export GOPATH=/Users/stemann/Developer/wastecontrol/wastecontrol
export PATH=$PATH:/Users/stemann/Developer/wastecontrol/wastecontrol/bin
export ENVIRONMENT=STAGE
cd /Users/stemann/Developer/wastecontrol/wastecontrol/src/wastecontrol
go install && heroku local web
