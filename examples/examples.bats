#!/usr/bin/env bats

setup() {
  cd "$BATS_TEST_DIRNAME"
  BASEDIR="$(git rev-parse --show-toplevel)"
  EXE="$BASEDIR/k6"

  if [ ! -x "$EXE" ]; then
    echo "    - building k6" >&3
    cd "$BASEDIR"
    xk6 build --with github.com/grafana/xk6-plugin=.
    cd "$BATS_TEST_DIRNAME"
  fi
}

test_math() {
  export K6_PLUGIN_RUNTIME=$1
  run $EXE run math.js
  [ $status -eq 0 ]
  echo "$output" | grep -q "2 + 2 = 4"
}

test_dns() {
  export K6_PLUGIN_RUNTIME=$1
  run $EXE run dns.js
  [ $status -eq 0 ]
  echo "$output" | grep -q 'k6.io: {'
}

test_dump() {
  export K6_PLUGIN_RUNTIME=$1
  let before=$(ls dump 2>/dev/null | wc -l)
  run $EXE run dump.js
  [ $status -eq 0 ]
  let after=$(ls dump 2>/dev/null | wc -l)
  [ $after -gt $before ]
}

@test 'math deno' {
  test_math deno
}

@test 'math bun' {
  test_math bun
}

@test 'dns deno' {
  test_dns deno
}

@test 'dns bun' {
  test_dns bun
}

@test 'dump deno' {
  test_dump deno
}

@test 'dump bun' {
  test_dump bun
}
