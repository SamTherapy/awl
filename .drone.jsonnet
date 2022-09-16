// SPDX-License-Identifier: BSD-3-Clause

local testing(version, arch) = {
  kind: 'pipeline',
  type: 'docker',
  name: version + '-' + arch,
  platform: {
    arch: arch,
  },
  steps: [
    {
      name: 'lint',
      image: 'rancher/drone-golangci-lint:latest',
    },
    {
      name: 'cache',
      image: 'golang:' + version,
      commands: [
        'go mod tidy'
      ],
      depends_on: [
        'lint',
      ],
      volumes: [
        {
          name: 'cache',
          path: '/go',
        },
      ],
    },
    {
      name: 'test',
      image: 'golang:' + version,
      commands: [
        'make test-ci',
      ],
      depends_on: [
        'cache',
      ],
      volumes: [
        {
          name: 'cache',
          path: '/go',
        },
      ],
    },
    {
      name: 'fuzz',
      image: 'golang:' + version,
      commands: [
        'make fuzz-ci',
      ],
      depends_on: [
        'cache',
      ],
      volumes: [
        {
          name: 'cache',
          path: '/go',
        },
      ],
    },
  ],
  trigger: {
    event: {
      exclude: [
        'tag',
      ],
    },
  },
  volumes: [
    {
      name: 'cache',
      temp: {},
    },
  ],
};

// "Inspired by" https://goreleaser.com/ci/drone/
local release() = {
  kind: 'pipeline',
  type: 'docker',
  name: 'release',
  trigger: {
    event: [
      'tag',
    ],
  },
  steps: [
    {
      name: 'fetch',
      image: 'alpine/git',
      commands: [
        'git fetch --tags',
      ],
    },
    {
      name: 'test',
      image: 'golang',
      commands: [
        'make test-ci',
      ],
      volumes: [
        {
          name: 'cache',
          path: '/go',
        },
      ],
    },
    {
      name: 'release',
      image: 'goreleaser/goreleaser',
      environment: {
        GITEA_TOKEN: {
          from_secret: 'GITEA_TOKEN',
        },
      },
      commands: [
        'goreleaser release',
      ],
      volumes: [
        {
          name: 'cache',
          path: '/go',
        },
      ],
    },
  ],
    volumes: [
    {
      name: 'cache',
      temp: {},
    },
  ],
};

[
  testing('1.19', 'amd64'),
  testing('1.19', 'arm64'),
  testing('1.18', 'amd64'),
  testing('1.18', 'arm64'),

  release(),
]
