local pipeline(version, arch) = {
  kind: "pipeline",
  name: version + "-" + arch ,
  platform: {
    arch: arch
  },
  steps: [
    {
      name: "test",
      image: "golang:" + version,
      commands: [
        "go test ./..."
      ]
    }
  ]
};

// logawl uses generics so 1.18 is the minimum
[
  pipeline("1.18", "amd64"),
  pipeline("1.18", "arm64"),
]