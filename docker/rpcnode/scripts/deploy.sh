#!/usr/bin/env sh

SKIP_BUILD=${SKIP_BUILD:-""}

# Set up env
export GOPATH=$HOME/go
export GOBIN=$GOPATH/bin
export BUILD_PATH=/tabilabs/tabi-v2/build
export PATH=$GOBIN:$PATH:/usr/local/go/bin:$BUILD_PATH
echo "export GOPATH=$HOME/go" >> /root/.bashrc
echo "GOBIN=$GOPATH/bin" >> /root/.bashrc
echo "export PATH=$GOBIN:$PATH:/usr/local/go/bin:$BUILD_PATH" >> /root/.bashrc
/bin/bash -c "source /root/.bashrc"
mkdir -p $GOBIN

# Step 1 build tabid
if [ -z "$SKIP_BUILD" ]
then
  /usr/bin/build.sh
fi
cp build/tabid "$GOBIN"/

# Run init to set up state sync configurations
/usr/bin/configure_init.sh

# Start the chain
/usr/bin/start_tabi.sh

tail -f /dev/null
