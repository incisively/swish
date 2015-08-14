#!/bin/bash
set -e

SWISH_URL=localhost:9001

SWISH=-1
TARGET_1=-1
TARGET_2=-1
TARGET_3=-1

echo "Setting up..."
socat TCP-LISTEN:12341,fork,reuseaddr SYSTEM:"echo HTTP/1.0 200; echo; echo Target 1" > /dev/null &
TARGET_1=$!
socat TCP-LISTEN:12342,fork,reuseaddr SYSTEM:"echo HTTP/1.0 200; echo; echo Target 2" > /dev/null &
TARGET_2=$!
socat TCP-LISTEN:12343,fork,reuseaddr SYSTEM:"echo HTTP/1.0 200; echo; echo Target 3" > /dev/null &
TARGET_3=$!

function finish {
  kill $SWISH
  kill $TARGET_1
  kill $TARGET_2
  kill $TARGET_3
}
trap finish EXIT


function assert_equal {
  echo "Expected '$1', got '$2'"
  if [ "$1" != "$2" ]; then
    echo "FAILED."
    exit 1
  else
    echo "PASSED!"
  fi
}

# start swish

./swish -bind=$SWISH_URL &
SWISH=$!
sleep 0.5

# tests

echo
echo "Test: Create listeners"
curl $SWISH_URL -d "listen=:12340&target=localhost:12341" --silent
curl $SWISH_URL -d "listen=:12349&target=localhost:12343" --silent
response=$(curl localhost:12340 --silent)
assert_equal 'Target 1' "$response"
response=$(curl localhost:12349 --silent)
assert_equal 'Target 3' "$response"

echo
echo "Test: Update listener"
curl $SWISH_URL -d "listen=:12340&target=localhost:12342" --silent
response=$(curl localhost:12340 --silent)
assert_equal 'Target 2' "$response"
response=$(curl localhost:12349 --silent)
assert_equal 'Target 3' "$response"

# echo
# echo "Test: Delete listener"
# curl "$SWISH_URL?port=12340" -XDELETE --silent
# response=$(curl localhost:12340 --silent)
# if [ $? -eq 0 ]; then
#   echo "FAILED. Expected request to deleted listener not to succeed."
#   exit 1
# fi

echo
echo "All passed!"
