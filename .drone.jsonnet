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
      name: "compile",
      image: "golang:" + version,
      commands: [
        "make awl"
      ],
    },
    {
      name: "lint",
      image: "rancher/drone-golangci-lint:latest",
      depends_on: [
        "compile",
      ],
    },
    {
      name: "test",
      image: "golang:" + version,
      commands: [
        "make test-ci"
      ],
      depends_on: [
        "lint",
      ],
    },
    {
      name: "fuzz",
      image: "golang:" + version,
      commands: [
        "make fuzz",
      ],
      depends_on: [
        "lint",
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
      ]
    },
    {
      name: "test",
      image: "golang",
      commands: [
        "make test"
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
  // testing("1.19", "arm64"),
  // testing("1.18", "amd64"),
  // testing("1.18", "arm64"),

  release()
]