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
      image: 'golangci/golangci-lint',
      commands: [
        'golangci-lint run ./...',
      ],
    },
    {
      name: 'cache',
      image: 'golang:' + version,
      commands: [
        'go mod tidy',
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
  clone: {
    disable: true,
  },
  trigger: {
    event: [
      'tag',
    ],
  },
  steps: [
    {
      name: 'clone',
      image: 'woodpeckerci/plugin-git',
      settings: {
        tags: true,
      },
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
        'apk add --no-cache scdoc',
        'goreleaser release',
      ],
      volumes: [
        {
          name: 'cache',
          path: '/go',
        },
      ],
    },
    {
      name: 'trigger',
      image: 'plugins/downstream',
      settings: {
        server: 'ci.git.froth.zone',
        token: {
          DRONE_TOKEN: {
            from_secret: 'DRONE_TOKEN',
          },
        },
        fork: true,
        repositories: [
          'packages/awl',
        ],
        parameters: [
          'TAG=${DRONE_TAG}',
        ],
      },
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
