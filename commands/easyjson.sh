#!/bin/sh
echo ""
echo "  --------------------------"
echo "  -----Easyjson generate----"
echo "  --------------------------"
echo ""
trap 'echo " stop" ' INT TERM
#chmod +x easyjson.sh && ./easyjson.sh

# set GOPATH and PATH
export GOPATH=$HOME/go
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
done=0

# install easyjson
#go get -u github.com/mailru/easyjson/...
#apt install golang-easyjson
#go mod tidy
#go get -u
#go mod vendor

echo "  1. Copy project to GOPATH"
# we need THISDIR to return back at the end
export THISDIR=$PWD
export GPROJECTDIR=$GOPATH/src/github.com/go-park-mail-ru/2019_1_Escapade
# create folder, -p - create parents folders
mkdir -p $GPROJECTDIR && \
# -r - copy folder
#cp -r $PWD/.. $GPROJECTDIR && \
cp -r $PWD/../internal $GPROJECTDIR && \
cp -r $PWD/../auth $GPROJECTDIR && \
cp -r $PWD/../chat_service $GPROJECTDIR && \
#cp -r $PWD/../go.mod $GPROJECTDIR && \
#cp -r $PWD/../go.sum $GPROJECTDIR && \
cp -r $PWD/../vendor $GPROJECTDIR && \

echo "  2.1 Apply easyjson to constants" && \
export CONSTANTSPATH=$GPROJECTDIR/internal/constants && \
cd $CONSTANTSPATH && \
easyjson . && \
cp $CONSTANTSPATH/constants_easyjson.go $THISDIR/../internal/constants && \

echo "  2.2 Apply easyjson to models" && \
export MODELSPATH=$GPROJECTDIR/internal/models && \
cd $MODELSPATH && \
easyjson . && \
cp $MODELSPATH/models_easyjson.go $THISDIR/../internal/models && \

echo "  2.3 Apply easyjson to config" && \
export CONFIGPATH=$GPROJECTDIR/internal/config && \
cd $CONFIGPATH && \
easyjson . && \
cp $CONFIGPATH/config_easyjson.go $THISDIR/../internal/config && \

echo "  2.4 Apply easyjson to photo" && \
export PHOTOPATH=$GPROJECTDIR/internal/photo && \
cd $PHOTOPATH && \
easyjson . && \
cp $PHOTOPATH/photo_easyjson.go $THISDIR/../internal/photo && \

echo "  2.5 Apply easyjson to game" && \
export GAMEPATH=$GPROJECTDIR/internal/game && \
cd $GAMEPATH && \
easyjson . && \
cp $GAMEPATH/game_easyjson.go $THISDIR/../internal/game && \
done=1

echo "  3. Remove project from GOPATH"
rm -R $GPROJECTDIR
cd $THISDIR

echo ""
if [ "$done" -eq 1 ]
then 
echo "  ----------Done!-----------"
else
echo "  ----------Error!-----------"
exit 1
fi