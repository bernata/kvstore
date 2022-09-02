#!/bin/bash

set +eu

HAS_FAILED=0

./bin/kvservice -p 8282 2>/dev/null &
KVSERVICE_PID=$!
sleep 1 # allow time for server to start to listen on port

echo "TEST: key not found"
RESULT=$(./bin/kvclient get --key foo 2>&1)
if [[ "$RESULT" == "kvclient: error: [404]: [404]: key 'foo' not found" ]]; then
  echo "  PASS"
else
  HAS_FAILED=1
  echo "  FAIL: " $RESULT
fi

echo

echo "TEST: write key and read"
RESULT=$(./bin/kvclient write --key foo/bar --value baz 2>&1)
if [[ "$RESULT" != "{}" ]]; then
  HAS_FAILED=1
  echo "  FAIL: " $RESULT
fi
RESULT=$(./bin/kvclient get --key foo/bar 2>&1)
if [[ "$RESULT" == *"\"value\": \"baz\""* ]]; then
  echo "  PASS"
else
  HAS_FAILED=1
  echo "  FAIL: " $RESULT
fi

echo

echo "TEST: write key; overwrite key and read"
RESULT=$(./bin/kvclient write --key foo/bar --value baz 2>&1)
if [[ "$RESULT" != "{}" ]]; then
  HAS_FAILED=1
  echo "  FAIL: " $RESULT
fi
RESULT=$(./bin/kvclient write --key foo/bar --value bam 2>&1)
if [[ "$RESULT" != "{}" ]]; then
  HAS_FAILED=1
  echo "  FAIL: " $RESULT
fi
RESULT=$(./bin/kvclient get --key foo/bar 2>&1)
if [[ "$RESULT" == *"\"value\": \"bam\""* ]]; then
  echo "  PASS"
else
  HAS_FAILED=1
  echo "  FAIL: " $RESULT
fi

echo

echo "TEST: write key; delete key and read"
RESULT=$(./bin/kvclient write --key foo/bar --value baz 2>&1)
if [[ "$RESULT" != "{}" ]]; then
  HAS_FAILED=1
  echo "  FAIL: " $RESULT
fi
RESULT=$(./bin/kvclient delete --key foo/bar 2>&1)
if [[ "$RESULT" != "{}" ]]; then
  HAS_FAILED=1
  echo "  FAIL: " $RESULT
fi
RESULT=$(./bin/kvclient get --key foo/bar 2>&1)
if [[ "$RESULT" == "kvclient: error: [404]: [404]: key 'foo/bar' not found" ]]; then
  echo "  PASS"
else
  HAS_FAILED=1
  echo "  FAIL: " $RESULT
fi

kill $KVSERVICE_PID

if [ $HAS_FAILED = 1 ]; then
  exit 255
fi