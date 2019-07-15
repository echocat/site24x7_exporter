package main

func deploy(branch string) {
	deployDocker(branch, dv)
}

func deployDocker(branch string, v dockerVariant) {
	deployDockerTag(v.imageName(branch))
	executeForVersionParts(branch, func(tagSuffix string) {
		deployDockerTag(v.imageName(tagSuffix))
	})
	if latestVersionPattern != nil && latestVersionPattern.MatchString(branch) {
		deployDockerTag(v.baseImageName())
	}
	if v.main {
		deployDockerTag(imagePrefix + ":" + branch)
		executeForVersionParts(branch, func(tagSuffix string) {
			deployDockerTag(imagePrefix + ":" + tagSuffix)
		})
		if latestVersionPattern != nil && latestVersionPattern.MatchString(branch) {
			deployDockerTag(imagePrefix + ":latest")
		}
	}
}

func deployDockerTag(tag string) {
	execute("docker", "push", tag)
}
