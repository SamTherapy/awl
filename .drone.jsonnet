// SPDX-License-Identifier: BSD-3-Clause

local testing(version, arch) = {
  kind: "pipeline",
  type: "docker",
  name: version + "-" + arch ,
  platform: {
    arch: arch
  },
  steps: [
    {
      name: "submodules",
      image: "alpine/git",
      commands: [
        "git submodule update --init --recursive"
      ]
    },
    {
      name: "lint",
      image: "rancher/drone-golangci-lint:latest",
      depends_on: [
        "submodules",
      ],
    },
    {
      name: "test",
      image: "golang:" + version,
      commands: [
        "go test -v -race ./... -cover"
      ],
      depends_on: [
        "submodules",
      ],
    },
  ],
  trigger: {
    event: {
      exclude: [
        "tag"
      ],
    }
  },
};

// "Inspired by" https://goreleaser.com/ci/drone/
local release() = {
  kind: "pipeline",
  type: "docker",
  name: "release",
  trigger: {
    event: [
      "tag"
    ],
  },
  steps: [
    {
      name: "fetch",
      image: "alpine/git",
      commands : [
        "git fetch --tags",
        "git submodule update --init --recursive"
      ]
    },
    {
      name: "test",
      image: "golang",
      commands: [
        "go test -race ./... -cover"
      ]
    },
    {
      name: "release",
      image: "goreleaser/goreleaser",
      environment: {
        "GITEA_TOKEN": {
          from_secret: "GITEA_TOKEN"
        }
      },
      commands: [
        "goreleaser release"
      ],
    }
  ]
};

[
  testing("1.19", "amd64"),
  testing("1.19", "arm64"),
  testing("1.18", "amd64"),
  testing("1.18", "arm64"),

  
  release()
]